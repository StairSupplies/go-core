package api_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/StairSupplies/go-core/api"
)

func ExampleEnvelope() {
	// Envelope is a map for wrapping JSON responses in a consistent structure
	envelope := api.Envelope{
		"message": "Hello, world!",
		"count":   42,
		"items":   []string{"apple", "banana", "cherry"},
	}

	// Convert to JSON for display
	data, _ := json.MarshalIndent(envelope, "", "  ")
	fmt.Println(string(data))
	// Output:
	// {
	//   "count": 42,
	//   "items": [
	//     "apple",
	//     "banana",
	//     "cherry"
	//   ],
	//   "message": "Hello, world!"
	// }
}

func ExampleError() {
	// Create an API error
	apiErr := api.Error{
		StatusCode: http.StatusNotFound,
		Message:    "Resource not found",
	}

	// Use the error as a standard Go error
	fmt.Println(apiErr.Error())
	// Output: API Error 404: Resource not found
}

func ExampleNewError() {
	// Create an API error from a standard Go error
	err := errors.New("something went wrong")
	apiErr := api.NewError(http.StatusBadRequest, err)

	fmt.Printf("Status: %d, Message: %s\n", apiErr.StatusCode, apiErr.Message)
	// Output: Status: 400, Message: something went wrong
}

func ExampleBadRequestError() {
	// Helper function for common error types
	err := errors.New("invalid input")
	apiErr := api.BadRequestError(err)

	fmt.Printf("Status: %d, Message: %s\n", apiErr.StatusCode, apiErr.Message)
	// Output: Status: 400, Message: invalid input
}

func ExampleWriteJSON() {
	// Create a test response recorder
	w := httptest.NewRecorder()

	// Write JSON data to the response
	data := map[string]interface{}{
		"name": "John Doe",
		"age":  30,
	}
	headers := http.Header{}
	headers.Set("X-Custom-Header", "custom-value")

	_ = api.WriteJSON(w, http.StatusOK, data, headers)

	// Inspect the response
	fmt.Println("Status Code:", w.Code)
	fmt.Println("Content-Type:", w.Header().Get("Content-Type"))
	fmt.Println("X-Custom-Header:", w.Header().Get("X-Custom-Header"))
	fmt.Println("Body:", w.Body.String())
	// Output:
	// Status Code: 200
	// Content-Type: application/json
	// X-Custom-Header: custom-value
	// Body: {
	//   "age": 30,
	//   "name": "John Doe"
	// }
}

func ExampleWriteSuccess() {
	// Create a test response recorder
	w := httptest.NewRecorder()

	// Write a success response
	data := map[string]string{"message": "Operation successful"}
	meta := map[string]int{"total": 1}

	_ = api.WriteSuccess(w, data, meta)

	// Inspect the response
	fmt.Println("Status Code:", w.Code)
	fmt.Println("Content-Type:", w.Header().Get("Content-Type"))
	fmt.Println("Body:", w.Body.String())
	// Output:
	// Status Code: 200
	// Content-Type: application/json
	// Body: {
	//   "status_code": 200,
	//   "data": {
	//     "message": "Operation successful"
	//   },
	//   "meta": {
	//     "total": 1
	//   }
	// }
}

func ExampleWriteError() {
	// Create a test response recorder
	w := httptest.NewRecorder()

	// Write an error response using NotFoundError helper
	err := api.NotFoundError(errors.New("user not found"))
	api.WriteError(w, err)

	// Inspect the response
	fmt.Println("Status Code:", w.Code)
	fmt.Println("Content-Type:", w.Header().Get("Content-Type"))
	fmt.Println("Body:", w.Body.String())
	// Output:
	// Status Code: 404
	// Content-Type: application/json
	// Body: {
	//   "error": {
	//     "status_code": 404,
	//     "message": "user not found"
	//   }
	// }
}

func ExampleWrapHandler() {
	// WrapHandler converts an api.HandlerFunc (which can return errors)
	// to an http.HandlerFunc (which doesn't return anything)

	// Create a handler that returns an error
	handler := func(w http.ResponseWriter, r *http.Request) error {
		// Check for some condition
		if r.URL.Query().Get("fail") == "true" {
			return api.BadRequestError(errors.New("requested failure"))
		}

		// Otherwise succeed
		return api.WriteSuccess(w, map[string]string{"message": "success"})
	}

	// Wrap the handler to handle errors automatically
	wrapped := api.WrapHandler(handler)

	// Create a test server with the wrapped handler
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// This is just for demonstration - in a real server, you'd register
			// the wrapped handler directly with your router
			wrapped(w, r)
		}))
	defer server.Close()

	// Make a sample request that fails
	req := httptest.NewRequest("GET", "/?fail=true", nil)
	w := httptest.NewRecorder()
	wrapped(w, req)

	// In a real example, we would make an HTTP request to the server
	// But for testing, we'll just check the response recorder
	fmt.Println("Status Code:", w.Code)
	fmt.Println("Content-Type:", w.Header().Get("Content-Type"))
	fmt.Println("Is error response:", w.Body.String()[0:10] == "{\n  \"error")
	// Output:
	// Status Code: 400
	// Content-Type: application/json
	// Is error response: true
}