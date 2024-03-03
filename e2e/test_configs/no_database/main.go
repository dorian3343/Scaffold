package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"service/components/controller"
	"service/configuration"
	"service/misc"
	"time"
)

// This is not part of the codebase, this is just an automated test

func main() {
	// measure time time to start
	start := time.Now()
	// Set pretty logging straight away
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	misc.WelcomeMessage()
	//Read YAML, Construct conf object, setup sqlite + Setup Logging
	conf, closeFile := configuration.Setup("./main.yml")
	defer closeFile()
	defer conf.DatabaseClosure()

	//Attach Controllers to HTTP + Start the HTTP server
	controller.SetupControllers(conf.Server.Services)

	//Serve static
	if conf.Server.Static != "" {
		// Check if the folder exists
		if _, err := os.Stat(conf.Server.Static); err != nil {
			if os.IsNotExist(err) {
				log.Fatal().Msgf("Static folder '%s' does not exist", conf.Server.Static)
			} else {
				log.Fatal().Msgf("Error checking static folder '%s': %v", conf.Server.Static, err)
			}
		}

		// If the folder exists, serve files from it
		http.Handle("/", http.FileServer(http.Dir(conf.Server.Static)))
	}

	end := time.Now()
	elapsed := end.Sub(start)
	log.Info().Msgf("Project built in : %s", elapsed)

	misc.StartHttp(conf.Server.Port)
}
