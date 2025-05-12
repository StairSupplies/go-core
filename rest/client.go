package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is a simple HTTP client for making requests
type Client struct {
	BaseURL     string
	HTTPClient  *http.Client
	Headers     map[string]string
	DefaultTimeout time.Duration
}

// NewClient creates a new rest client with default settings
func NewClient() *Client {
	return &Client{
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		Headers:     make(map[string]string),
		DefaultTimeout: 30 * time.Second,
	}
}

// WithBaseURL sets the base URL for all requests
func (c *Client) WithBaseURL(baseURL string) *Client {
	c.BaseURL = baseURL
	return c
}

// WithHeader adds a header to all requests
func (c *Client) WithHeader(key, value string) *Client {
	c.Headers[key] = value
	return c
}

// WithTimeout sets the default timeout for requests
func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.DefaultTimeout = timeout
	c.HTTPClient.Timeout = timeout
	return c
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

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for non-2xx responses
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP error: %s, body: %s", resp.Status, string(respBody))
	}

	// If no response is expected, return nil
	if response == nil {
		return nil
	}

	// Parse the response
	if err := json.Unmarshal(respBody, response); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
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