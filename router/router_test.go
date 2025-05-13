package router

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/StairSupplies/go-core/api"
	"github.com/go-chi/chi/v5"
)

func TestRouterCreation(t *testing.T) {
	r := New()
	if r == nil {
		t.Fatal("Expected router to be created")
	}

	// Test with custom options
	customRouter := NewWithOptions(Options{
		EnableLogging:   false,
		EnableRecovery:  true,
		EnableRequestID: true,
		EnableTimeout:   false,
	})

	if customRouter == nil {
		t.Fatal("Expected router with custom options to be created")
	}
}

func TestWithErrorHandler(t *testing.T) {
	// Create test handlers
	successHandler := func(w http.ResponseWriter, r *http.Request) error {
		return api.WriteSuccess(w, map[string]string{"message": "success"})
	}

	errorHandler := func(w http.ResponseWriter, r *http.Request) error {
		return api.BadRequestError(errors.New("bad request"))
	}

	// Create a router with the handlers
	r := New()
	r.Get("/success", WithErrorHandler(successHandler))
	r.Get("/error", WithErrorHandler(errorHandler))

	// Test success handler
	req := httptest.NewRequest("GET", "/success", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Verify response structure
	var successResp api.SuccessResponse
	if err := json.Unmarshal(w.Body.Bytes(), &successResp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if successResp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d in response, got %d", http.StatusOK, successResp.StatusCode)
	}

	// Test error handler
	req = httptest.NewRequest("GET", "/error", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Verify error response structure
	var errorResp map[string]api.Error
	if err := json.Unmarshal(w.Body.Bytes(), &errorResp); err != nil {
		t.Fatalf("Failed to unmarshal error response: %v", err)
	}

	if apiErr, ok := errorResp["error"]; ok {
		if apiErr.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status code %d in error response, got %d", http.StatusBadRequest, apiErr.StatusCode)
		}
		if apiErr.Message != "bad request" {
			t.Errorf("Expected error message '%s', got '%s'", "bad request", apiErr.Message)
		}
	} else {
		t.Errorf("Expected error field in response, got %v", errorResp)
	}
}

func TestRouterGroup(t *testing.T) {
	r := New()
	apiRouter := r.Group(func(r chi.Router) {
		r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("test"))
		})
	})

	// Mount the api router
	r.Mount("/api", apiRouter)

	// Test the mounted router
	req := httptest.NewRequest("GET", "/api/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	body, _ := io.ReadAll(w.Body)
	if string(body) != "test" {
		t.Errorf("Expected body '%s', got '%s'", "test", string(body))
	}
}

func TestWithMiddleware(t *testing.T) {
	// Create a custom middleware
	testMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test", "test")
			next.ServeHTTP(w, r)
		})
	}
	
	// Create router with middleware first
	r := New().WithMiddleware(testMiddleware)
	
	// Then add a test handler
	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test"))
	})
	
	// Test the handler with middleware
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	
	if w.Header().Get("X-Test") != "test" {
		t.Errorf("Expected header X-Test to be 'test', got '%s'", w.Header().Get("X-Test"))
	}
}