package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/metalim/jsonmap"
	"github.com/rs/zerolog/log"
	"strings"
)

type Model struct {
	Name     string
	db       *sql.DB
	template string
	json     *jsonmap.Map
}

func Create(name string, db *sql.DB, template string, JSON *jsonmap.Map) Model {
	return Model{Name: name, db: db, template: template, json: JSON}
}

// Fills out the query template with data from the json
func (m Model) Querybuilder(x []byte) (string, error) {
	json1 := jsonmap.New()

	if len(x) == 0 {
		log.Warn().Msg("Empty JSON template so Query is sent as is.")
		return m.template, nil
	}

	err := json.Unmarshal(x, json1)
	if err != nil {
		return "", errors.New("failed to decode JSON data: " + err.Error())
	}

	T := generateStructFromJsonMap(*m.json)

	if matchesSpec(*json1, T) {
		arrayData, err := MapToArray(json1)
		if err != nil {
			return "", err
		}
		// Get the number of placeholders in the template string
		numPlaceholders := strings.Count(m.template, "%s")
		// Truncate arrayData if it exceeds the number of placeholders
		if numPlaceholders < len(arrayData) {
			arrayData = arrayData[:numPlaceholders]
		}

		// Convert arrayData to a slice of interface{}
		args := make([]interface{}, len(arrayData))
		for i, v := range arrayData {
			args[i] = v
		}
		return fmt.Sprintf(m.template, args...), nil

	} else {
		err := "JSON request does not match spec"
		log.Warn().Err(errors.New(err)).Msg("Malformed JSON")
		return "", errors.New(err)
	}

}

// Queries the database
func (m Model) Query(query string) (*sql.Rows, error) {
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
