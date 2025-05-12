package httputils

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

// Envelope is a map for wrapping JSON responses
type Envelope map[string]any

// APIError represents an API error response with status code and message
type APIError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("API Error %d: %s", e.StatusCode, e.Message)
}

// NewAPIError creates a new API error with the given status code and error
func NewAPIError(statusCode int, err error) APIError {
	return APIError{
		StatusCode: statusCode,
		Message:    err.Error(),
	}
}

// ServerError returns a 500 Internal Server Error
func ServerError(err error) APIError {
	return NewAPIError(http.StatusInternalServerError, err)
}

// BadRequestError returns a 400 Bad Request Error
func BadRequestError(err error) APIError {
	return NewAPIError(http.StatusBadRequest, err)
}

// NotFoundError returns a 404 Not Found Error
func NotFoundError(err error) APIError {
	return NewAPIError(http.StatusNotFound, err)
}

// UnauthorizedError returns a 401 Unauthorized Error
func UnauthorizedError(err error) APIError {
	return NewAPIError(http.StatusUnauthorized, err)
}

// ForbiddenError returns a 403 Forbidden Error
func ForbiddenError(err error) APIError {
	return NewAPIError(http.StatusForbidden, err)
}

// UnprocessableEntityError returns a 422 Unprocessable Entity Error
func UnprocessableEntityError(err error) APIError {
	return NewAPIError(http.StatusUnprocessableEntity, err)
}

// SuccessResponse is a standard structure for successful API responses
type SuccessResponse struct {
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data"`
	Meta       interface{} `json:"meta,omitempty"`
}

// WriteJSON writes a JSON response with given status and data
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

// WriteSuccess writes a success response with status 200
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

// WriteError writes an error response
func WriteError(w http.ResponseWriter, err error) {
	var apiErr APIError
	var statusCode int

	// Check if the error is already an APIError
	if e, ok := err.(APIError); ok {
		apiErr = e
		statusCode = e.StatusCode
	} else {
		// Default to internal server error
		apiErr = ServerError(err)
		statusCode = http.StatusInternalServerError
	}

	WriteJSON(w, statusCode, Envelope{"error": apiErr}, nil)
}

// ErrorHandler is a middleware that handles API errors
type APIFunc func(w http.ResponseWriter, r *http.Request) error

// WrapHandler wraps an API function with error handling
func WrapHandler(h APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			// Log the error
			slog.Error("HTTP API Error", "err", err.Error(), "path", r.URL.Path)

			// Write the error response
			WriteError(w, err)
		}
	}
}
