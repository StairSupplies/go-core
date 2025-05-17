package router

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/StairSupplies/go-core/api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func TestRouterCreation(t *testing.T) {
	t.Run("default options", func(t *testing.T) {
		r := New()
		if r == nil {
			t.Fatal("Expected router to be created")
		}

		// Test that middleware are configured
		if !r.options.EnableLogging {
			t.Error("Expected EnableLogging to be true")
		}
		if !r.options.EnableRecovery {
			t.Error("Expected EnableRecovery to be true")
		}
		if !r.options.EnableRequestID {
			t.Error("Expected EnableRequestID to be true")
		}
		if !r.options.EnableTimeout {
			t.Error("Expected EnableTimeout to be true")
		}
		if !r.options.EnableHealthcheck {
			t.Error("Expected EnableHealthcheck to be true")
		}
	})

	t.Run("custom options", func(t *testing.T) {
		customOptions := Options{
			EnableLogging:     false,
			EnableRecovery:    true,
			EnableRequestID:   true,
			EnableTimeout:     false,
			EnableHealthcheck: true,
			TimeoutDuration:   30 * time.Second,
			LoggerOptions: LoggerOptions{
				LogRequestHeaders:  true,
				LogResponseHeaders: true,
				LogRequestBody:     true,
				SkipPaths:          []string{"/skip"},
			},
		}

		customRouter := NewWithOptions(customOptions)

		if customRouter == nil {
			t.Fatal("Expected router with custom options to be created")
		}

		// Verify options were applied
		if customRouter.options.EnableLogging != false {
			t.Error("Expected EnableLogging to be false")
		}
		if customRouter.options.EnableTimeout != false {
			t.Error("Expected EnableTimeout to be false")
		}
		if customRouter.options.TimeoutDuration != 30*time.Second {
			t.Error("Expected TimeoutDuration to be 30 seconds")
		}
		if !customRouter.options.LoggerOptions.LogRequestHeaders {
			t.Error("Expected LogRequestHeaders to be true")
		}
		if len(customRouter.options.LoggerOptions.SkipPaths) != 1 ||
			customRouter.options.LoggerOptions.SkipPaths[0] != "/skip" {
			t.Error("Expected SkipPaths to be ['/skip']")
		}
	})
}

func TestWithErrorHandler(t *testing.T) {
	// Create test handlers
	successHandler := func(w http.ResponseWriter, r *http.Request) error {
		return api.WriteSuccess(w, map[string]string{"message": "success"})
	}

	errorHandler := func(w http.ResponseWriter, r *http.Request) error {
		return api.BadRequestError(errors.New("bad request"))
	}

	panicHandler := func(w http.ResponseWriter, r *http.Request) error {
		panic("test panic")
		return nil
	}

	// Create a router with the handlers
	r := New()
	r.Get("/success", WithErrorHandler(successHandler))
	r.Get("/error", WithErrorHandler(errorHandler))
	r.Get("/panic", WithErrorHandler(panicHandler))

	// Test success handler
	t.Run("success response", func(t *testing.T) {
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

		dataMap, ok := successResp.Data.(map[string]interface{})
		if !ok {
			t.Fatalf("Expected response.Data to be a map")
		}

		if dataMap["message"] != "success" {
			t.Errorf("Expected message to be 'success', got '%v'", dataMap["message"])
		}
	})

	// Test error handler
	t.Run("error response", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/error", nil)
		w := httptest.NewRecorder()
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
	})

	// Test panic handler (middleware should recover)
	t.Run("panic recovery", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/panic", nil)
		w := httptest.NewRecorder()

		// This should not panic due to recovery middleware
		r.ServeHTTP(w, req)

		if w.Code < 500 {
			t.Errorf("Expected 5xx status code for panic, got %d", w.Code)
		}
	})
}

func TestRouterGroup(t *testing.T) {
	r := New()
	apiRouter := r.Group(func(r chi.Router) {
		r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("api test"))
		})

		r.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("user " + id))
		})
	})

	// Mount the api router
	r.Mount("/api", apiRouter)

	// Test the mounted router
	t.Run("basic route", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		body, _ := io.ReadAll(w.Body)
		if string(body) != "api test" {
			t.Errorf("Expected body '%s', got '%s'", "api test", string(body))
		}
	})

	// Test route with URL parameter
	t.Run("route with parameters", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/users/123", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		body, _ := io.ReadAll(w.Body)
		if string(body) != "user 123" {
			t.Errorf("Expected body '%s', got '%s'", "user 123", string(body))
		}
	})
}

func TestWithMiddleware(t *testing.T) {
	// Create a custom middleware
	testMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test", "test-value")
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
	
	if w.Header().Get("X-Test") != "test-value" {
		t.Errorf("Expected header X-Test to be 'test-value', got '%s'", w.Header().Get("X-Test"))
	}
	
	body, _ := io.ReadAll(w.Body)
	if string(body) != "test" {
		t.Errorf("Expected body '%s', got '%s'", "test", string(body))
	}
}

func TestHealthcheck(t *testing.T) {
	// Test with healthcheck enabled (default)
	t.Run("healthcheck enabled", func(t *testing.T) {
		r := New()
		req := httptest.NewRequest("GET", "/healthz", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		body, _ := io.ReadAll(w.Body)
		if string(body) != "OK" {
			t.Errorf("Expected body '%s', got '%s'", "OK", string(body))
		}
	})

	// Test with healthcheck disabled
	t.Run("healthcheck disabled", func(t *testing.T) {
		options := DefaultOptions()
		options.EnableHealthcheck = false
		r := NewWithOptions(options)

		req := httptest.NewRequest("GET", "/healthz", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestMiddlewareConfiguration(t *testing.T) {
	// Test that RequestID middleware is applied
	t.Run("request id middleware", func(t *testing.T) {
		r := New()
		r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			requestID := middleware.GetReqID(r.Context())
			w.Write([]byte(requestID))
		})

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		body, _ := io.ReadAll(w.Body)
		if len(string(body)) == 0 {
			t.Error("Expected Request ID to be present, got empty string")
		}
	})

	// Test that timeout middleware configuration is applied
	t.Run("timeout middleware configuration", func(t *testing.T) {
		options := DefaultOptions()
		options.TimeoutDuration = 50 * time.Millisecond
		r := NewWithOptions(options)

		// Verify the option was set correctly
		if r.options.TimeoutDuration != 50*time.Millisecond {
			t.Errorf("Expected TimeoutDuration to be 50ms, got %v", r.options.TimeoutDuration)
		}

		// Add a simple handler
		r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("test"))
		})

		// Test that the handler works
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
	})
}

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()

	// Verify default values
	if !opts.EnableLogging {
		t.Error("Expected EnableLogging to be true by default")
	}
	if !opts.EnableRecovery {
		t.Error("Expected EnableRecovery to be true by default")
	}
	if !opts.EnableRequestID {
		t.Error("Expected EnableRequestID to be true by default")
	}
	if !opts.EnableTimeout {
		t.Error("Expected EnableTimeout to be true by default")
	}
	if !opts.EnableHealthcheck {
		t.Error("Expected EnableHealthcheck to be true by default")
	}
	if opts.TimeoutDuration != 60*time.Second {
		t.Errorf("Expected TimeoutDuration to be 60s, got %v", opts.TimeoutDuration)
	}
}

func TestServeHTTP(t *testing.T) {
	r := New()
	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test"))
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Test that the ServeHTTP method works correctly
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	body, _ := io.ReadAll(w.Body)
	if string(body) != "test" {
		t.Errorf("Expected body '%s', got '%s'", "test", string(body))
	}
}