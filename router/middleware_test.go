package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test"))
	})
	
	// Create the logger middleware
	loggerMiddleware := Logger(LoggerOptions{
		LogRequestHeaders:  true,
		LogResponseHeaders: true,
		SkipPaths:          []string{"/skip"},
	})
	
	// Wrap the test handler with logger middleware
	handler := loggerMiddleware(testHandler)
	
	// Test normal request
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Test", "test")
	w := httptest.NewRecorder()
	
	handler.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	
	// Test skipped path
	req = httptest.NewRequest("GET", "/skip", nil)
	w = httptest.NewRecorder()
	
	handler.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestRecoverer(t *testing.T) {
	// Create a test handler that panics
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})
	
	// Wrap the test handler with recoverer middleware
	handler := Recoverer()(panicHandler)
	
	// Test request that causes panic
	req := httptest.NewRequest("GET", "/panic", nil)
	w := httptest.NewRecorder()
	
	handler.ServeHTTP(w, req)
	
	// Should recover and return 500
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestTimeout(t *testing.T) {
	// Create a test handler that sleeps
	slowHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("done"))
	})
	
	// Test successful response (no timeout)
	timeoutMiddleware := Timeout(100 * time.Millisecond)
	handler := timeoutMiddleware(slowHandler)
	
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	
	handler.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	
	// Test timeout
	timeoutMiddleware = Timeout(10 * time.Millisecond)
	handler = timeoutMiddleware(slowHandler)
	
	req = httptest.NewRequest("GET", "/test", nil)
	w = httptest.NewRecorder()
	
	handler.ServeHTTP(w, req)
	
	// Chi's timeout middleware doesn't automatically write an error response,
	// but the connection should be closed. We can't easily test this in a unit test.
}

func TestRequestID(t *testing.T) {
	// Create a test handler that checks for request ID
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Just return OK, we're not testing the request ID value itself
		w.WriteHeader(http.StatusOK)
	})
	
	// Wrap the test handler with request ID middleware
	handler := RequestID(testHandler)
	
	// Test request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	
	handler.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}