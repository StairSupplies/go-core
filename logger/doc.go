/*
Package logger provides structured logging using Uber's Zap.

It offers a global logger with context-aware logging and support for both structured
and formatted logging.

# Initialization

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

# Basic Logging

logger messages at different levels:

	// Debug level
	logger.Debug("Processing item",
	    zap.Int("item_id", item.ID),
	    zap.String("status", item.Status),
	)

	// Info level
	logger.Info("Application started",
	    zap.String("version", version),
	    zap.String("environment", env),
	)

	// Warning level
	logger.Warn("API rate limit approaching",
	    zap.Int("current", current),
	    zap.Int("limit", limit),
	)

	// Error level
	logger.Error("Failed to process payment",
	    zap.String("user_id", userID),
	    zap.Error(err),
	)

	// Fatal level (will exit the application)
	logger.Fatal("Failed to connect to database",
	    zap.Error(err),
	)

# Formatted Logging

logger with string formatting (using the sugared logger):

	// Debug level with formatting
	logger.Debugf("Processing item %d with status %s", item.ID, item.Status)

	// Info level with formatting
	logger.Infof("Application started with version %s in %s environment", version, env)

# Structured Logging

Log with structured key-value pairs (using the sugared logger):

	// Debug level with structured key-value pairs
	logger.Debugw("Processing item", 
		"item_id", item.ID, 
		"status", item.Status,
	)

	// Info level with structured key-value pairs
	logger.Infow("Application started", 
		"version", version, 
		"environment", env,
	)

	// Warning level with structured key-value pairs
	logger.Warnw("API rate limit approaching", 
		"current", current, 
		"limit", limit,
	)

	// Error level with structured key-value pairs
	logger.Errorw("Failed to process payment", 
		"user_id", userID, 
		"error", err.Error(),
	)

# Context-Aware Logging

Use context to pass loggers:

	// Add logger to context
	logger := logger.With(zap.String("request_id", requestID))
	ctx = logger.ContextWithLogger(ctx, logger)

	// Get logger from context
	func HandleRequest(ctx context.Context) {
	    logger := logger.WithContext(ctx)
	    logger.Info("Processing request")
	}

# Structured Fields

Create child loggers with additional fields:

	// Using the structured logger
	logger := logger.With(
	    zap.String("user_id", user.ID),
	    zap.String("request_id", requestID),
	)
	logger.Info("User logged in")

	// Using the sugared logger
	logger := logger.WithFields(map[string]interface{}{
	    "user_id": user.ID,
	    "request_id": requestID,
	})
	logger.Info("User logged in")

# Cleanup

Flush any buffered logger entries before exit:

	func main() {
	    // Initialize and use logger...

	    // Before exit, sync the logger
	    defer logger.Sync()
	}
*/
package logger
