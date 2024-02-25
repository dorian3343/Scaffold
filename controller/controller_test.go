package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

/*Testing for a no model Controller using the fallback string*/
func TestController_ServeHTTP_BasicString(t *testing.T) {
	x, _ := json.Marshal("Hello World")
	c := Create("/basicTest", nil, x, false)
	req := httptest.NewRequest("GET", "http://google.com", nil)
	w := httptest.NewRecorder()
	c.ServeHTTP(w, req)
	res := w.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
	}

	var bodyBytes []byte
	if res.Body != nil {
		bodyBytes, _ = io.ReadAll(res.Body)
	}

	// Check the response body
	expectedBody := "\"Hello World\"" //
	if string(bodyBytes) != expectedBody {
		t.Errorf("Expected response body %s, got %s", expectedBody, string(bodyBytes))
	}
}

/*Testing for a no model Controller using the fallback string*/
func TestController_ServeHTTP_BasicInt(t *testing.T) {
	x, _ := json.Marshal(69)
	c := Create("/basicTest", nil, x, false)
	req := httptest.NewRequest("GET", "http://google.com", nil)
	w := httptest.NewRecorder()
	c.ServeHTTP(w, req)
	res := w.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
	}

	var bodyBytes []byte
	if res.Body != nil {
		bodyBytes, _ = io.ReadAll(res.Body)
	}

	// Check the response body
	expectedBody := "69"
	if string(bodyBytes) != expectedBody {
		t.Errorf("Expected response body %s, got %s", expectedBody, string(bodyBytes))
	}
}

func TestController_ServeHTTP_Struct(t *testing.T) {
	// Basic struct to test
	type SomeStruct struct {
		Field1 string
		Field2 string
	}
	// Prepare the input for the controller
	inputData := SomeStruct{Field1: "value1", Field2: "value2"}
	// Marshal the input to JSON
	requestData, err := json.Marshal(inputData)
	if err != nil {
		t.Errorf("Error marshaling input data: %v", err)
		return
	}

	// Create a request using the input data
	c := Create("/basicTest", nil, requestData, false)
	req := httptest.NewRequest("GET", "http://google.com", nil)
	w := httptest.NewRecorder()

	// Call the ServeHTTP method of the controller
	c.ServeHTTP(w, req)
	res := w.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
	}

	// Parse the response body as JSON
	var responseBody SomeStruct
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&responseBody); err != nil {
		t.Errorf("Error decoding response body: %v", err)
		return
	}

	// Check if the response matches the expected output
	expectedOutput := SomeStruct{Field1: "value1", Field2: "value2"}
	if !reflect.DeepEqual(responseBody, expectedOutput) {
		t.Errorf("Expected response body %#v, got %#v", expectedOutput, responseBody)
	}
}

func TestController_Create(t *testing.T) {
	expectedName := "/test"
	expectedFallback, _ := json.Marshal(69)

	// Call the Create function
	c := Create(expectedName, nil, expectedFallback, false)

	// Check if the fields of the created controller match the expected values
	if c.Name != expectedName {
		t.Errorf("Expected name %s, got %s", expectedName, c.Name)
	}

	if !reflect.DeepEqual(c.Fallback, expectedFallback) {
		t.Errorf("Expected fallback %#v, got %#v", expectedFallback, c.Fallback)
	}
}
