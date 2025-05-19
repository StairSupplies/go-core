package rest_test

import (
	"fmt"
	"net/http"
	"time"

	"github.com/StairSupplies/go-core/logger"
	"github.com/StairSupplies/go-core/rest"
)

func Example() {
	// Create a client with default options
	client, err := rest.NewClient()
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}

	// Make a request
	// Note: In a real example we would make a real HTTP request
	// but for the example we just show the client creation
	fmt.Println("Client created successfully")

	_ = client // Prevent unused variable warning

	// Output: Client created successfully
}

func ExampleNewClient() {
	// Create a client with various options
	client, err := rest.NewClient(
		rest.WithBaseURL("https://api.example.com"),
		rest.WithTimeout(10*time.Second),
		rest.WithHeader("X-API-Key", "your-api-key"),
	)
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}

	fmt.Println("Client created successfully")
	_ = client // Prevent unused variable warning

	// Output: Client created successfully
}

func ExampleClient_Get() {
	// Create a client with a base URL
	client, err := rest.NewClient(
		rest.WithBaseURL("https://api.example.com"),
	)
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}

	// In a real example, we would make an actual API call
	// But for testing purposes, we'll just show the interface
	fmt.Println("GET request demonstrated")
	_ = client // Prevent unused variable warning

	// Output: GET request demonstrated
}

func ExampleClient_Post() {
	// Create a client
	client, err := rest.NewClient(
		rest.WithBaseURL("https://api.example.com"),
	)
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}

	// In a real example, we would make an actual API call
	// But for testing purposes, we'll just show the interface
	fmt.Println("POST request demonstrated")
	_ = client // Prevent unused variable warning

	// Output: POST request demonstrated
}

func ExampleClient_Put() {
	// Create a client
	client, err := rest.NewClient(
		rest.WithBaseURL("https://api.example.com"),
	)
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}

	// In a real example, we would make an actual API call
	// But for testing purposes, we'll just show the interface
	fmt.Println("PUT request demonstrated")
	_ = client // Prevent unused variable warning

	// Output: PUT request demonstrated
}

func ExampleClient_Patch() {
	// Create a client
	client, err := rest.NewClient(
		rest.WithBaseURL("https://api.example.com"),
	)
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}

	// In a real example, we would make an actual API call
	// But for testing purposes, we'll just show the interface
	fmt.Println("PATCH request demonstrated")
	_ = client // Prevent unused variable warning

	// Output: PATCH request demonstrated
}

func ExampleClient_Delete() {
	// Create a client
	client, err := rest.NewClient(
		rest.WithBaseURL("https://api.example.com"),
	)
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}

	// In a real example, we would make an actual API call
	// But for testing purposes, we'll just show the interface
	fmt.Println("DELETE request demonstrated")
	_ = client // Prevent unused variable warning

	// Output: DELETE request demonstrated
}

func ExampleWithBaseURL() {
	// Create a client with a base URL
	client, err := rest.NewClient(
		rest.WithBaseURL("https://api.example.com"),
	)
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}

	fmt.Println("Client created with base URL")
	_ = client // Prevent unused variable warning

	// Output: Client created with base URL
}

func ExampleWithTimeout() {
	// Create a client with a custom timeout
	client, err := rest.NewClient(
		rest.WithTimeout(10*time.Second),
	)
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}

	fmt.Println("Client created with custom timeout")
	_ = client // Prevent unused variable warning

	// Output: Client created with custom timeout
}

func ExampleWithHeader() {
	// Create a client with custom headers
	client, err := rest.NewClient(
		rest.WithHeader("X-API-Key", "your-api-key"),
		rest.WithHeader("User-Agent", "MyApp/1.0"),
	)
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}

	fmt.Println("Client created with custom headers")
	_ = client // Prevent unused variable warning

	// Output: Client created with custom headers
}

func ExampleWithHeaders() {
	// Create a client with multiple headers at once
	headers := map[string]string{
		"X-API-Key":  "your-api-key",
		"User-Agent": "MyApp/1.0",
		"Accept":     "application/json",
	}

	client, err := rest.NewClient(
		rest.WithHeaders(headers),
	)
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}

	fmt.Println("Client created with multiple headers")
	_ = client // Prevent unused variable warning

	// Output: Client created with multiple headers
}

func ExampleWithHTTPClient() {
	// Create a custom HTTP client
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	// Create a rest client with the custom HTTP client
	client, err := rest.NewClient(
		rest.WithHTTPClient(httpClient),
	)
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}

	fmt.Println("Client created with custom HTTP client")
	_ = client // Prevent unused variable warning

	// Output: Client created with custom HTTP client
}

func ExampleWithLogger() {
	// Create a custom logger
	log, err := logger.New(
		logger.WithLevel("debug"),
		logger.WithServiceName("api-client"),
	)
	if err != nil {
		fmt.Printf("Error creating logger: %v\n", err)
		return
	}

	// Create a client with the custom logger
	client, err := rest.NewClient(
		rest.WithLogger(log),
	)
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}

	fmt.Println("Client created with custom logger")
	_ = client // Prevent unused variable warning

	// Output: Client created with custom logger
}

func ExampleWithServiceName() {
	// Create a client with a service name
	client, err := rest.NewClient(
		rest.WithServiceName("user-service"),
	)
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}

	fmt.Println("Client created with service name")
	_ = client // Prevent unused variable warning

	// Output: Client created with service name
}

func ExampleNewClientError() {
	// Create a custom client error
	err := rest.NewClientError(
		rest.ErrNotFound,
		"User with ID 123 not found",
		"USER_NOT_FOUND",
	)

	// Use the error in your application
	fmt.Printf("Error type: %T\n", err)

	// Output: Error type: *rest.ClientError
}