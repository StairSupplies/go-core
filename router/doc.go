/*
Package router provides an opinionated HTTP router with middleware for Go web applications.

This package builds on the github.com/go-chi/chi/v5 router, adding preconfigured middleware,
standardized error handling, and integration with the api package. It's designed to provide
a consistent approach to building HTTP services with sensible defaults.

# Features

  - Integration with the go-core/api package for error handling
  - Structured logging with the go-core/logger package
  - Request tracing with unique request IDs
  - Panic recovery with proper error responses
  - Timeout handling
  - Configurable middleware options

# Basic Usage

Creating a new router with default middleware:

	// Create a new router with default middleware
	r := router.New()

	// Add routes
	r.Get("/api/health", healthCheckHandler)
	r.Post("/api/users", createUserHandler)

	// Start the server
	http.ListenAndServe(":8080", r)

# Error Handling

The router integrates with the api package for standardized error handling:

	// Define a handler that returns errors
	func getUserHandler(w http.ResponseWriter, r *http.Request) error {
	    id := chi.URLParam(r, "id")
	    user, err := userService.GetUser(id)
	    if err != nil {
	        return api.NotFoundError(fmt.Errorf("user not found: %w", err))
	    }
	    return api.WriteSuccess(w, user)
	}

	// Wrap the handler in the router
	r.Get("/api/users/{id}", router.WithErrorHandler(getUserHandler))

# Route Groups and Middleware

You can create route groups with shared middleware:

	// Create a router
	r := router.New()

	// Create an authenticated API group
	apiRouter := r.Group(func(r chi.Router) {
	    // Apply auth middleware only to this group
	    r.Use(authMiddleware)

	    // Add protected routes
	    r.Get("/users", router.WithErrorHandler(listUsersHandler))
	    r.Post("/users", router.WithErrorHandler(createUserHandler))
	})

	// Mount the API router
	r.Mount("/api", apiRouter)

# Custom Configuration

You can customize the router's middleware:

	r := router.NewWithOptions(router.Options{
	    EnableLogging:   true,
	    EnableRecovery:  true,
	    EnableRequestID: true,
	    EnableTimeout:   true,
	    TimeoutDuration: 30 * time.Second,
	    LoggerOptions: router.LoggerOptions{
	        LogRequestHeaders:  true,
	        LogResponseHeaders: false,
	        SkipPaths:          []string{"/health", "/metrics"},
	    },
	})

# Logging

The router uses the go-core/logger package for structured logging of requests:

	// Initialize the logger first
	logger.Init(logger.Config{
	    Level:       "info",
	    Development: true,
	    ServiceName: "api-service",
	})

	// Create a new router (logging is enabled by default)
	r := router.New()
*/
package router