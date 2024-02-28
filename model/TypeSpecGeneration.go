package model

import (
	"errors"
	"fmt"
	"github.com/metalim/jsonmap"
	"github.com/rs/zerolog/log"
	"reflect"
	"service/misc"
	"slices"
)

/*
	The file TypeSpecGeneration should contain any function that a model uses to handle
	matching the json request against the generated specification
*/

// Function made to match if json request matches a struct's specification
func matchesKeys(Y []jsonmap.Key, T reflect.Type) bool {
	var Keys = make([]string, 0)
	var TKeys = make([]string, 0)

	for i := 0; i < T.NumField(); i++ {
		TKeys = append(TKeys, T.Field(i).Name)
	}

	for i := 0; i < len(Y); i++ {
		str := fmt.Sprintf("%v", Y[i])
		Keys = append(Keys, misc.Capitalize(str))
	}
	return slices.Equal(Keys, TKeys)
}

// Generates a type from a jsonmap.map, the intended usage is to  generate types from the configuration,
// to match requests against a spec
func generateStructFromJsonMap(f jsonmap.Map) reflect.Type {
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
