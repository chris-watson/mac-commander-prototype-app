package integration_test

import (
	"os"
	"testing"
	"time"

	"github.com/chris-watson/mac-windows-installer-app/app"
)

const testPort = 8899

func TestMain(m *testing.M) {

	// start server in a goroutine
	go func() {
		app.Start(testPort)
	}()

	// give server some time to start
	time.Sleep(100 * time.Millisecond)

	// run tests
	code := m.Run()

	// stop server
	// Note: In a production application,
	// we would implement a graceful shutdown
	// in our application code and use that here.
	os.Exit(code)

}
