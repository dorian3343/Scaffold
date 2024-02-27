package main

import (
	_ "github.com/glebarez/go-sqlite"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"service/configuration"
	"service/controller"
	"service/misc"
)

// This is not part of the codebase, this is just an automated test

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	misc.WelcomeMessage()
	//Read YAML, Construct conf object, setup sqlite + Setup Logging
	conf, closeFile := configuration.Setup("./main.yml")
	defer closeFile()
	defer conf.DatabaseClosure()

	//Attach Controllers to HTTP + Start the HTTP server
	controller.SetupControllers(conf.Server.Services)

	misc.StartHttp(conf.Server.Port)
}
