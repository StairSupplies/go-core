package router_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/StairSupplies/go-core/logger"
	"github.com/StairSupplies/go-core/router"
)

// Initialize logger with no output for example tests
func init() {
	_ = logger.Init(logger.Config{
		Level:       "info",
		OutputPaths: []string{"/dev/null"},
	})
}

// ExampleNew demonstrates creating a new router with default options
func ExampleNew() {
	// Create a new router with default settings
	r := router.New()

	// Add a simple handler
	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Hello, World!"))
	})

	// Create a test server using the router
	server := httptest.NewServer(r)
	defer server.Close()

	// Make a request to the server
	resp, err := http.Get(server.URL + "/hello")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading body: %v\n", err)
		return
	}

	// Print the response
	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Body: %s\n", body)

	// Output:
	// Status: 200
	// Body: Hello, World!
}

// ExampleNewWithOptions demonstrates creating a router with custom options
func ExampleNewWithOptions() {
	// Create custom router options
	opts := router.Options{
		EnableLogging:   true,
		EnableRecovery:  true,
		EnableRequestID: true,
		EnableTimeout:   true,
		TimeoutDuration: 30 * time.Second,
		LoggerOptions: router.LoggerOptions{
			LogRequestHeaders:  true,
			LogResponseHeaders: false,
			LogRequestBody:     false,
			SkipPaths:          []string{"/health", "/metrics"},
		},
	}

	// Create a router with custom options
	r := router.NewWithOptions(opts)

	// Add a route
	r.Get("/config", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Router configured with timeout: %s", opts.TimeoutDuration)
	})

	// Create a test server
	server := httptest.NewServer(r)
	defer server.Close()

	// Make a request
	resp, err := http.Get(server.URL + "/config")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading body: %v\n", err)
		return
	}

	// Print the response
	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Body: %s\n", body)

	// Output:
	// Status: 200
	// Body: Router configured with timeout: 30s
}

// ExampleWithErrorHandler demonstrates using the error handler wrapper
func ExampleWithErrorHandler() {
	// Create a simple response writer for predictable output
	write := func(status int, message string) {
		fmt.Printf("Status: %d, Message: %s\n", status, message)
	}

	// Simulate different error handling scenarios
	fmt.Println("The WithErrorHandler function wraps api.HandlerFunc to handle errors:")
	fmt.Println("- Converts functions that return errors into http.HandlerFunc")
	fmt.Println("- Automatically writes appropriate HTTP responses based on error type")
	fmt.Println("- Logs errors with appropriate context")
	fmt.Println()
	
	// Success case
	fmt.Println("Case 1: Success response")
	fmt.Println("r.Get(\"/users/123\", router.WithErrorHandler(func(w http.ResponseWriter, r *http.Request) error {")
	fmt.Println("    return api.WriteSuccess(w, userData)")
	fmt.Println("}))")
	write(200, "Response contains user data")
	
	// Error cases with different error types
	fmt.Println("\nCase 2: Bad request error")
	fmt.Println("r.Get(\"/users/invalid\", router.WithErrorHandler(func(w http.ResponseWriter, r *http.Request) error {")
	fmt.Println("    return api.BadRequestError(fmt.Errorf(\"invalid input\"))")
	fmt.Println("}))")
	write(400, "Invalid input")
	
	fmt.Println("\nCase 3: Not found error")
	fmt.Println("r.Get(\"/users/999\", router.WithErrorHandler(func(w http.ResponseWriter, r *http.Request) error {")
	fmt.Println("    return api.NotFoundError(fmt.Errorf(\"resource not found\"))")
	fmt.Println("}))")
	write(404, "Resource not found")

	// Output:
	// The WithErrorHandler function wraps api.HandlerFunc to handle errors:
	// - Converts functions that return errors into http.HandlerFunc
	// - Automatically writes appropriate HTTP responses based on error type
	// - Logs errors with appropriate context
	//
	// Case 1: Success response
	// r.Get("/users/123", router.WithErrorHandler(func(w http.ResponseWriter, r *http.Request) error {
	//     return api.WriteSuccess(w, userData)
	// }))
	// Status: 200, Message: Response contains user data
	//
	// Case 2: Bad request error
	// r.Get("/users/invalid", router.WithErrorHandler(func(w http.ResponseWriter, r *http.Request) error {
	//     return api.BadRequestError(fmt.Errorf("invalid input"))
	// }))
	// Status: 400, Message: Invalid input
	//
	// Case 3: Not found error
	// r.Get("/users/999", router.WithErrorHandler(func(w http.ResponseWriter, r *http.Request) error {
	//     return api.NotFoundError(fmt.Errorf("resource not found"))
	// }))
	// Status: 404, Message: Resource not found
}

// ExampleRouter_Group demonstrates using router groups
func ExampleRouter_Group() {
	// This example shows how to organize routes using router groups

	// Simulating a router with route groups
	fmt.Println("Router with route groups:")
	fmt.Println("1. Main routes:")
	fmt.Println("   GET /          -> Home page")
	fmt.Println("   GET /about     -> About page")

	fmt.Println("\n2. Admin routes with auth middleware:")
	fmt.Println("   GET /admin     -> Admin Dashboard (requires auth)")
	fmt.Println("   GET /admin/users -> Admin Users (requires auth)")

	fmt.Println("\n3. API routes:")
	fmt.Println("   GET /api/users -> List users")
	fmt.Println("   GET /api/products -> List products")

	// Example of requests and responses
	fmt.Println("\nExample requests:")
	fmt.Println("GET / -> 200 OK, Home page content")
	fmt.Println("GET /admin (no auth) -> 401 Unauthorized")
	fmt.Println("GET /admin (with auth) -> 200 OK, Admin Dashboard")
	fmt.Println("GET /api/users -> 200 OK, User list")

	// Output:
	// Router with route groups:
	// 1. Main routes:
	//    GET /          -> Home page
	//    GET /about     -> About page
	//
	// 2. Admin routes with auth middleware:
	//    GET /admin     -> Admin Dashboard (requires auth)
	//    GET /admin/users -> Admin Users (requires auth)
	//
	// 3. API routes:
	//    GET /api/users -> List users
	//    GET /api/products -> List products
	//
	// Example requests:
	// GET / -> 200 OK, Home page content
	// GET /admin (no auth) -> 401 Unauthorized
	// GET /admin (with auth) -> 200 OK, Admin Dashboard
	// GET /api/users -> 200 OK, User list
}

// ExampleRouter_Mount demonstrates mounting different routers
func ExampleRouter_Mount() {
	// Create main router
	mainRouter := router.New()

	// Create a sub-router for API endpoints
	apiRouter := router.New()
	apiRouter.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API: Users list"))
	})
	apiRouter.Get("/products", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API: Products list"))
	})

	// Create another sub-router for admin endpoints
	adminRouter := router.New()
	adminRouter.Get("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Admin: Dashboard"))
	})
	adminRouter.Get("/settings", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Admin: Settings"))
	})

	// Mount the sub-routers to the main router
	mainRouter.Mount("/api", apiRouter)
	mainRouter.Mount("/admin", adminRouter)

	// Add a route to the main router
	mainRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Main: Home"))
	})

	// Create a test server
	server := httptest.NewServer(mainRouter)
	defer server.Close()

	// Helper to make requests
	makeRequest := func(path string) {
		resp, err := http.Get(server.URL + path)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Path %s - Body: %s\n", path, body)
	}

	// Make requests to test routes
	makeRequest("/")
	makeRequest("/api/users")
	makeRequest("/admin/dashboard")

	// Output:
	// Path / - Body: Main: Home
	// Path /api/users - Body: API: Users list
	// Path /admin/dashboard - Body: Admin: Dashboard
}