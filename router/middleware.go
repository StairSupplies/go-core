package router

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/StairSupplies/go-core/api"
	"github.com/StairSupplies/go-core/logger"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// Logger is a middleware that logs the start and end of each request.
// It provides structured logging with details about the request and response.
// Options can be used to customize what information is logged.
func Logger(opts LoggerOptions) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip logging for specified paths
			for _, path := range opts.SkipPaths {
				if r.URL.Path == path {
					next.ServeHTTP(w, r)
					return
				}
			}

			// Start timer
			start := time.Now()
			
			// Create a response writer wrapper to capture status code
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			
			// Prepare request logger with common fields
			requestLog := logger.With(
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("request_id", middleware.GetReqID(r.Context())),
				zap.String("remote_addr", r.RemoteAddr),
			)
			
			// Log request headers if enabled
			if opts.LogRequestHeaders {
				for k, v := range r.Header {
					if len(v) > 0 {
						requestLog = requestLog.With(zap.String("req_header_"+k, v[0]))
					}
				}
			}
			
			// Add logger to request context
			ctx := logger.ContextWithLogger(r.Context(), requestLog)
			r = r.WithContext(ctx)
			
			// Log request start
			requestLog.Info("HTTP request started")
			
			// Process request
			next.ServeHTTP(ww, r)
			
			// Log response
			duration := time.Since(start)
			responseLog := requestLog.With(
				zap.Int("status", ww.Status()),
				zap.Duration("duration", duration),
				zap.Int("size", ww.BytesWritten()),
			)
			
			// Log response headers if enabled
			if opts.LogResponseHeaders {
				for k, v := range ww.Header() {
					if len(v) > 0 {
						responseLog = responseLog.With(zap.String("resp_header_"+k, v[0]))
					}
				}
			}
			
			responseLog.Info("HTTP request completed")
		})
	}
}

// Recoverer recovers from panics and logs them.
// It ensures that the application doesn't crash when a panic occurs in a handler,
// and instead returns a 500 Internal Server Error response.
func Recoverer() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					stack := debug.Stack()
					
					// Log the panic
					reqID := middleware.GetReqID(r.Context())
					ctx := r.Context()
					
					logger.WithContext(ctx).Error("panic recovered",
						zap.Any("panic", rvr),
						zap.String("request_id", reqID),
						zap.String("stack", string(stack)),
					)
					
					// Return a 500 error
					err := fmt.Errorf("internal server error")
					api.WriteError(w, api.ServerError(err))
				}
			}()
			
			next.ServeHTTP(w, r)
		})
	}
}

// RequestID sets a unique ID for each request.
// This is a wrapper around chi's RequestID middleware.
// Request IDs are used for tracing requests in logs and responses.
func RequestID(next http.Handler) http.Handler {
	return middleware.RequestID(next)
}

// Timeout sets a timeout for the request.
// This is a wrapper around chi's Timeout middleware.
// It helps prevent long-running requests from consuming resources indefinitely.
func Timeout(duration time.Duration) func(next http.Handler) http.Handler {
	return middleware.Timeout(duration)
}