package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/StairSupplies/go-core/logger"
)

// Client is an enhanced HTTP client for making API requests
type Client struct {
	BaseURL     string
	HTTPClient  *http.Client
	Headers     map[string]string
	Retries     int
	Timeout     time.Duration
	Logger      *logger.Logger
	ServiceName string
}

// NewClient creates a new rest client with the provided options
func NewClient(opts ...ClientOption) (*Client, error) {
	// Initialize client with default values
	c := &Client{
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		Headers: make(map[string]string),
		Retries: 3,
		Timeout: 30 * time.Second,
	}

	// Apply options
	for _, opt := range opts {
		opt(c)
	}

	// Create default logger if not provided
	if c.Logger == nil {
		// Start with standard options
		defaultOptions := []logger.Option{
			logger.WithServiceName(c.ServiceName),
			logger.WithInitialFields(map[string]interface{}{
				"component": "client",
			}),
		}

		// Create the logger with the combined options
		defaultLogger, err := logger.New(defaultOptions...)
		if err != nil {
			return nil, err
		}
		c.Logger = defaultLogger
	}

	// Configure client timeout
	c.HTTPClient.Timeout = c.Timeout

	return c, nil
}

// Request performs an HTTP request and returns the response
func (c *Client) Request(ctx context.Context, method, path string, body interface{}, response interface{}) error {
	url := path
	if c.BaseURL != "" {
		url = c.BaseURL + path
	}

	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set default headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Set custom headers
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}

	// Perform the request with retries
	var resp *http.Response

	for attempt := 0; attempt <= c.Retries; attempt++ {
		resp, err = c.HTTPClient.Do(req)
		if err == nil {
			break
		}

		// Check if the error is due to context cancellation
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return err
		}

		if attempt < c.Retries {
			// Exponential backoff
			backoffTime := time.Duration(attempt+1) * 100 * time.Millisecond
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoffTime):
				// Continue with retry
			}
		}
	}

	if err != nil {
		return &ClientError{
			Err:     ErrConnectionFailed,
			Message: err.Error(),
		}
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &ClientError{
			Err:     ErrConnectionFailed,
			Message: fmt.Sprintf("failed to read response body: %s", err),
		}
	}

	// Check for non-2xx responses
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Try to parse as error response
		var errResp struct {
			Message string `json:"message,omitempty"`
			Code    string `json:"code,omitempty"`
		}

		if jsonErr := json.Unmarshal(respBody, &errResp); jsonErr == nil && errResp.Message != "" {
			return &ClientError{
				Err:     getErrorByStatusCode(resp.StatusCode),
				Message: errResp.Message,
				Code:    errResp.Code,
			}
		}

		// Fall back to generic error
		return &ClientError{
			Err:     getErrorByStatusCode(resp.StatusCode),
			Message: string(respBody),
			Code:    fmt.Sprintf("%d", resp.StatusCode),
		}
	}

	// If no response is expected, return nil
	if response == nil {
		return nil
	}

	// Parse the response
	if err := json.Unmarshal(respBody, response); err != nil {
		return &ClientError{
			Err:     ErrInvalidRequest,
			Message: fmt.Sprintf("failed to parse response: %s", err),
		}
	}

	return nil
}

// Get makes a GET request
func (c *Client) Get(ctx context.Context, path string, response interface{}) error {
	return c.Request(ctx, http.MethodGet, path, nil, response)
}

// Post makes a POST request
func (c *Client) Post(ctx context.Context, path string, body interface{}, response interface{}) error {
	return c.Request(ctx, http.MethodPost, path, body, response)
}

// Put makes a PUT request
func (c *Client) Put(ctx context.Context, path string, body interface{}, response interface{}) error {
	return c.Request(ctx, http.MethodPut, path, body, response)
}

// Patch makes a PATCH request
func (c *Client) Patch(ctx context.Context, path string, body interface{}, response interface{}) error {
	return c.Request(ctx, http.MethodPatch, path, body, response)
}

// Delete makes a DELETE request
func (c *Client) Delete(ctx context.Context, path string, response interface{}) error {
	return c.Request(ctx, http.MethodDelete, path, nil, response)
}
