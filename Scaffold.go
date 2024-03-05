package main

import (
	"flag"
	"fmt"
	_ "github.com/glebarez/go-sqlite"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"service/cmd"
	"service/components/controller"
	"service/configuration"
	"service/misc"
	"time"
)

func main() {
	//Print version
	if len(os.Args) > 1 && os.Args[1] == "version" {
		cmd.PrintVersion()
		os.Exit(0)
	}
	//Create a new empty project
	if len(os.Args) > 1 && os.Args[1] == "init" {
		if len(os.Args) > 2 {
			cmd.ProjectInit(os.Args[2])
			os.Exit(0)
		} else {
			fmt.Println("Error: Missing project name")
			os.Exit(1)
		}
	}

	//Generate docs for a project
	if len(os.Args) > 1 && os.Args[1] == "auto-doc" {
		if len(os.Args) > 2 {
			cmd.GenerateDoc(os.Args[2])
			os.Exit(0)
		} else {
			fmt.Println("Error: Missing project name")
			os.Exit(1)
		}
	}

	// Run a scaffold app in current dir or a specified
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
		cmd.PrintGuide()
	}

}
