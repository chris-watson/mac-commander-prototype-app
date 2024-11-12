package infra

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chris-watson/mac-windows-installer-app/pkg/adapter/handler"
	sysAdapter "github.com/chris-watson/mac-windows-installer-app/pkg/adapter/system"
	service "github.com/chris-watson/mac-windows-installer-app/pkg/service"
)

func StartServer(port int) {

	platform := getPlatform()
	fmt.Println("Running on", platform)

	var cmdSvc *service.CommanderService
	switch platform {
	case Windows:
		cmdSvc = service.NewCommanderService(sysAdapter.NewWindowsCommander())
	case MacOS:
		cmdSvc = service.NewCommanderService(sysAdapter.NewMacCommander())
	}

	handler := handler.NewHandler(cmdSvc)

	mux := http.NewServeMux()
	mux.HandleFunc("/execute", handler.HandleCommand)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	// channel to listen for errors coming from the listener
	serverErrors := make(chan error, 1)

	// start server in a goroutine
	go func() {
		log.Println("Server starting on port:", port)
		serverErrors <- srv.ListenAndServe()
	}()

	// channel to listen for interrupt/terminate signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// block until we receive a signal or an error
	select {
	case err := <-serverErrors:
		log.Fatalf("Error starting server: %v", err)
	case <-shutdown:
		log.Println("Starting shutdown...")

		// create deadline for shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// graceful shutdown
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Could not stop server gracefully: %v", err)
			if err := srv.Close(); err != nil {
				log.Printf("Could not force close server: %v", err)
			}
		}
	}
}
