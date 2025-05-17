package rest

import (
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/StairSupplies/go-core/logger"
)

// ClientOption is a function that configures a Client
type ClientOption func(*Client)

// optionNames maps option functions to their string names
var optionNames = make(map[uintptr]string)

// registerOption registers an option function with its name
func registerOption(option ClientOption, name string) ClientOption {
	// Get function pointer
	ptr := reflect.ValueOf(option).Pointer()
	optionNames[ptr] = name
	return option
}

// OptionToString converts a ClientOption function to its string name
// Returns name of the option if it's a known option, otherwise returns a generic name
func OptionToString(option ClientOption) string {
	// Get function pointer
	ptr := reflect.ValueOf(option).Pointer()
	
	// Check if we have a registered name
	if name, ok := optionNames[ptr]; ok {
		return name
	}
	
	// Fallback to reflection for getting the function name
	fnName := runtime.FuncForPC(ptr).Name()
	if fnName == "" {
		return "UnknownOption"
	}
	
	// Extract just the function name (without package path)
	parts := strings.Split(fnName, ".")
	if len(parts) == 0 {
		return "UnknownOption"
	}
	
	return parts[len(parts)-1]
}

// WithBaseURL sets the base URL for the client
func WithBaseURL(baseURL string) ClientOption {
	return registerOption(func(c *Client) {
		c.BaseURL = baseURL
	}, "WithBaseURL")
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return registerOption(func(c *Client) {
		c.HTTPClient = httpClient
	}, "WithHTTPClient")
}

// WithHeader adds a header to all requests
func WithHeader(key, value string) ClientOption {
	return registerOption(func(c *Client) {
		c.Headers[key] = value
	}, "WithHeader")
}

// WithHeaders sets all headers for requests
func WithHeaders(headers map[string]string) ClientOption {
	return registerOption(func(c *Client) {
		for k, v := range headers {
			c.Headers[k] = v
		}
	}, "WithHeaders")
}

// WithTimeout sets the client timeout
func WithTimeout(timeout time.Duration) ClientOption {
	return registerOption(func(c *Client) {
		c.Timeout = timeout
	}, "WithTimeout")
}

// WithLogger sets a custom logger for the client.
func WithLogger(log *logger.Logger) ClientOption {
	return registerOption(func(c *Client) {
		c.Logger = log
	}, "WithLogger")
}

// WithServiceName sets the service name for the client.
func WithServiceName(name string) ClientOption {
	return registerOption(func(c *Client) {
		c.ServiceName = name
	}, "WithServiceName")
}
