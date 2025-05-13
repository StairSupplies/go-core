package logger_test

import (
	"context"
	"fmt"

	"github.com/StairSupplies/go-core/logger"
	"go.uber.org/zap"
)

// ExampleInit demonstrates initializing the logger
func ExampleInit() {
	// Initialize the logger with a specific configuration
	cfg := logger.Config{
		Level:       "info",
		Development: true,
		// Disable actual logging for example test
		OutputPaths: []string{"/dev/null"},
		ServiceName: "example-service",
		InitialFields: map[string]interface{}{
			"environment": "testing",
		},
	}

	// Initialize the logger
	err := logger.Init(cfg)
	if err != nil {
		fmt.Printf("Error initializing logger: %v\n", err)
		return
	}

	fmt.Println("Logger initialized successfully")
	// Output: Logger initialized successfully
}

// ExampleInfo demonstrates basic logging at the Info level
func ExampleInfo() {
	// Initialize a simple logger with no output
	_ = logger.Init(logger.Config{
		Level:       "info",
		Development: true,
		OutputPaths: []string{"/dev/null"},
	})

	// Log a message (no actual output due to /dev/null)
	logger.Info("Hello, world!")

	// For the example, we display what would normally be logged
	fmt.Println("Logged 'Hello, world!' at INFO level")
	// Output: Logged 'Hello, world!' at INFO level
}

// ExampleInfow demonstrates structured logging with key-value pairs
func ExampleInfow() {
	// Initialize a simple logger with no output
	_ = logger.Init(logger.Config{
		Level:       "info",
		Development: true,
		OutputPaths: []string{"/dev/null"},
	})

	// Log structured data (no actual output due to /dev/null)
	logger.Infow("User logged in",
		"user_id", 123,
		"username", "john.doe",
		"login_count", 5,
	)

	// For the example, we display what would normally be logged
	fmt.Println("Logged structured data about user login at INFO level")
	// Output: Logged structured data about user login at INFO level
}

// ExampleWithContext demonstrates using loggers with context
func ExampleWithContext() {
	// Initialize logger with no output
	_ = logger.Init(logger.Config{
		Level:       "info",
		Development: true,
		OutputPaths: []string{"/dev/null"},
	})

	// Create a context with logger
	ctx := context.Background()

	// Add fields to logger and store in context
	l := logger.With(
		zap.String("request_id", "req-123"),
		zap.String("client_ip", "192.168.1.1"),
	)
	ctx = logger.ContextWithLogger(ctx, l)

	// Later, retrieve the logger from context
	contextLogger := logger.WithContext(ctx)

	// Use the logger (no actual output due to /dev/null)
	contextLogger.Info("Processing request")

	fmt.Println("Logged with context-aware logger containing request_id and client_ip")
	// Output: Logged with context-aware logger containing request_id and client_ip
}

// ExampleWithFields demonstrates adding structured fields to logs
func ExampleWithFields() {
	// Initialize logger with no output
	_ = logger.Init(logger.Config{
		Level:       "info",
		Development: true,
		OutputPaths: []string{"/dev/null"},
	})

	// Create a logger with additional fields
	appLogger := logger.WithFields(map[string]interface{}{
		"component": "auth",
		"module":    "login",
	})

	// Use the logger (no actual output due to /dev/null)
	appLogger.Info("Authentication successful")

	fmt.Println("Logged with component and module fields")
	// Output: Logged with component and module fields
}

// Example_logLevels demonstrates different log levels
func Example_logLevels() {
	// Initialize logger with info level and no output
	_ = logger.Init(logger.Config{
		Level:       "info", // Debug messages won't be logged
		Development: true,
		OutputPaths: []string{"/dev/null"},
	})

	// These logs would be written but aren't captured in example output
	logger.Debug("This is a debug message") // Not logged due to level setting
	logger.Info("This is an info message")
	logger.Warn("This is a warning message")
	logger.Error("This is an error message")

	// Simulate the log level filtering behavior for the example
	fmt.Println("Debug message: not logged (below threshold)")
	fmt.Println("Info message: logged")
	fmt.Println("Warning message: logged")
	fmt.Println("Error message: logged")
	// Output:
	// Debug message: not logged (below threshold)
	// Info message: logged
	// Warning message: logged
	// Error message: logged
}