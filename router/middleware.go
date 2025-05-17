package router

import (
	"net/http"
	"time"

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

			// Create a logger
			reqLogger, _ := logger.NewLogger(logger.Config{})
			
			// Prepare request logger with common fields
			requestLog := reqLogger.With(
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
			ctx := logger.NewContext(r.Context(), requestLog)
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