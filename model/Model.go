package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/metalim/jsonmap"
	"github.com/rs/zerolog/log"
	"reflect"
	"slices"
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
func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]-'a'+'A') + s[1:]
}
func generateStructFromJsonMap(f jsonmap.Map) reflect.Type {
	fields := make([]reflect.StructField, 0, len(f.Keys()))
	for _, Name := range f.Keys() {
		x, ok := f.Get(Name)
		if !ok {
			log.Fatal().Msg("Something went wrong ")
		} else {
			var T any
			switch x {
			case "string":
				T = ""
			case "integer":
				T = 0
			default:
				log.Fatal().Msg("Unrecognized type in JSON template")
			}
			field := reflect.StructField{
				Name: capitalize(Name),
				Type: reflect.TypeOf(T),
			}
			fields = append(fields, field)

		}
	}
	structType := reflect.StructOf(fields)
	return structType
}

// fills out the query template with data from the json
func (m Model) Querybuilder(x []byte) (string, error) {
	if len(x) == 0 {
		log.Warn().Msg("Empty JSON template so Query is sent as is.")
		return m.template, nil
	}
	json1 := jsonmap.New()
	err := json.Unmarshal(x, json1)
	if err != nil {
		return "", errors.New("failed to decode JSON data: " + err.Error())
	}
	err = json.Unmarshal(x, json1)
	if err != nil {
		return "", errors.New("failed to decode JSON data: " + err.Error())
	}
	var Keys = make([]string, 0)
	val := json1.Keys()
	for i := 0; i < len(val); i++ {
		str := fmt.Sprintf("%v", val[i])
		Keys = append(Keys, capitalize(str))
	}
	var TKeys = make([]string, 0)
	T := generateStructFromJsonMap(*m.json)
	fmt.Println(T)
	for i := 0; i < T.NumField(); i++ {
		fmt.Println(T.Field(i).Name)
		TKeys = append(TKeys, T.Field(i).Name) //
	}
	if slices.Equal(TKeys, Keys) {
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
		log.Warn().Msg("JSON request does not match spec")
		return "", errors.New("JSON request does not match spec")

	}

}

func (m Model) Query(query string) (*sql.Rows, error) {
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
