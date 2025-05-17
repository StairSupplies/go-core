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

func ExampleNotFoundError() {
	// Helper function for not found errors
	err := errors.New("user not found")
	apiErr := api.NotFoundError(err)

	fmt.Printf("Status: %d, Message: %s\n", apiErr.StatusCode, apiErr.Message)
	// Output: Status: 404, Message: user not found
}

func ExampleServerError() {
	// Helper function for internal server errors
	err := errors.New("database connection failed")
	apiErr := api.ServerError(err)

	fmt.Printf("Status: %d, Message: %s\n", apiErr.StatusCode, apiErr.Message)
	// Output: Status: 500, Message: database connection failed
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
	// Instead of using a real example that prints log output,
	// we'll show the typical usage pattern without executing it.
	fmt.Println("// Create a handler that returns an error")
	fmt.Println("handler := func(w http.ResponseWriter, r *http.Request) error {")
	fmt.Println("    if someCondition {")
	fmt.Println("        return api.BadRequestError(errors.New(\"bad request\"))")
	fmt.Println("    }")
	fmt.Println("    return api.WriteSuccess(w, responseData)")
	fmt.Println("}")
	fmt.Println("")
	fmt.Println("// Wrap the handler to automatically handle errors")
	fmt.Println("wrapped := api.WrapHandler(handler)")
	fmt.Println("")
	fmt.Println("// Register with your router")
	fmt.Println("router.Get(\"/endpoint\", wrapped)")
	fmt.Println("")
	fmt.Println("// The wrapped handler will:")
	fmt.Println("// 1. Call your handler function")
	fmt.Println("// 2. If an error is returned, log it and write an error response")
	fmt.Println("// 3. Otherwise, your handler handles the response writing")
	// Output: // Create a handler that returns an error
	// handler := func(w http.ResponseWriter, r *http.Request) error {
	//     if someCondition {
	//         return api.BadRequestError(errors.New("bad request"))
	//     }
	//     return api.WriteSuccess(w, responseData)
	// }
	// 
	// // Wrap the handler to automatically handle errors
	// wrapped := api.WrapHandler(handler)
	// 
	// // Register with your router
	// router.Get("/endpoint", wrapped)
	// 
	// // The wrapped handler will:
	// // 1. Call your handler function
	// // 2. If an error is returned, log it and write an error response
	// // 3. Otherwise, your handler handles the response writing
}