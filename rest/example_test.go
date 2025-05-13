package rest_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/StairSupplies/go-core/rest"
)

// ExampleNewClient demonstrates creating a new REST client with default settings
func ExampleNewClient() {
	// Create a new client with default settings
	client := rest.NewClient()

	// You can chain configuration methods
	client.WithBaseURL("https://api.example.com")
	client.WithHeader("Authorization", "Bearer token123")
	client.WithTimeout(10 * time.Second)

	// Print client configuration
	fmt.Printf("Base URL: %s\n", client.BaseURL)
	fmt.Printf("Timeout: %s\n", client.DefaultTimeout)
	fmt.Printf("Authorization Header: %s\n", client.Headers["Authorization"])

	// Output:
	// Base URL: https://api.example.com
	// Timeout: 10s
	// Authorization Header: Bearer token123
}

// ExampleClient_Get demonstrates making a GET request
func ExampleClient_Get() {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method and path
		if r.Method != http.MethodGet || r.URL.Path != "/users/123" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		// Set response headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Write response body
		fmt.Fprint(w, `{"id": 123, "name": "John Doe", "email": "john@example.com"}`)
	}))
	defer server.Close()

	// Create client with the test server URL
	client := rest.NewClient().WithBaseURL(server.URL)

	// Define a struct to hold the response
	type User struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	// Make the request
	var user User
	err := client.Get(context.Background(), "/users/123", &user)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Print the response
	fmt.Printf("User ID: %d\n", user.ID)
	fmt.Printf("Name: %s\n", user.Name)
	fmt.Printf("Email: %s\n", user.Email)

	// Output:
	// User ID: 123
	// Name: John Doe
	// Email: john@example.com
}

// ExampleClient_Post demonstrates making a POST request
func ExampleClient_Post() {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method and path
		if r.Method != http.MethodPost || r.URL.Path != "/users" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		// Read and validate request body
		var user map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Add ID to the user
		user["id"] = 456
		
		// Write response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}))
	defer server.Close()

	// Create client
	client := rest.NewClient().WithBaseURL(server.URL)

	// Request data
	requestBody := map[string]interface{}{
		"name":  "Jane Smith",
		"email": "jane@example.com",
	}

	// Response data
	var responseData map[string]interface{}

	// Make the request
	err := client.Post(context.Background(), "/users", requestBody, &responseData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Print the response
	fmt.Printf("User created with ID: %.0f\n", responseData["id"])
	fmt.Printf("Name: %s\n", responseData["name"])
	fmt.Printf("Email: %s\n", responseData["email"])

	// Output:
	// User created with ID: 456
	// Name: Jane Smith
	// Email: jane@example.com
}

// ExampleClient_WithHeader demonstrates customizing headers
func ExampleClient_WithHeader() {
	// Create a test server that echoes headers
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return the received headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		
		// Create a response with the received headers
		response := map[string]string{
			"authorization": r.Header.Get("Authorization"),
			"x-api-key":     r.Header.Get("X-API-Key"),
			"user-agent":    r.Header.Get("User-Agent"),
		}
		
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client with custom headers
	client := rest.NewClient().
		WithBaseURL(server.URL).
		WithHeader("Authorization", "Bearer token456").
		WithHeader("X-API-Key", "api-key-789").
		WithHeader("User-Agent", "CustomApp/1.0")

	// Make request
	var response map[string]string
	err := client.Get(context.Background(), "/headers", &response)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Print headers received by the server
	fmt.Printf("Authorization: %s\n", response["authorization"])
	fmt.Printf("X-API-Key: %s\n", response["x-api-key"])
	fmt.Printf("User-Agent: %s\n", response["user-agent"])

	// Output:
	// Authorization: Bearer token456
	// X-API-Key: api-key-789
	// User-Agent: CustomApp/1.0
}

// ExampleClient_Request demonstrates using the Request method directly
func ExampleClient_Request() {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method and path
		if r.Method != http.MethodPatch || r.URL.Path != "/products/789" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		// Set response headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Write response body
		fmt.Fprint(w, `{"id": 789, "name": "Updated Product", "price": 29.99, "status": "active"}`)
	}))
	defer server.Close()

	// Create client
	client := rest.NewClient().WithBaseURL(server.URL)

	// Define request and response types
	type PatchRequest struct {
		Name   string  `json:"name"`
		Price  float64 `json:"price"`
		Status string  `json:"status"`
	}

	type Product struct {
		ID     int     `json:"id"`
		Name   string  `json:"name"`
		Price  float64 `json:"price"`
		Status string  `json:"status"`
	}

	// Create request payload
	patchData := PatchRequest{
		Name:   "Updated Product",
		Price:  29.99,
		Status: "active",
	}

	// Make the request using Request method directly
	var product Product
	err := client.Request(
		context.Background(),
		http.MethodPatch,
		"/products/789",
		patchData,
		&product,
	)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Print the response
	fmt.Printf("Product ID: %d\n", product.ID)
	fmt.Printf("Updated name: %s\n", product.Name)
	fmt.Printf("Updated price: %.2f\n", product.Price)
	fmt.Printf("Updated status: %s\n", product.Status)

	// Output:
	// Product ID: 789
	// Updated name: Updated Product
	// Updated price: 29.99
	// Updated status: active
}