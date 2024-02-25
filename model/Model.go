package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/metalim/jsonmap"
	"strings"
)

type Model struct {
	db       *sql.DB
	template string
	json     interface{}
}

func Create(db *sql.DB, template string, JSON interface{}) Model {
	return Model{db: db, template: template, json: JSON}
}

// Converts maps  into Array
func MapToArray(s interface{}) ([]string, error) {
	if s == nil {
		return nil, errors.New("Nil value Passed")
	}
	switch v := s.(type) {
	case *jsonmap.Map:
		if v == nil {
			return nil, errors.New("Nil value passed")
		}
		// Handle map type
		var values = make([]string, 0)
		val := v.Values()
		for i := 0; i < len(val); i++ {
			str := fmt.Sprintf("%v", val[i])
			values = append(values, str)
		}

		return values, nil
	default:
		return nil, errors.New("unsupported type")
	}
}

// fills out the query template with data from the json
func (m Model) Querybuilder(x []byte) (string, error) {
	json1 := jsonmap.New()
	err := json.Unmarshal(x, json1)
	if err != nil {
		return "", errors.New("failed to decode JSON data: " + err.Error())
	}
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
}

func (m Model) Query(query string) (*sql.Rows, error) {
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
