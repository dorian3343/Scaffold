package model

import (
	"errors"
	"fmt"
	"github.com/metalim/jsonmap"
	"github.com/rs/zerolog/log"
	"reflect"
	"service/misc"
)

/*
	The file TypeSpecGeneration should contain any function that a model uses to handle
	matching the json request against the generated specification
*/

// Matches Json Request to a struct (JSON specification)
func matchesSpec(Y jsonmap.Map, T reflect.Type) bool {
	// Check if the number of keys in the JSON map matches the number of fields in the struct
	if len(Y.Keys()) != T.NumField() {
		return false
	}
	for i := 0; i < T.NumField(); i++ {
		fieldName := T.Field(i).Name
		fieldT := T.Field(i).Type

		//check if field exists
		if fieldValue, ok := Y.Get(fieldName); !ok {
			log.Error().Msgf("Missing field '%s' in JSON request", fieldName)
			return false
		} else {
			// Compare the types of the field in the struct and in the JSON map
			if fieldT != reflect.TypeOf(fieldValue) {
				if fieldT == reflect.TypeOf(int(0)) && reflect.TypeOf(fieldValue) == reflect.TypeOf(float64(0)) {
					log.Warn().Msg("Adapted type to int from float64")
				} else {
					log.Error().Msgf("Wrong Type in field '%s' in JSON request. Got type '%s' expected type '%s'", fieldName, fieldT, reflect.TypeOf(fieldValue))
					return false
				}
			}
		}
	}

	return true
}

// Generates a type from a jsonmap.map, the intended usage is to  generate types from the configuration,
// to match requests against a spec
func GenerateStructFromJsonMap(f jsonmap.Map) reflect.Type {
	fields := make([]reflect.StructField, 0, len(f.Keys()))
	for _, Name := range f.Keys() {
		x, ok := f.Get(Name)
		if !ok {
			//If this error shows up something went very wrong
			log.Fatal().Msg("Missing field in jsonmap")
		} else {
			//Variable made to set the StructFields type
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
				Name: misc.Capitalize(Name),
				Type: reflect.TypeOf(T),
			}
			fields = append(fields, field)

		}
	}
	structType := reflect.StructOf(fields)
	return structType
}

// Converts jsonmaps into arrays, used to build a query out of the values
func MapToArray(s interface{}) ([]string, error) {
	if s == nil {
		return nil, errors.New("Nil value Passed")
	}
	switch v := s.(type) {
	case *jsonmap.Map:
		if v == nil {
			return nil, errors.New("Nil value passed")
		}
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
