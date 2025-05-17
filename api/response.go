package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/StairSupplies/go-core/logger"
	"go.uber.org/zap"
)

// Envelope is a map for wrapping JSON responses in a consistent structure.
// It allows for flexible response structures while maintaining a standardized format.
type Envelope map[string]any

// Error represents an API error response with status code and message.
// It implements the error interface for seamless integration with Go's error handling.
type Error struct {
	StatusCode int    `json:"status_code"` // HTTP status code
	Message    string `json:"message"`     // Human-readable error message
}

// Error implements the error interface, returning a formatted error message.
func (e Error) Error() string {
	return fmt.Sprintf("API Error %d: %s", e.StatusCode, e.Message)
}

// NewError creates a new API error with the given status code and error.
// It extracts the error message from the provided error.
func NewError(statusCode int, err error) Error {
	return Error{
		StatusCode: statusCode,
		Message:    err.Error(),
	}
}

// ServerError returns a 500 Internal Server Error.
// Use this for unexpected errors that are not the client's fault.
func ServerError(err error) Error {
	return NewError(http.StatusInternalServerError, err)
}

// BadRequestError returns a 400 Bad Request Error.
// Use this when the client sends an invalid request (e.g., validation failures).
func BadRequestError(err error) Error {
	return NewError(http.StatusBadRequest, err)
}

// NotFoundError returns a 404 Not Found Error.
// Use this when the requested resource does not exist.
func NotFoundError(err error) Error {
	return NewError(http.StatusNotFound, err)
}

// UnauthorizedError returns a 401 Unauthorized Error.
// Use this when authentication is required but not provided or invalid.
func UnauthorizedError(err error) Error {
	return NewError(http.StatusUnauthorized, err)
}

// ForbiddenError returns a 403 Forbidden Error.
// Use this when the client is authenticated but does not have permission.
func ForbiddenError(err error) Error {
	return NewError(http.StatusForbidden, err)
}

// UnprocessableEntityError returns a 422 Unprocessable Entity Error.
// Use this when the request is well-formed but cannot be processed due to semantic errors.
func UnprocessableEntityError(err error) Error {
	return NewError(http.StatusUnprocessableEntity, err)
}

// SuccessResponse is a standard structure for successful API responses.
// It provides a consistent response format with status code, data, and optional metadata.
type SuccessResponse struct {
	StatusCode int         `json:"status_code"`    // HTTP status code
	Data       interface{} `json:"data"`           // Response payload
	Meta       interface{} `json:"meta,omitempty"` // Optional metadata (pagination, counts, etc.)
}

// WriteJSON writes a JSON response with the given status and data.
// It handles JSON serialization, content-type headers, and status code setting.
// Additional headers can be provided to be included in the response.
func WriteJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

// WriteSuccess writes a success response with status 200.
// It wraps the provided data in a SuccessResponse structure.
// Optional metadata can be provided as the final parameter.
func WriteSuccess(w http.ResponseWriter, data any, meta ...any) error {
	var metaData interface{}
	if len(meta) > 0 {
		metaData = meta[0]
	}

	resp := SuccessResponse{
		StatusCode: http.StatusOK,
		Data:       data,
		Meta:       metaData,
	}

	return WriteJSON(w, http.StatusOK, resp, nil)
}

// WriteError writes an error response.
// It handles both api.Error instances and standard Go errors.
// Standard errors are converted to 500 Internal Server Error responses.
func WriteError(w http.ResponseWriter, err error) {
	var apiErr Error
	var statusCode int

	// Check if the error is already an API Error
	if e, ok := err.(Error); ok {
		apiErr = e
		statusCode = e.StatusCode
	} else {
		// Default to internal server error
		apiErr = ServerError(err)
		statusCode = http.StatusInternalServerError
	}

	WriteJSON(w, statusCode, Envelope{"error": apiErr}, nil)
}

// HandlerFunc is a function that handles an HTTP request and may return an error.
// This pattern allows for cleaner controller logic by separating error handling.
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// WrapHandler wraps an API function with error handling.
// It automatically handles logging errors and writing error responses.
// This allows controller functions to focus on business logic and simply return errors.
func WrapHandler(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			// Get logger from context or use default
			log := logger.WithContext(r.Context())

			// Log the error
			log.Error("HTTP API Error",
				zap.String("path", r.URL.Path),
				zap.String("err", err.Error()),
			)

			// Write the error response
			WriteError(w, err)
		}
	}
}
