package rest

import (
	"errors"
	"fmt"
	"net/http"
)

// Common error types for HTTP clients
var (
	// ErrInvalidRequest indicates that the request was invalid.
	ErrInvalidRequest = errors.New("invalid request")

	// ErrUnauthorized indicates that the request was rejected due to authentication failure.
	ErrUnauthorized = errors.New("unauthorized")

	// ErrResourceNotFound indicates that the requested resource could not be found.
	ErrResourceNotFound = errors.New("resource not found")

	// ErrServerError indicates that the server encountered an error processing the request.
	ErrServerError = errors.New("server error")

	// ErrConnectionFailed indicates that the client could not establish a connection.
	ErrConnectionFailed = errors.New("connection failed")

	// ErrTimeout indicates that the request timed out before receiving a response.
	ErrTimeout = errors.New("request timed out")

	// ErrUnprocessableEntity indicates that the server understood the request but was unable to process it.
	ErrUnprocessableEntity = errors.New("unprocessable entity")

	// ErrInvalidResponse indicates that the response from the server was invalid or malformed.
	ErrInvalidResponse = errors.New("invalid response")
)

// ClientError represents an error from the client.
type ClientError struct {
	Err     error  // Underlying error, typically one of the sentinel errors
	Message string // Detailed error message explaining what went wrong
	Code    string // Optional error code, typically derived from the API response
}

// Error returns the error message.
func (e *ClientError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s: %s", e.Err.Error(), e.Message)
	}
	return e.Err.Error()
}

// Unwrap returns the underlying error.
func (e *ClientError) Unwrap() error {
	return e.Err
}

// NewClientError creates a new client error.
func NewClientError(err error, message string, code string) *ClientError {
	return &ClientError{
		Err:     err,
		Message: message,
		Code:    code,
	}
}

// getErrorByStatusCode maps HTTP status codes to error types
func getErrorByStatusCode(statusCode int) error {
	switch {
	case statusCode == http.StatusUnauthorized:
		return ErrUnauthorized
	case statusCode == http.StatusNotFound:
		return ErrResourceNotFound
	case statusCode == http.StatusUnprocessableEntity:
		return ErrUnprocessableEntity
	case statusCode >= 400 && statusCode < 500:
		return ErrInvalidRequest
	case statusCode >= 500:
		return ErrServerError
	default:
		return fmt.Errorf("unexpected status: %s (%d)", http.StatusText(statusCode), statusCode)
	}
}