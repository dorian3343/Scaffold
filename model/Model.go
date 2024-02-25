package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
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

func structToArray(s interface{}) ([]string, error) {
	switch v := s.(type) {
	case map[string]interface{}:
		// Handle map type
		var values []string
		for x, val := range v {
			fmt.Println(x)
			values = append(values, fmt.Sprintf("%v", val))
		}
		return values, nil
	case struct{}:
		// Handle struct type
		st := reflect.TypeOf(s)
		sv := reflect.ValueOf(s)
		values := make([]string, st.NumField())
		for i := 0; i < st.NumField(); i++ {
			fieldValue := sv.Field(i)
			values[i] = fmt.Sprintf("%v", fieldValue.Interface())
		}
		return values, nil
	default:
		return nil, errors.New("unsupported type")
	}
}

func (m Model) Querybuilder(x []byte) (string, error) {
	var jsonData interface{}

	err := json.Unmarshal(x, &jsonData)
	if err != nil {
		return "", errors.New("failed to decode JSON data: " + err.Error())
	}

	arrayData, err := structToArray(jsonData)
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

	// Format the array data into a string using the provided template
	return fmt.Sprintf(m.template, args...), nil
}

func (m Model) Query(query string) (*sql.Rows, error) {
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
