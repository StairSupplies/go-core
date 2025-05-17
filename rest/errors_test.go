package rest

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestClientError_Error(t *testing.T) {
	// Test with message
	t.Run("with message", func(t *testing.T) {
		err := &ClientError{
			Err:     ErrNotFound,
			Message: "resource not found",
		}

		want := "resource not found: resource not found"
		got := err.Error()

		if got != want {
			t.Errorf("ClientError.Error() = %q, want %q", got, want)
		}
	})

	// Test without message
	t.Run("without message", func(t *testing.T) {
		err := &ClientError{
			Err: ErrServerError,
		}

		want := "server error"
		got := err.Error()

		if got != want {
			t.Errorf("ClientError.Error() = %q, want %q", got, want)
		}
	})
}

func TestClientError_Unwrap(t *testing.T) {
	underlying := ErrNotFound
	err := &ClientError{
		Err:     underlying,
		Message: "test message",
	}

	// Test unwrap
	unwrapped := err.Unwrap()
	if unwrapped != underlying {
		t.Errorf("ClientError.Unwrap() = %v, want %v", unwrapped, underlying)
	}

	// Test with errors.Is
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("errors.Is() failed to find the underlying error")
	}
}

func TestNewClientError(t *testing.T) {
	err := NewClientError(ErrTimeout, "test message", "TIMEOUT")

	if err.Err != ErrTimeout {
		t.Errorf("Expected Err to be ErrTimeout, got %v", err.Err)
	}

	if err.Message != "test message" {
		t.Errorf("Expected Message to be 'test message', got %q", err.Message)
	}

	if err.Code != "TIMEOUT" {
		t.Errorf("Expected Code to be 'TIMEOUT', got %q", err.Code)
	}
}

func TestGetErrorByStatusCode(t *testing.T) {
	tests := []struct {
		statusCode int
		expected   error
	}{
		{http.StatusUnauthorized, ErrUnauthorized},
		{http.StatusNotFound, ErrNotFound},
		{http.StatusUnprocessableEntity, ErrUnprocessableEntity},
		{http.StatusBadRequest, ErrInvalidRequest},
		{http.StatusConflict, ErrInvalidRequest},
		{http.StatusInternalServerError, ErrServerError},
		{http.StatusServiceUnavailable, ErrServerError},
		{0, fmt.Errorf("unexpected status: %s (%d)", http.StatusText(0), 0)},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("status_%d", tt.statusCode), func(t *testing.T) {
			got := getErrorByStatusCode(tt.statusCode)

			// Handle special case for unexpected status
			if tt.statusCode == 0 {
				if got == nil || got.Error() != tt.expected.Error() {
					t.Errorf("getErrorByStatusCode(%d) = %v, want %v", tt.statusCode, got, tt.expected)
				}
				return
			}

			if got != tt.expected {
				t.Errorf("getErrorByStatusCode(%d) = %v, want %v", tt.statusCode, got, tt.expected)
			}
		})
	}
}

func TestSentinelErrors(t *testing.T) {
	// Ensure all sentinel errors exist and have appropriate messages
	tests := []struct {
		err  error
		want string
	}{
		{ErrInvalidRequest, "invalid request"},
		{ErrUnauthorized, "unauthorized"},
		{ErrNotFound, "resource not found"},
		{ErrServerError, "server error"},
		{ErrConnectionFailed, "connection failed"},
		{ErrTimeout, "request timed out"},
		{ErrUnprocessableEntity, "unprocessable entity"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if tt.err == nil {
				t.Errorf("Expected %s error to be defined", tt.want)
				return
			}

			if tt.err.Error() != tt.want {
				t.Errorf("Error message = %q, want %q", tt.err.Error(), tt.want)
			}
		})
	}
}