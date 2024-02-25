package main

import (
	_ "github.com/glebarez/go-sqlite"
	"github.com/rs/zerolog/log"
	"service/configuration"
	"service/controller"
	"service/database"
	"service/misc"
)

func main() {
	misc.WelcomeMessage()
	//Read YAML, Construct conf object + Setup Logging
	conf, closeFile := configuration.Setup()
	defer closeFile()

	// Setup sqlite
	if conf.Database.Path == "" || conf.Database.InitQuery == "" {
		log.Warn().Msg("Missing Database in main.yml : Models are disabled")
	} else {
		// call closeDB to defer the db close
		_, closeDB := database.Setup(conf.Database.InitQuery, conf.Database.Path)
		defer closeDB()
	}

	//Attach Controllers to HTTP + Start the HTTP server
	controller.SetupControllers(conf.Server.Services)
	misc.StartHttp(conf.Server.Port)
}
