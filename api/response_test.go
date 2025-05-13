package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestError_Error(t *testing.T) {
	// Create API error
	apiErr := Error{
		StatusCode: http.StatusBadRequest,
		Message:    "invalid request",
	}

	// Test Error() method
	expected := "API Error 400: invalid request"
	if got := apiErr.Error(); got != expected {
		t.Errorf("Error.Error() = %v, want %v", got, expected)
	}
}

func TestNewError(t *testing.T) {
	// Create error
	err := errors.New("test error")

	// Create API error
	apiErr := NewError(http.StatusBadRequest, err)

	// Validate result
	if apiErr.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected StatusCode to be %d, got %d", http.StatusBadRequest, apiErr.StatusCode)
	}
	if apiErr.Message != "test error" {
		t.Errorf("Expected Message to be 'test error', got '%s'", apiErr.Message)
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

func TestWriteJSON(t *testing.T) {
	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create test data
	data := map[string]string{"key": "value"}

	// Create custom headers
	headers := http.Header{}
	headers.Set("X-Custom-Header", "test")

	// Call function
	err := WriteJSON(rr, http.StatusOK, data, headers)
	if err != nil {
		t.Fatalf("WriteJSON returned an error: %v", err)
	}

	// Check response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// Check headers
	if rr.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type to be 'application/json', got '%s'", rr.Header().Get("Content-Type"))
	}
	if rr.Header().Get("X-Custom-Header") != "test" {
		t.Errorf("Expected X-Custom-Header to be 'test', got '%s'", rr.Header().Get("X-Custom-Header"))
	}

	// Check response body
	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["key"] != "value" {
		t.Errorf("Expected response['key'] to be 'value', got '%s'", response["key"])
	}
}

func TestWriteSuccess(t *testing.T) {
	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create test data
	data := map[string]string{"key": "value"}
	meta := map[string]int{"count": 1}

	// Call function
	err := WriteSuccess(rr, data, meta)
	if err != nil {
		t.Fatalf("WriteSuccess returned an error: %v", err)
	}

	// Check response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// Check response body
	var response SuccessResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Check status code in response
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode to be %d, got %d", http.StatusOK, response.StatusCode)
	}

	// Check data in response
	dataMap, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to convert response.Data to map")
	}

	if dataMap["key"] != "value" {
		t.Errorf("Expected data['key'] to be 'value', got '%v'", dataMap["key"])
	}

	// Check meta in response
	metaMap, ok := response.Meta.(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to convert response.Meta to map")
	}

	if metaMap["count"] != float64(1) {
		t.Errorf("Expected meta['count'] to be 1, got '%v'", metaMap["count"])
	}
}

func TestWriteError(t *testing.T) {
	// Test with regular error
	t.Run("regular error", func(t *testing.T) {
		rr := httptest.NewRecorder()
		err := errors.New("test error")

		WriteError(rr, err)

		// Check status code (should be 500 for regular error)
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, rr.Code)
		}

		// Parse response
		var response map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		// Check error in response
		errorData, ok := response["error"].(map[string]interface{})
		if !ok {
			t.Fatalf("Failed to convert response['error'] to map")
		}

		if errorData["status_code"].(float64) != float64(http.StatusInternalServerError) {
			t.Errorf("Expected error.status_code to be %d, got %v",
				http.StatusInternalServerError, errorData["status_code"])
		}

		if errorData["message"] != "test error" {
			t.Errorf("Expected error.message to be 'test error', got '%v'", errorData["message"])
		}
	})

	// Test with API Error
	t.Run("API Error", func(t *testing.T) {
		rr := httptest.NewRecorder()
		err := Error{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid request",
		}

		WriteError(rr, err)

		// Check status code
		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}

		// Parse response
		var response map[string]interface{}
		err2 := json.Unmarshal(rr.Body.Bytes(), &response)
		if err2 != nil {
			t.Fatalf("Failed to unmarshal response: %v", err2)
		}

		// Check error in response
		errorData, ok := response["error"].(map[string]interface{})
		if !ok {
			t.Fatalf("Failed to convert response['error'] to map")
		}

		if errorData["status_code"].(float64) != float64(http.StatusBadRequest) {
			t.Errorf("Expected error.status_code to be %d, got %v",
				http.StatusBadRequest, errorData["status_code"])
		}

		if errorData["message"] != "invalid request" {
			t.Errorf("Expected error.message to be 'invalid request', got '%v'", errorData["message"])
		}
	})
}

func TestWrapHandler(t *testing.T) {
	// Test with handler that doesn't return error
	t.Run("no error", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		handler := WrapHandler(func(w http.ResponseWriter, r *http.Request) error {
			w.Write([]byte("success"))
			return nil
		})

		handler(rr, req)

		// Check response body
		body, _ := io.ReadAll(rr.Body)
		if string(body) != "success" {
			t.Errorf("Expected response body to be 'success', got '%s'", string(body))
		}
	})

	// Test with handler that returns error
	t.Run("with error", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		handler := WrapHandler(func(w http.ResponseWriter, r *http.Request) error {
			return BadRequestError(errors.New("invalid request"))
		})

		handler(rr, req)

		// Check status code
		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
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

		if errorData["message"] != "invalid request" {
			t.Errorf("Expected error.message to be 'invalid request', got '%v'", errorData["message"])
		}
	})
}