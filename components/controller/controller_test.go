package controller

import (
	"bytes"
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
	c := Create("basicTest", nil, x, "", "", "")
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
	c := Create("basicTest", nil, x, "", "", "")
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
	c := Create("basicTest", nil, requestData, "", "", "")
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
	expectedName := "test"
	expectedFallback, _ := json.Marshal(69)

	// Call the Create function
	c := Create(expectedName, nil, expectedFallback, "", "", "")

	// Check if the fields of the created controller match the expected values
	if c.Name != expectedName {
		t.Errorf("Expected name %s, got %s", expectedName, c.Name)
	}

	if !reflect.DeepEqual(c.Fallback, expectedFallback) {
		t.Errorf("Expected fallback %#v, got %#v", expectedFallback, c.Fallback)
	}
}

func TestSetupControllers(t *testing.T) {
	// Initialize Services map
	Services := make(map[string]Controller)

	// Marshal fallback responses
	Fallback1, _ := json.Marshal(69)
	Fallback2, _ := json.Marshal("Hello World")

	// Create controllers and add them to Services map
	Services["/get_int"] = Create("int_controller", nil, Fallback1, "", "", "")
	Services["/get_str"] = Create("str_controller", nil, Fallback2, "", "", "")

	// Setup controllers
	SetupControllers(Services)

	// Test requests
	testCases := []struct {
		endpoint         string
		expectedResponse []byte
	}{
		{"/get_int", Fallback1},
		{"/get_str", Fallback2},
	}

	for _, tc := range testCases {
		t.Run(tc.endpoint, func(t *testing.T) {
			// Create a request
			req := httptest.NewRequest("GET", tc.endpoint, nil)
			rr := httptest.NewRecorder()

			// Handler for the given endpoint
			handler := http.HandlerFunc(Services[tc.endpoint].ServeHTTP)

			// Serve the request
			handler.ServeHTTP(rr, req)

			// Check the response status code
			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}

			// Check the response body
			if !bytes.Equal(rr.Body.Bytes(), tc.expectedResponse) {
				t.Errorf("handler returned unexpected body: got %s want %s", rr.Body.String(), tc.expectedResponse)
			}
		})
	}
}
