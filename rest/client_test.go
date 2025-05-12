package rest

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	
	// Check default values
	if client.HTTPClient.Timeout != 30*time.Second {
		t.Errorf("Expected default timeout to be 30s, got %v", client.HTTPClient.Timeout)
	}
	
	if client.BaseURL != "" {
		t.Errorf("Expected default BaseURL to be empty, got '%s'", client.BaseURL)
	}
	
	if len(client.Headers) != 0 {
		t.Errorf("Expected default Headers to be empty, got %v", client.Headers)
	}
}

func TestClientWithBaseURL(t *testing.T) {
	client := NewClient().WithBaseURL("https://api.example.com")
	
	if client.BaseURL != "https://api.example.com" {
		t.Errorf("Expected BaseURL to be 'https://api.example.com', got '%s'", client.BaseURL)
	}
}

func TestClientWithHeader(t *testing.T) {
	client := NewClient().WithHeader("X-API-Key", "test-key")
	
	if client.Headers["X-API-Key"] != "test-key" {
		t.Errorf("Expected Headers to contain X-API-Key with value 'test-key'")
	}
}

func TestClientWithTimeout(t *testing.T) {
	client := NewClient().WithTimeout(5 * time.Second)
	
	if client.DefaultTimeout != 5*time.Second {
		t.Errorf("Expected DefaultTimeout to be 5s, got %v", client.DefaultTimeout)
	}
	
	if client.HTTPClient.Timeout != 5*time.Second {
		t.Errorf("Expected HTTPClient.Timeout to be 5s, got %v", client.HTTPClient.Timeout)
	}
}

func TestRequest(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method
		if r.Method != http.MethodPost {
			t.Errorf("Expected request method to be POST, got %s", r.Method)
		}
		
		// Check content type header
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type header to be 'application/json', got '%s'", r.Header.Get("Content-Type"))
		}
		
		// Check custom header
		if r.Header.Get("X-API-Key") != "test-key" {
			t.Errorf("Expected X-API-Key header to be 'test-key', got '%s'", r.Header.Get("X-API-Key"))
		}
		
		// Read request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}
		
		// Parse request body
		var requestData map[string]interface{}
		err = json.Unmarshal(body, &requestData)
		if err != nil {
			t.Fatalf("Failed to parse request body: %v", err)
		}
		
		// Check request data
		if requestData["name"] != "test" {
			t.Errorf("Expected request data name to be 'test', got '%v'", requestData["name"])
		}
		
		// Write response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": 123, "name": "test"}`))
	}))
	defer server.Close()
	
	// Create client
	client := NewClient().
		WithBaseURL(server.URL).
		WithHeader("X-API-Key", "test-key")
	
	// Define request and response types
	type RequestData struct {
		Name string `json:"name"`
	}
	
	type ResponseData struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	
	// Make request
	requestData := RequestData{Name: "test"}
	var responseData ResponseData
	
	err := client.Request(context.Background(), http.MethodPost, "/test", requestData, &responseData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	// Check response data
	if responseData.ID != 123 {
		t.Errorf("Expected response ID to be 123, got %d", responseData.ID)
	}
	if responseData.Name != "test" {
		t.Errorf("Expected response Name to be 'test', got '%s'", responseData.Name)
	}
}

func TestRequestWithNoBaseURL(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": 123, "name": "test"}`))
	}))
	defer server.Close()
	
	// Create client without base URL
	client := NewClient()
	
	// Define response type
	type ResponseData struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	
	// Make request with full URL
	var responseData ResponseData
	err := client.Request(context.Background(), http.MethodGet, server.URL, nil, &responseData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	// Check response data
	if responseData.ID != 123 {
		t.Errorf("Expected response ID to be 123, got %d", responseData.ID)
	}
}

func TestRequestWithNilResponse(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": 123, "name": "test"}`))
	}))
	defer server.Close()
	
	// Create client
	client := NewClient().WithBaseURL(server.URL)
	
	// Make request with nil response
	err := client.Request(context.Background(), http.MethodGet, "/test", nil, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestRequestWithNonSuccessStatusCode(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Bad Request"}`))
	}))
	defer server.Close()
	
	// Create client
	client := NewClient().WithBaseURL(server.URL)
	
	// Define response type
	type ResponseData struct {
		ID int `json:"id"`
	}
	
	// Make request
	var responseData ResponseData
	err := client.Request(context.Background(), http.MethodGet, "/test", nil, &responseData)
	
	// Check error
	if err == nil {
		t.Fatal("Expected error for non-success status code, got nil")
	}
	
	if !contains(err.Error(), "HTTP error: 400") {
		t.Errorf("Expected error message to contain 'HTTP error: 400', got '%s'", err.Error())
	}
}

func TestHttpMethods(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method
		switch r.URL.Path {
		case "/get":
			if r.Method != http.MethodGet {
				t.Errorf("Expected request method to be GET, got %s", r.Method)
			}
		case "/post":
			if r.Method != http.MethodPost {
				t.Errorf("Expected request method to be POST, got %s", r.Method)
			}
		case "/put":
			if r.Method != http.MethodPut {
				t.Errorf("Expected request method to be PUT, got %s", r.Method)
			}
		case "/patch":
			if r.Method != http.MethodPatch {
				t.Errorf("Expected request method to be PATCH, got %s", r.Method)
			}
		case "/delete":
			if r.Method != http.MethodDelete {
				t.Errorf("Expected request method to be DELETE, got %s", r.Method)
			}
		}
		
		// Write response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"method": "` + r.Method + `"}`))
	}))
	defer server.Close()
	
	// Create client
	client := NewClient().WithBaseURL(server.URL)
	
	// Define response type
	type ResponseData struct {
		Method string `json:"method"`
	}
	
	// Test GET
	t.Run("GET", func(t *testing.T) {
		var responseData ResponseData
		err := client.Get(context.Background(), "/get", &responseData)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if responseData.Method != "GET" {
			t.Errorf("Expected response Method to be 'GET', got '%s'", responseData.Method)
		}
	})
	
	// Test POST
	t.Run("POST", func(t *testing.T) {
		var responseData ResponseData
		err := client.Post(context.Background(), "/post", map[string]string{"key": "value"}, &responseData)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if responseData.Method != "POST" {
			t.Errorf("Expected response Method to be 'POST', got '%s'", responseData.Method)
		}
	})
	
	// Test PUT
	t.Run("PUT", func(t *testing.T) {
		var responseData ResponseData
		err := client.Put(context.Background(), "/put", map[string]string{"key": "value"}, &responseData)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if responseData.Method != "PUT" {
			t.Errorf("Expected response Method to be 'PUT', got '%s'", responseData.Method)
		}
	})
	
	// Test PATCH
	t.Run("PATCH", func(t *testing.T) {
		var responseData ResponseData
		err := client.Patch(context.Background(), "/patch", map[string]string{"key": "value"}, &responseData)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if responseData.Method != "PATCH" {
			t.Errorf("Expected response Method to be 'PATCH', got '%s'", responseData.Method)
		}
	})
	
	// Test DELETE
	t.Run("DELETE", func(t *testing.T) {
		var responseData ResponseData
		err := client.Delete(context.Background(), "/delete", &responseData)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if responseData.Method != "DELETE" {
			t.Errorf("Expected response Method to be 'DELETE', got '%s'", responseData.Method)
		}
	})
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return s != "" && s != substr && len(s) >= len(substr) && s[0:len(substr)] == substr
}