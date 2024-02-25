package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/metalim/jsonmap"
	"reflect"
	"testing"
)

const sampleJSON = `{"name":"John","Age":30}`

func TestModel_Querybuilder(t *testing.T) {
	// Test case for successful scenario
	json1 := jsonmap.New()
	err := json.Unmarshal([]byte(sampleJSON), json1)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		//panic(err)
	}

	// marshal, keeping order
	output, err := json.Marshal(json1)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		//panic(err)
	}
	model := Model{template: "SELECT * FROM users WHERE name = %s AND age = %s", json: &json1}
	expectedQuery := "SELECT * FROM users WHERE name = John AND age = 30"
	query, err := model.Querybuilder(output)
	if err != nil {
		fmt.Println(err)
		t.Errorf("Unexpected error: %v", err)
	}
	if query != expectedQuery {
		t.Errorf("Expected query: %s, got: %s", expectedQuery, query)
	}
	fmt.Println("Test completed successfully.")
}

func TestMapToArray(t *testing.T) {
	json1 := jsonmap.New()

	err := json.Unmarshal([]byte(sampleJSON), json1)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		//panic(err)
	}

	mapExpected := []string{"John", "30"}
	res, err := MapToArray(json1)
	if err != nil {
		t.Errorf("Error converting map to array: %v", err)
	}
	if !reflect.DeepEqual(res, mapExpected) {
		t.Errorf("Map conversion result incorrect. Expected: %v, Got: %v", mapExpected, res)
	}

	// Test case for unsupported type
	unsupportedType := 10
	_, err = MapToArray(unsupportedType)
	if err == nil {
		t.Error("Expected error for unsupported type, but got none")
	}
	expectedError := errors.New("unsupported type")
	if err.Error() != expectedError.Error() {
		t.Errorf("Expected error message: %v, Got: %v", expectedError, err)
	}
}