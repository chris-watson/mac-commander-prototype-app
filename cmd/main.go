package main

import (
	"os"
	"strconv"

	"github.com/chris-watson/mac-windows-installer-app/app"
)

func main() {
	port := 8080
	if envPort := os.Getenv("COMMANDER_APP_PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			port = p
		}
	}
	app.Start(port)
}
