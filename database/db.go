package database

import (
	"database/sql"
	"github.com/rs/zerolog/log"
	"os"
)

func Setup(query string, databaseName string) (*sql.DB, func()) {
	/* Check if Database file exists */
	_, err := os.Stat(databaseName)
	if os.IsNotExist(err) {
		_, err2 := os.Create(databaseName)
		if err2 != nil {
			log.Fatal().Err(err).Msg("Something went wrong with creating Database")
		}
	}
	// Create the Connection
	db, err := sql.Open("sqlite", databaseName)
	if err != nil {
		log.Fatal().Err(err).Msg("Fatal Error opening sqlite")
	}

	_, err = db.Exec(query)
	if err != nil {
		log.Fatal().Err(err).Msg("Fatal Error during table setup")
	}

	closeDB := func() {
		if err := db.Close(); err != nil {
			log.Fatal().Err(err).Msg("Failed to close the database")
		}
	}
	return db, closeDB
}
