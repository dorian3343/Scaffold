package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/glebarez/go-sqlite"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"service/configuration"
	"service/database"
	"slices"
	"strconv"
)

func main() {
	//Setup actual config
	conf, err := configuration.Create("./main.yml")
	if err != nil {
		log.Fatal().Err(err).Msg("Something went wrong with generating config from main.yml")
	}
	targetLog := conf.Server.TargetLog

	var multi zerolog.LevelWriter
	if targetLog != "" {
		/* Setup logging :  Get logging file and set MultiLevelWriting*/
		file, err := os.OpenFile(targetLog, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal().Err(err).Msg("Error opening log file")
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Fatal().Err(err).Msg("Error while closing log file")
			}
		}(file)

		multi = zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stdout}, file)
	} else {
		multi = zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	log.Logger = zerolog.New(multi).With().Timestamp().Logger()

	/* Setup sqlite */
	if conf.Database.Path == "" || conf.Database.InitQuery == "" {
		log.Warn().Msg("Missing Database in main.yml : Models are disabled")
	} else {
		db := database.Setup(conf.Database.InitQuery, conf.Database.Path)
		if db == nil {
			log.Fatal().Msg("Database is a nil pointer")
			return
		}
		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {
				log.Fatal().Err(err).Msg("Fatal Error closing sqlite")
			}
		}(db)
	}
	/* Setup http handling : Controllers + Models + main http server */
	fmt.Println("---------Tree------------")
	var wrn []string

	for route, handler := range conf.Server.Services {

		log.Trace().Msgf("New Route : | %s : %s", route, handler.Name)
		if route == "" {
			log.Fatal().Err(errors.New("Missing route")).Msg("Something went wrong with setting up Controllers")
		}

		http.Handle(route, handler)
		if handler.Name == "" {
			wrn = append(wrn, fmt.Sprintf("Empty controller for Route: '%s'", route))
		}

		// Check for empty fallbacks
		e := []byte("null")
		if slices.Equal(e, handler.Fallback) {
			wrn = append(wrn, fmt.Sprintf("Empty Fallback for Route: '%s'", route))
		}

	}
	fmt.Println("---------Tree------------")

	/* print tree warnings */
	if len(wrn) != 0 {
		for i := 0; i < len(wrn); i++ {
			log.Warn().Msg(wrn[i])
		}
	}

	log.Info().Msg("Finished initialization : Starting server...")
	err = http.ListenAndServe(":"+strconv.Itoa(conf.Server.Port), nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Fatal Error with http server")
	}

}
