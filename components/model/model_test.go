package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/metalim/jsonmap"
	"reflect"
	"testing"
)

const sampleJSON = `{"Name":"John","Age":30}`

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

func TestMatchesSpec(t *testing.T) {
	testjsonmap := jsonmap.New()
	testjsonmap.Set("Name", "string")
	testjsonmap.Set("Age", "integer")

	testjsonmap2 := jsonmap.New()
	testjsonmap2.Set("Age", "string")
	testjsonmap2.Set("Name", "integer")

	testjsonmap3 := jsonmap.New()
	testjsonmap2.Set("ageless", "string")
	testjsonmap2.Set("nameful", "integer")

	json1 := jsonmap.New()
	err := json.Unmarshal([]byte(sampleJSON), json1)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
	}
	T := generateStructFromJsonMap(*testjsonmap)
	if !matchesSpec(*json1, T) {
		t.Error("Request does not match spec")
	}
	T = generateStructFromJsonMap(*testjsonmap2)
	if matchesSpec(*json1, T) {
		t.Error("Request does not match spec")
	}
	T = generateStructFromJsonMap(*testjsonmap3)
	if matchesSpec(*json1, T) {
		t.Error("Request does not match spec")
	}
}
