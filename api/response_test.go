package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestError_Error(t *testing.T) {
	tests := []struct {
		name     string
		apiErr   Error
		expected string
	}{
		{
			name: "bad request error",
			apiErr: Error{
				StatusCode: http.StatusBadRequest,
				Message:    "invalid request",
			},
			expected: "API Error 400: invalid request",
		},
		{
			name: "not found error",
			apiErr: Error{
				StatusCode: http.StatusNotFound,
				Message:    "resource not found",
			},
			expected: "API Error 404: resource not found",
		},
		{
			name: "server error",
			apiErr: Error{
				StatusCode: http.StatusInternalServerError,
				Message:    "internal server error",
			},
			expected: "API Error 500: internal server error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.apiErr.Error(); got != tc.expected {
				t.Errorf("Error.Error() = %v, want %v", got, tc.expected)
			}
		})
	}
}

func TestNewError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		err        error
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "basic error",
			statusCode: http.StatusBadRequest,
			err:        errors.New("test error"),
			wantStatus: http.StatusBadRequest,
			wantMsg:    "test error",
		},
		{
			name:       "empty error",
			statusCode: http.StatusNotFound,
			err:        errors.New(""),
			wantStatus: http.StatusNotFound,
			wantMsg:    "",
		},
		{
			name:       "wrapped error",
			statusCode: http.StatusInternalServerError,
			err:        fmt.Errorf("wrapped: %w", errors.New("original")),
			wantStatus: http.StatusInternalServerError,
			wantMsg:    "wrapped: original",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			apiErr := NewError(tc.statusCode, tc.err)

			if apiErr.StatusCode != tc.wantStatus {
				t.Errorf("Expected StatusCode to be %d, got %d", tc.wantStatus, apiErr.StatusCode)
			}
			if apiErr.Message != tc.wantMsg {
				t.Errorf("Expected Message to be '%s', got '%s'", tc.wantMsg, apiErr.Message)
			}
		})
	}
}

func TestErrorFactories(t *testing.T) {
	err := errors.New("test error")

	tests := []struct {
		name     string
		factory  func(error) Error
		expected int
	}{
		{"ServerError", ServerError, http.StatusInternalServerError},
		{"BadRequestError", BadRequestError, http.StatusBadRequest},
		{"NotFoundError", NotFoundError, http.StatusNotFound},
		{"UnauthorizedError", UnauthorizedError, http.StatusUnauthorized},
		{"ForbiddenError", ForbiddenError, http.StatusForbidden},
		{"UnprocessableEntityError", UnprocessableEntityError, http.StatusUnprocessableEntity},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			apiErr := tc.factory(err)

			if apiErr.StatusCode != tc.expected {
				t.Errorf("Expected StatusCode to be %d, got %d", tc.expected, apiErr.StatusCode)
			}
			if apiErr.Message != "test error" {
				t.Errorf("Expected Message to be 'test error', got '%s'", apiErr.Message)
			}
		})
	}
}

func TestEnvelope(t *testing.T) {
	// Test creating and using an envelope
	envelope := Envelope{
		"message": "Hello",
		"count":   5,
		"items":   []string{"apple", "banana"},
	}

	// Make sure we can access fields
	if msg, ok := envelope["message"].(string); !ok || msg != "Hello" {
		t.Errorf("Expected envelope[\"message\"] to be \"Hello\", got %v", envelope["message"])
	}

	if count, ok := envelope["count"].(int); !ok || count != 5 {
		t.Errorf("Expected envelope[\"count\"] to be 5, got %v", envelope["count"])
	}

	// Test serialization
	data, err := json.Marshal(envelope)
	if err != nil {
		t.Fatalf("Failed to marshal envelope: %v", err)
	}

	// Unmarshal back and compare
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal envelope: %v", err)
	}

	if result["message"] != "Hello" {
		t.Errorf("Expected result[\"message\"] to be \"Hello\", got %v", result["message"])
	}
}

func TestWriteJSON(t *testing.T) {
	tests := []struct {
		name           string
		data           interface{}
		headers        http.Header
		status         int
		expectedStatus int
		expectedType   string
		expectedBody   string
	}{
		{
			name:           "simple data",
			data:           map[string]string{"key": "value"},
			headers:        nil,
			status:         http.StatusOK,
			expectedStatus: http.StatusOK,
			expectedType:   "application/json",
			expectedBody:   "{\n  \"key\": \"value\"\n}\n",
		},
		{
			name:           "with custom headers",
			data:           map[string]int{"count": 42},
			headers:        http.Header{"X-Custom-Header": []string{"test"}},
			status:         http.StatusCreated,
			expectedStatus: http.StatusCreated,
			expectedType:   "application/json",
			expectedBody:   "{\n  \"count\": 42\n}\n",
		},
		{
			name:           "with array data",
			data:           []string{"apple", "banana"},
			headers:        nil,
			status:         http.StatusOK,
			expectedStatus: http.StatusOK,
			expectedType:   "application/json",
			expectedBody:   "[\n  \"apple\",\n  \"banana\"\n]\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			err := WriteJSON(rr, tc.status, tc.data, tc.headers)
			if err != nil {
				t.Fatalf("WriteJSON returned an error: %v", err)
			}

			// Check status code
			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, rr.Code)
			}

			// Check content type
			if rr.Header().Get("Content-Type") != tc.expectedType {
				t.Errorf("Expected Content-Type to be '%s', got '%s'", tc.expectedType, rr.Header().Get("Content-Type"))
			}

			// Check custom headers
			if tc.headers != nil {
				for key, values := range tc.headers {
					if !reflect.DeepEqual(rr.Header()[key], values) {
						t.Errorf("Expected header %s to be %v, got %v", key, values, rr.Header()[key])
					}
				}
			}

			// Check response body
			if rr.Body.String() != tc.expectedBody {
				t.Errorf("Expected body to be:\n%s\nGot:\n%s", tc.expectedBody, rr.Body.String())
			}
		})
	}
}

func TestWriteSuccess(t *testing.T) {
	tests := []struct {
		name        string
		data        interface{}
		meta        interface{}
		expectedKey string
	}{
		{
			name:        "with data only",
			data:        map[string]string{"key": "value"},
			meta:        nil,
			expectedKey: "value",
		},
		{
			name:        "with data and meta",
			data:        map[string]string{"key": "value"},
			meta:        map[string]int{"count": 1},
			expectedKey: "value",
		},
		{
			name:        "with array data",
			data:        []string{"apple", "banana"},
			meta:        nil,
			expectedKey: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			var err error
			if tc.meta == nil {
				err = WriteSuccess(rr, tc.data)
			} else {
				err = WriteSuccess(rr, tc.data, tc.meta)
			}

			if err != nil {
				t.Fatalf("WriteSuccess returned an error: %v", err)
			}

			// Check response status code
			if rr.Code != http.StatusOK {
				t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
			}

			// Check content type
			if rr.Header().Get("Content-Type") != "application/json" {
				t.Errorf("Expected Content-Type to be 'application/json', got '%s'", rr.Header().Get("Content-Type"))
			}

			// Parse the response
			var response SuccessResponse
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			// Check status code in response
			if response.StatusCode != http.StatusOK {
				t.Errorf("Expected StatusCode in response to be %d, got %d", http.StatusOK, response.StatusCode)
			}

			// Check data structure based on test case
			switch typedData := tc.data.(type) {
			case map[string]string:
				dataMap, ok := response.Data.(map[string]interface{})
				if !ok {
					t.Fatalf("Failed to convert response.Data to map")
				}

				if key, exists := typedData["key"]; exists && tc.expectedKey != "" {
					if dataMap["key"] != key {
						t.Errorf("Expected data['key'] to be '%s', got '%v'", key, dataMap["key"])
					}
				}
			case []string:
				dataArray, ok := response.Data.([]interface{})
				if !ok {
					t.Fatalf("Failed to convert response.Data to array")
				}

				if len(dataArray) != len(typedData) {
					t.Errorf("Expected data array length to be %d, got %d", len(typedData), len(dataArray))
				}
			}

			// Check meta if provided
			if tc.meta != nil {
				if response.Meta == nil {
					t.Fatalf("Expected meta data but found none")
				}

				metaMap, ok := response.Meta.(map[string]interface{})
				if !ok {
					t.Fatalf("Failed to convert response.Meta to map")
				}

				countMeta, ok := tc.meta.(map[string]int)
				if ok && countMeta["count"] != int(metaMap["count"].(float64)) {
					t.Errorf("Expected meta['count'] to be %d, got %v", countMeta["count"], metaMap["count"])
				}
			} else if response.Meta != nil {
				t.Errorf("Expected meta to be nil, got %v", response.Meta)
			}
		})
	}
}

func TestWriteError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:           "regular error",
			err:            errors.New("test error"),
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "test error",
		},
		{
			name:           "api error",
			err:            BadRequestError(errors.New("invalid request")),
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "invalid request",
		},
		{
			name:           "not found error",
			err:            NotFoundError(errors.New("resource not found")),
			expectedStatus: http.StatusNotFound,
			expectedMsg:    "resource not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			WriteError(rr, tc.err)

			// Check status code
			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, rr.Code)
			}

			// Check content type
			if rr.Header().Get("Content-Type") != "application/json" {
				t.Errorf("Expected Content-Type to be 'application/json', got '%s'", rr.Header().Get("Content-Type"))
			}

			// Parse response
			var response map[string]interface{}
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			// Check error in response
			errorData, ok := response["error"].(map[string]interface{})
			if !ok {
				t.Fatalf("Failed to convert response['error'] to map")
			}

			if int(errorData["status_code"].(float64)) != tc.expectedStatus {
				t.Errorf("Expected error.status_code to be %d, got %v", tc.expectedStatus, errorData["status_code"])
			}

			if errorData["message"] != tc.expectedMsg {
				t.Errorf("Expected error.message to be '%s', got '%v'", tc.expectedMsg, errorData["message"])
			}
		})
	}
}

func TestWrapHandler(t *testing.T) {
	tests := []struct {
		name           string
		handler        func(w http.ResponseWriter, r *http.Request) error
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "no error",
			handler: func(w http.ResponseWriter, r *http.Request) error {
				w.Write([]byte("success"))
				return nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
		{
			name: "with bad request error",
			handler: func(w http.ResponseWriter, r *http.Request) error {
				return BadRequestError(errors.New("invalid request"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "",  // Body checked separately for error responses
		},
		{
			name: "with server error",
			handler: func(w http.ResponseWriter, r *http.Request) error {
				return errors.New("unexpected error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "", // Body checked separately for error responses
		},
		{
			name: "with custom response",
			handler: func(w http.ResponseWriter, r *http.Request) error {
				return WriteJSON(w, http.StatusCreated, map[string]string{"status": "created"}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   "{\n  \"status\": \"created\"\n}\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/", nil)

			handler := WrapHandler(tc.handler)
			handler(rr, req)

			// Check status code
			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, rr.Code)
			}

			// For regular responses, check body directly
			if tc.expectedBody != "" {
				body, _ := io.ReadAll(rr.Body)
				if string(body) != tc.expectedBody {
					t.Errorf("Expected response body to be '%s', got '%s'", tc.expectedBody, string(body))
				}
			}

			// For error responses, verify JSON structure
			if tc.expectedStatus >= 400 {
				var response map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("Failed to unmarshal error response: %v", err)
				}

				errorData, ok := response["error"].(map[string]interface{})
				if !ok {
					t.Fatalf("Failed to convert response['error'] to map")
				}

				if int(errorData["status_code"].(float64)) != tc.expectedStatus {
					t.Errorf("Expected error.status_code to be %d, got %v", tc.expectedStatus, errorData["status_code"])
				}
			}
		})
	}
}