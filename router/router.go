package router

import (
	"net/http"
	"time"

	"github.com/StairSupplies/go-core/api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Router extends chi.Router with additional functionality
// such as preconfigured middleware and error handling.
type Router struct {
	chi.Router
	options Options
}

// Options configures the router and its middleware.
// It allows for fine-grained control over which middleware to enable
// and how they should behave.
type Options struct {
	// EnableLogging enables the logging middleware
	EnableLogging bool
	// EnableRecovery enables panic recovery middleware
	EnableRecovery bool
	// EnableRequestID enables request ID middleware
	EnableRequestID bool
	// EnableTimeout enables timeout middleware
	EnableTimeout bool
	// EnableHealthcheck enables the healthcheck middleware
	EnableHealthcheck bool
	// TimeoutDuration sets the timeout for requests
	TimeoutDuration time.Duration
	// LoggerOptions configures the logger middleware
	LoggerOptions LoggerOptions
}

// LoggerOptions configures the logger middleware.
// It allows for controlling what information is logged for each request.
type LoggerOptions struct {
	// LogRequestHeaders determines if request headers should be logged
	LogRequestHeaders bool
	// LogResponseHeaders determines if response headers should be logged
	LogResponseHeaders bool
	// LogRequestBody determines if request body should be logged
	LogRequestBody bool
	// SkipPaths lists paths that should not be logged
	SkipPaths []string
}

// DefaultOptions returns the default router options.
// These defaults provide a balance of functionality and performance.
func DefaultOptions() Options {
	return Options{
		EnableLogging:     true,
		EnableRecovery:    true,
		EnableRequestID:   true,
		EnableTimeout:     true,
		EnableHealthcheck: true,
		TimeoutDuration:   60 * time.Second,
		LoggerOptions: LoggerOptions{
			LogRequestHeaders:  false,
			LogResponseHeaders: false,
			LogRequestBody:     false,
			SkipPaths:          []string{"/health", "/metrics"},
		},
	}
}

// New creates a new router with default options.
// This is the recommended way to create a router for most applications.
func New() *Router {
	return NewWithOptions(DefaultOptions())
}

// NewWithOptions creates a new router with the specified options.
// Use this when you need to customize the router's behavior.
func NewWithOptions(options Options) *Router {
	r := chi.NewRouter()

	// Apply middleware based on options
	if options.EnableRequestID {
		r.Use(middleware.RequestID)
	}

	if options.EnableRecovery {
		r.Use(Recoverer())
	}

	if options.EnableLogging {
		r.Use(Logger(options.LoggerOptions))
	}

	if options.EnableTimeout {
		r.Use(middleware.Timeout(options.TimeoutDuration))
	}

	if options.EnableHealthcheck {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
	}

	return &Router{
		Router:  r,
		options: options,
	}
}

// WithErrorHandler wraps an api.HandlerFunc to handle errors.
// It converts functions that return errors into standard http.HandlerFunc.
func WithErrorHandler(h api.HandlerFunc) http.HandlerFunc {
	return api.WrapHandler(h)
}

// Group creates a new router group with the same options.
// This is useful for applying middleware to a group of routes.
func (r *Router) Group(fn func(r chi.Router)) chi.Router {
	subRouter := &Router{
		Router:  chi.NewRouter(),
		options: r.options,
	}

	// Apply middleware to subRouter if needed
	if r.options.EnableRequestID {
		subRouter.Use(middleware.RequestID)
	}

	if r.options.EnableRecovery {
		subRouter.Use(Recoverer())
	}

	if r.options.EnableLogging {
		subRouter.Use(Logger(r.options.LoggerOptions))
	}

	if r.options.EnableTimeout {
		subRouter.Use(middleware.Timeout(r.options.TimeoutDuration))
	}

	if r.options.EnableHealthcheck {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
	}

	fn(subRouter)
	return subRouter
}

// Mount attaches another router to the specified pattern.
// This is useful for organizing routes into separate modules.
func (r *Router) Mount(pattern string, h http.Handler) {
	r.Router.Mount(pattern, h)
}

// WithMiddleware creates a copy of the router with additional middleware.
// This allows you to add custom middleware to a router.
func (r *Router) WithMiddleware(middlewares ...func(http.Handler) http.Handler) *Router {
	newRouter := &Router{
		Router:  r.Router,
		options: r.options,
	}

	for _, m := range middlewares {
		newRouter.Use(m)
	}

	return newRouter
}

// ServeHTTP implements the http.Handler interface.
// This allows the router to be used directly with http.ListenAndServe.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.Router.ServeHTTP(w, req)
}
