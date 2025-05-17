package rest

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/StairSupplies/go-core/logger"
)

func TestWithBaseURL(t *testing.T) {
	client := &Client{}
	baseURL := "https://api.example.com"
	
	WithBaseURL(baseURL)(client)
	
	if client.BaseURL != baseURL {
		t.Errorf("WithBaseURL() = %q, want %q", client.BaseURL, baseURL)
	}
}

func TestWithHTTPClient(t *testing.T) {
	client := &Client{}
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	WithHTTPClient(httpClient)(client)
	
	if client.HTTPClient != httpClient {
		t.Errorf("WithHTTPClient() did not set the expected HTTP client")
	}
}

func TestWithHeader(t *testing.T) {
	client := &Client{
		Headers: make(map[string]string),
	}
	
	WithHeader("X-API-Key", "test-key")(client)
	
	if client.Headers["X-API-Key"] != "test-key" {
		t.Errorf("WithHeader() = %q, want %q", client.Headers["X-API-Key"], "test-key")
	}
}

func TestWithHeaders(t *testing.T) {
	client := &Client{
		Headers: make(map[string]string),
	}
	
	headers := map[string]string{
		"X-API-Key": "test-key",
		"User-Agent": "test-agent",
	}
	
	WithHeaders(headers)(client)
	
	for k, v := range headers {
		if client.Headers[k] != v {
			t.Errorf("WithHeaders() key %q = %q, want %q", k, client.Headers[k], v)
		}
	}
}

func TestWithTimeout(t *testing.T) {
	client := &Client{}
	timeout := 5 * time.Second
	
	WithTimeout(timeout)(client)
	
	if client.Timeout != timeout {
		t.Errorf("WithTimeout() = %v, want %v", client.Timeout, timeout)
	}
}

func TestWithLogger(t *testing.T) {
	client := &Client{}
	log, _ := logger.New()
	
	WithLogger(log)(client)
	
	if client.Logger != log {
		t.Errorf("WithLogger() did not set the expected logger")
	}
}

func TestWithServiceName(t *testing.T) {
	client := &Client{}
	serviceName := "test-service"
	
	WithServiceName(serviceName)(client)
	
	if client.ServiceName != serviceName {
		t.Errorf("WithServiceName() = %q, want %q", client.ServiceName, serviceName)
	}
}

func TestOptionToString_MoreOptions(t *testing.T) {
	// Additional tests for OptionToString beyond what's in client_test.go
	tests := []struct {
		name   string
		option ClientOption
		want   string
	}{
		{"WithHTTPClient", WithHTTPClient(&http.Client{}), "WithHTTPClient"},
		{"WithHeaders", WithHeaders(map[string]string{"key": "value"}), "WithHeaders"},
		{"WithLogger", WithLogger(nil), "WithLogger"},
		{"WithServiceName", WithServiceName("test"), "WithServiceName"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := OptionToString(tt.option)
			// The function name may contain additional info like function address
			// or fully qualified package name. We just want to check if it contains 
			// the expected function name.
			if !strings.Contains(got, tt.want) {
				t.Errorf("OptionToString() = %v, should contain %v", got, tt.want)
			}
		})
	}

	// Test with a non-public or anonymous function should return a more complex name
	// that may include package path or be empty
	customOption := func(c *Client) {}
	result := OptionToString(customOption)
	if result == "" {
		// This is acceptable; anonymous functions may not have a name
		// Alternatively, it could return a complex function name with package
		t.Logf("OptionToString for anonymous function returned empty string")
	}
}