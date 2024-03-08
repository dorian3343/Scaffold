package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/metalim/jsonmap"
	"github.com/rs/zerolog/log"
	"reflect"
	"strings"
)

type Model struct {
	Name               string
	db                 *sql.DB
	queryTemplate      string
	json               *jsonmap.Map
	generatedTypeCache *reflect.Type
}

func Create(name string, db *sql.DB, template string, JSON *jsonmap.Map) Model {
	return Model{Name: name, db: db, queryTemplate: template, json: JSON, generatedTypeCache: nil}
}
func (m Model) GetQuery() string {
	return m.queryTemplate
}

// Fills out the query queryTemplate with data from the json
func (m Model) Querybuilder(x []byte) (string, error) {
	jsonRequest := jsonmap.New()

	if len(x) == 0 {
		log.Warn().Msg("Empty JSON queryTemplate so Query is sent as is.")
		return m.queryTemplate, nil
	}

	err := json.Unmarshal(x, jsonRequest)
	if err != nil {
		log.Debug().Msg("Failed to decode :" + string(x))
		return "", errors.New("Failed to decode JSON data: " + err.Error())
	}
	//Basic type caching
	var GeneratedType reflect.Type
	if m.generatedTypeCache == nil {
		GeneratedType = GenerateStructFromJsonMap(*m.json)
		m.generatedTypeCache = &GeneratedType
		log.Trace().Msg("Caching Type for model : " + m.Name)
	} else {
		GeneratedType = *m.generatedTypeCache
	}

	if matchesSpec(*jsonRequest, GeneratedType) {
		arrayData, err := MapToArray(jsonRequest)
		if err != nil {
			return "", err
		}
		// Get the number of placeholders in the queryTemplate string
		numPlaceholders := strings.Count(m.queryTemplate, "%s")
		// Truncate arrayData if it exceeds the number of placeholders
		if numPlaceholders < len(arrayData) {
			arrayData = arrayData[:numPlaceholders]
		}

		// Convert arrayData to a slice of interface{}
		args := make([]interface{}, len(arrayData))
		for i, v := range arrayData {
			args[i] = v
		}
		return fmt.Sprintf(m.queryTemplate, args...), nil

	} else {
		err := "JSON request does not match spec"
		log.Warn().Err(errors.New(err)).Msg("Malformed JSON")
		return "", errors.New(err)
	}

}

func (m Model) Query(query string) (*sql.Rows, error) {
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
