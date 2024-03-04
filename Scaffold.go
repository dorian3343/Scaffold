package main

import (
	"flag"
	"fmt"
	_ "github.com/glebarez/go-sqlite"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"service/components/controller"
	"service/configuration"
	"service/misc"
	"time"
)

func printVersion() {
	body, err := os.ReadFile("VERSION")
	if err != nil {
		fmt.Println("Something went wrong reading version: " + err.Error())
	} else {
		fmt.Println(string(body))
	}
}

func main() {
	//Print version
	if len(os.Args) > 1 && os.Args[1] == "version" {
		printVersion()
		os.Exit(0)
	}
	if len(os.Args) > 1 && os.Args[1] == "run" {
		entrypoint := "./main.yml"

		if len(os.Args) > 2 {
			entrypoint = os.Args[2] + "/main.yml"
		}

		// Core of Scaffold,this should be used in e2e testing not the cli

		// measure time time to start

		start := time.Now()
		// Set pretty logging straight away
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		misc.WelcomeMessage()
		//Read YAML, Construct conf object, setup sqlite + Setup Logging
		conf, closeFile := configuration.Setup(entrypoint)
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
	if flag.NArg() == 0 {
		// Print help message
		fmt.Println(`Scaffold is a tool for building APIs fast and easy

Usage:
	Scaffold <command> [argument]

List of commands:
	version   print out your scaffold version
	run       run the Scaffold from a config in a specified directory`)
		fmt.Println("\nTool by Dorian Kalaczyński")
		os.Exit(0)
	}

}
