package router_test

import (
	"fmt"
	"net/http"
	"time"

	"github.com/StairSupplies/go-core/api"
	"github.com/StairSupplies/go-core/router"
	"github.com/go-chi/chi/v5"
)

func Example() {
	// Create a router with default options
	r := router.New()

	// Add a simple route
	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// Add a route with URL parameters
	r.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		w.Write([]byte(fmt.Sprintf("User ID: %s", id)))
	})

	// Use the router as http.Handler
	// http.ListenAndServe(":8080", r)
	fmt.Println("Router created with 2 routes")

	// Output: Router created with 2 routes
}

func ExampleNewWithOptions() {
	// Create custom options
	options := router.Options{
		EnableLogging:     true,
		EnableRecovery:    true,
		EnableRequestID:   true,
		EnableTimeout:     true,
		EnableHealthcheck: true,
		TimeoutDuration:   30 * time.Second,
		LoggerOptions: router.LoggerOptions{
			LogRequestHeaders:  true,
			LogResponseHeaders: false,
			LogRequestBody:     false,
			SkipPaths:          []string{"/healthz", "/metrics"},
		},
	}

	// Create a router with custom options
	r := router.NewWithOptions(options)

	// Add routes as needed
	r.Get("/api/users", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("List of users"))
	})

	fmt.Println("Router created with custom options")

	// Output: Router created with custom options
}

func ExampleWithErrorHandler() {
	// Create a router
	r := router.New()

	// Define a handler that returns an error
	userHandler := func(w http.ResponseWriter, r *http.Request) error {
		// Get user ID from URL
		id := chi.URLParam(r, "id")

		// Check if ID is valid
		if id == "" {
			return api.BadRequestError(fmt.Errorf("missing user ID"))
		}

		// Return success response
		return api.WriteSuccess(w, map[string]string{
			"id":   id,
			"name": "John Doe",
		})
	}

	// Use WithErrorHandler to handle errors from the handler
	r.Get("/users/{id}", router.WithErrorHandler(userHandler))

	fmt.Println("Route with error handling registered")

	// Output: Route with error handling registered
}

func ExampleRouter_Group() {
	// Create a router
	r := router.New()

	// Create an API router group
	apiRouter := r.Group(func(r chi.Router) {
		// Add routes to the API group
		r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("List of users"))
		})

		r.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			w.Write([]byte(fmt.Sprintf("User: %s", id)))
		})

		r.Get("/products", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("List of products"))
		})
	})

	// Mount the API router at /api
	r.Mount("/api", apiRouter)

	fmt.Println("API routes registered under /api")

	// Output: API routes registered under /api
}

func ExampleRouter_WithMiddleware() {
	// Create a router
	r := router.New()

	// Define a custom middleware
	authMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check for authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}

	// Create a router with custom middleware
	protectedRouter := r.WithMiddleware(authMiddleware)

	// Add protected routes
	protectedRouter.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Admin area"))
	})

	fmt.Println("Router with auth middleware created")

	// Output: Router with auth middleware created
}

func ExampleTimeout() {
	// Create a router
	r := router.New()

	// Define a custom timeout handler
	customTimeout := router.Timeout(5 * time.Second)

	// Apply to specific routes
	r.With(customTimeout).Get("/slow-operation", func(w http.ResponseWriter, r *http.Request) {
		// This operation will timeout if it takes longer than 5 seconds
		time.Sleep(1 * time.Second)
		w.Write([]byte("Operation completed"))
	})

	fmt.Println("Route with custom timeout registered")

	// Output: Route with custom timeout registered
}

func ExampleRequestID() {
	// Create a router
	r := router.New()

	// Define a handler that uses the request ID
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Get the request ID from the context
		// In a real application, you would use middleware.GetReqID(r.Context())
		w.Write([]byte("Request processed"))
	}

	// Apply RequestID middleware to specific routes
	r.With(router.RequestID).Get("/api/resource", handler)

	fmt.Println("Route with RequestID middleware registered")

	// Output: Route with RequestID middleware registered
}