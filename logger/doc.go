/*
Package logger provides structured logging using Uber's Zap.

It offers both global and instance-based loggers with context-aware logging 
and support for both structured and formatted logging.

# Global vs Instance-Based Loggers

This package provides two approaches to logging:

1. Global logger:
   - Singleton logger accessible from anywhere via package-level functions
   - Suitable for simple applications where a single logger configuration is sufficient

2. Instance-based loggers:
   - Multiple independent logger instances with separate configurations
   - Ideal for libraries or applications that need different logging configurations
     for different components or when being consumed by other applications

# Global Logger Initialization

Initialize the global logger:

	cfg := logger.Config{
	    Level:       "info",
	    Development: true,
	    ServiceName: "user-service",
	    InitialFields: map[string]interface{}{
	        "version": "1.0.0",
	    },
	}

	if err := logger.Init(cfg); err != nil {
	    panic(err)
	}

# Instance-Based Logger Creation

Create multiple logger instances:

	// Create an API logger
	apiLogger, err := logger.NewLogger(logger.Config{
	    Level:       "info",
	    Development: true,
	    ServiceName: "api-service",
	    InitialFields: map[string]interface{}{
	        "component": "api",
	    },
	})
	if err != nil {
	    panic(err)
	}

	// Create a database logger
	dbLogger, err := logger.NewLogger(logger.Config{
	    Level:       "debug",
	    Development: true,
	    ServiceName: "db-service",
	    InitialFields: map[string]interface{}{
	        "component": "database",
	    },
	})
	if err != nil {
	    panic(err)
	}

	// Use the loggers independently
	apiLogger.Info("API server started")
	dbLogger.Debug("Database connection established")

# Basic Logging

Log messages at different levels:

	// Using global logger
	logger.Debug("Processing item",
	    zap.Int("item_id", item.ID),
	    zap.String("status", item.Status),
	)

	// Using logger instance
	myLogger.Info("Application started",
	    zap.String("version", version),
	    zap.String("environment", env),
	)

# Formatted Logging

Log with string formatting (using the sugared logger):

	// Using global logger
	logger.Debugf("Processing item %d with status %s", item.ID, item.Status)

	// Using logger instance
	myLogger.Infof("Application started with version %s in %s environment", version, env)

# Structured Logging

Log with structured key-value pairs (using the sugared logger):

	// Using global logger
	logger.Debugw("Processing item", 
		"item_id", item.ID, 
		"status", item.Status,
	)

	// Using logger instance
	myLogger.Infow("Application started", 
		"version", version, 
		"environment", env,
	)

# Context-Aware Logging

Use context to pass loggers:

	// Global logger approach
	globalLoggerWithFields := logger.With(zap.String("request_id", requestID))
	ctx = logger.ContextWithLogger(ctx, globalLoggerWithFields)

	// Instance logger approach
	myLoggerWithFields := myLogger.With(zap.String("request_id", requestID))
	ctx = logger.ContextWithLogger(ctx, myLoggerWithFields)

	// Get logger from context
	func HandleRequest(ctx context.Context) {
	    log := logger.WithContext(ctx)
	    log.Info("Processing request")
	}

# Structured Fields

Create child loggers with additional fields:

	// Using the global structured logger
	loggerWithFields := logger.With(
	    zap.String("user_id", user.ID),
	    zap.String("request_id", requestID),
	)
	loggerWithFields.Info("User logged in")

	// Using a logger instance with the sugared logger
	myLoggerWithFields := myLogger.WithFields(map[string]interface{}{
	    "user_id": user.ID,
	    "request_id": requestID,
	})
	myLoggerWithFields.Info("User logged in")

# Cleanup

Flush any buffered logger entries before exit:

	func main() {
	    // Initialize and use logger...

	    // Before exit, sync the global logger
	    defer logger.Sync()

	    // Or sync a logger instance
	    defer myLogger.Sync()
	}

# Using Multiple Loggers in Libraries

When creating a library that might be used by applications with their own loggers:

	type MyLibrary struct {
	    logger *logger.Logger
	}

	func NewMyLibrary(opts ...Option) *MyLibrary {
	    // Default logger
	    l, _ := logger.NewLogger(logger.Config{
	        Level:       "info",
	        ServiceName: "my-library",
	    })

	    lib := &MyLibrary{
	        logger: l,
	    }

	    // Apply options
	    for _, opt := range opts {
	        opt(lib)
	    }

	    return lib
	}

	// Option for configuring the library
	type Option func(*MyLibrary)

	// WithLogger allows injecting a custom logger
	func WithLogger(logger *logger.Logger) Option {
	    return func(l *MyLibrary) {
	        l.logger = logger
	    }
	}

	// Library methods use the instance logger
	func (l *MyLibrary) DoSomething() {
	    l.logger.Info("Doing something")
	}
*/
package logger