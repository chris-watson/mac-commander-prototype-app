package app

import (
	"github.com/chris-watson/mac-windows-installer-app/pkg/infra"
)

func Start(port int) {
	infra.StartServer(port)
}
