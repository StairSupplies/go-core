package logger_test

import (
	"context"
	"fmt"
	"time"

	"github.com/StairSupplies/go-core/logger"
	"go.uber.org/zap"
)

func Example() {
	// Create a simple logger with default settings
	log, err := logger.New()
	if err != nil {
		fmt.Printf("Error creating logger: %v\n", err)
		return
	}

	// Log at different levels
	log.Info("This is an informational message")
	log.Warn("This is a warning message")
	log.Debug("This won't be visible at default 'info' level")
	log.Error("This is an error message")

	// No Output: Log output is not captured in examples
}

func ExampleNew() {
	// Create a logger with custom options
	log, err := logger.New(
		logger.WithLevel("debug"),
		logger.WithDevelopmentMode(true),
		logger.WithServiceName("example-service"),
		logger.WithInitialFields(map[string]interface{}{
			"environment": "development",
			"version":     "1.0.0",
		}),
	)
	if err != nil {
		fmt.Printf("Error creating logger: %v\n", err)
		return
	}

	// Now all log entries will include the service name and initial fields
	log.Info("Service started")

	// No Output: Log output is not captured in examples
}

func ExampleLogger_WithFields() {
	// Create a base logger
	log, err := logger.New()
	if err != nil {
		fmt.Printf("Error creating logger: %v\n", err)
		return
	}

	// Create a child logger with additional fields
	requestLogger := log.WithFields(map[string]interface{}{
		"request_id": "abc-123",
		"user_id":    "user-456",
		"path":       "/api/items",
	})

	// Log with the request context
	requestLogger.Info("Processing request")
	requestLogger.Error("Request failed")

	// No Output: Log output is not captured in examples
}

func ExampleLogger_With() {
	// Create a base logger
	log, err := logger.New()
	if err != nil {
		fmt.Printf("Error creating logger: %v\n", err)
		return
	}

	// Create a child logger with additional structured fields
	requestLogger := log.With(
		zap.String("request_id", "abc-123"),
		zap.String("user_id", "user-456"),
		zap.String("path", "/api/items"),
	)

	// Log with the request context
	requestLogger.Info("Processing request")

	// No Output: Log output is not captured in examples
}

func ExampleLogger_Infow() {
	// Create a logger
	log, err := logger.New()
	if err != nil {
		fmt.Printf("Error creating logger: %v\n", err)
		return
	}

	// Log with structured key-value pairs (equivalent to zapcore.Field)
	log.Infow("User logged in",
		"user_id", "user-123",
		"login_method", "oauth",
		"client_ip", "192.168.1.1",
	)

	// No Output: Log output is not captured in examples
}

func ExampleLogger_Debugf() {
	// Create a logger with debug level
	log, err := logger.New(logger.WithLevel("debug"))
	if err != nil {
		fmt.Printf("Error creating logger: %v\n", err)
		return
	}

	// Log with formatting (like fmt.Printf)
	userId := "user-123"
	log.Debugf("Processing request for user %s", userId)

	// No Output: Log output is not captured in examples
}

func ExampleNewContext() {
	// Create a logger
	log, err := logger.New(
		logger.WithServiceName("api-service"),
	)
	if err != nil {
		fmt.Printf("Error creating logger: %v\n", err)
		return
	}

	// Create a request-specific logger
	requestLog := log.With(
		zap.String("request_id", "req-abc123"),
		zap.String("path", "/api/users"),
	)

	// Store the logger in the context
	baseCtx := context.Background()
	ctx := logger.NewContext(baseCtx, requestLog)

	// Later, retrieve the logger from context
	ProcessRequest(ctx)

	// No Output: Log output is not captured in examples
}

// Helper function to demonstrate WithContext
func ProcessRequest(ctx context.Context) {
	// Get the logger from context
	log := logger.WithContext(ctx)

	// Log with the context-specific logger
	log.Info("Processing request")
}

func ExampleWithContext() {
	// Create a context with no logger
	ctx := context.Background()

	// Try to get a logger from the context
	// This will return a default logger since none exists in the context
	log := logger.WithContext(ctx)

	// Use the logger
	log.Info("Using default logger from context")

	// No Output: Log output is not captured in examples
}

func ExampleNewNopLogger() {
	// Create a no-op logger for testing
	log := logger.NewNopLogger()

	// These logs will be silently discarded
	log.Info("This will not be logged")
	log.Error("This error will also be discarded",
		zap.String("reason", "using no-op logger"),
	)

	// No Output: Log output is not captured in examples
}

func ExampleLogger_Debug() {
	// Create a logger with debug level
	log, err := logger.New(logger.WithLevel("debug"))
	if err != nil {
		fmt.Printf("Error creating logger: %v\n", err)
		return
	}

	// Log with structured fields
	log.Debug("Database connection established",
		zap.String("host", "localhost"),
		zap.Int("port", 5432),
		zap.Duration("connect_time", 50*time.Millisecond),
	)

	// No Output: Log output is not captured in examples
}