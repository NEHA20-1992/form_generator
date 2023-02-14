package main

import (
	l "github.com/NEHA20-1992/form_generator/pkg/logger"
	"github.com/NEHA20-1992/form_generator/pkg/server"
)

func main() {
	l.BootstrapLogger.Infoln("Tausi Server starting")
	server.ApiServer.Run()
}
