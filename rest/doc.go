/*
Package rest provides a fluent REST client for making HTTP requests with error handling and retry capabilities.

This package simplifies making HTTP API requests with features like automatic JSON
serialization/deserialization, retries with exponential backoff, and standardized error handling.
It's designed to make API integration more reliable and consistent.

# Features

  - Fluent interface for HTTP methods (GET, POST, PUT, PATCH, DELETE)
  - Automatic JSON request/response serialization
  - Configurable with functional options pattern
  - Automatic retry with exponential backoff
  - Standardized error handling with typed errors
  - Integration with the go-core/logger package
  - Context support for cancellation and timeouts

# Basic Usage

Creating a client and making requests:

	// Create a new client with a base URL
	client, err := rest.NewClient(
		rest.WithBaseURL("https://api.example.com"),
		rest.WithTimeout(10 * time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Make a GET request
	var user User
	err = client.Get(ctx, "/users/123", &user)
	if err != nil {
		log.Fatal(err)
	}

	// Make a POST request
	newUser := User{Name: "Jane Doe", Email: "jane@example.com"}
	var createdUser User
	err = client.Post(ctx, "/users", newUser, &createdUser)
	if err != nil {
		log.Fatal(err)
	}

# Client Configuration

The client can be configured with various options:

	client, err := rest.NewClient(
		rest.WithBaseURL("https://api.example.com"),
		rest.WithTimeout(10 * time.Second),
		rest.WithHeader("X-API-Key", "your-api-key"),
		rest.WithServiceName("user-service"),
	)

# Error Handling

The package provides standardized error handling:

	var user User
	err := client.Get(ctx, "/users/123", &user)
	if err != nil {
		switch {
		case errors.Is(err, rest.ErrNotFound):
			fmt.Println("User not found")
		case errors.Is(err, rest.ErrUnauthorized):
			fmt.Println("Authentication required")
		case errors.Is(err, rest.ErrServerError):
			fmt.Println("Server error occurred")
		default:
			fmt.Printf("Unknown error: %v", err)
		}
	}

# Advanced HTTP Client Configuration

For more control, you can provide a custom HTTP client:

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	client, err := rest.NewClient(
		rest.WithHTTPClient(httpClient),
	)

# Logging Integration

The client integrates with the go-core/logger package:

	logger, err := logger.New(
		logger.WithLevel("debug"),
		logger.WithServiceName("api-client"),
	)
	if err != nil {
		log.Fatal(err)
	}

	client, err := rest.NewClient(
		rest.WithLogger(logger),
	)
*/
package rest