package logger

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestInit(t *testing.T) {
	// Test with default config
	t.Run("default config", func(t *testing.T) {
		cfg := Config{
			Level:       "info",
			Development: false,
			ServiceName: "test-service",
			InitialFields: map[string]interface{}{
				"version": "1.0.0",
			},
		}

		err := Init(cfg)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Check that logger is initialized
		logger := GetLogger()
		if logger == nil {
			t.Fatal("Expected logger to be initialized")
		}

		// Clean up
		_ = Sync()
	})

	// Test with invalid log level
	t.Run("invalid log level", func(t *testing.T) {
		cfg := Config{
			Level: "invalid",
		}

		err := Init(cfg)
		if err == nil {
			t.Fatal("Expected error for invalid log level, got nil")
		}
	})

	// Test with development mode
	t.Run("development mode", func(t *testing.T) {
		cfg := Config{
			Level:       "debug",
			Development: true,
		}

		err := Init(cfg)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Clean up
		_ = Sync()
	})
}

// Helper function to capture logs during tests
func captureOutput(f func()) ([]map[string]interface{}, error) {
	// Create a pipe
	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	// Save original outputs
	originalStdout := os.Stdout
	originalStderr := os.Stderr

	// Set output to pipe
	os.Stdout = w
	os.Stderr = w

	// Create logger with stdout output
	cfg := Config{
		Level:       "debug",
		Development: false,
		OutputPaths: []string{"stdout"},
	}

	err = Init(cfg)
	if err != nil {
		return nil, err
	}

	// Execute the function that logs
	f()

	// Sync to flush logs
	_ = Sync()

	// Restore original outputs
	os.Stdout = originalStdout
	os.Stderr = originalStderr

	// Close the write end of the pipe
	err = w.Close()
	if err != nil {
		return nil, err
	}

	// Read all logs from the pipe
	output, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// Split logs by newline
	lines := splitLines(string(output))

	// Parse JSON logs
	logs := make([]map[string]interface{}, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			continue
		}

		var log map[string]interface{}
		err = json.Unmarshal([]byte(line), &log)
		if err != nil {
			return nil, err
		}

		logs = append(logs, log)
	}

	return logs, nil
}

// Helper function to split output by lines
func splitLines(s string) []string {
	var lines []string
	var line string

	for _, char := range s {
		if char == '\n' {
			lines = append(lines, line)
			line = ""
		} else {
			line += string(char)
		}
	}

	if line != "" {
		lines = append(lines, line)
	}

	return lines
}

func TestLoggerMethods(t *testing.T) {
	// Create a test logger with an observer core
	core, recorded := observer.New(zapcore.InfoLevel)
	globalLogger = zap.New(core)
	globalSugared = globalLogger.Sugar()

	// Test all log level methods
	t.Run("Debug", func(t *testing.T) {
		recorded.TakeAll() // Clear previous logs
		Debug("debug message", zap.String("key", "value"))
		logs := recorded.TakeAll()
		if len(logs) > 0 {
			t.Error("Debug logs should not be recorded at Info level")
		}
	})

	t.Run("Info", func(t *testing.T) {
		Info("info message", zap.String("key", "value"))
		logs := recorded.TakeAll()
		if len(logs) != 1 {
			t.Fatalf("Expected 1 log, got %d", len(logs))
		}
		if logs[0].Message != "info message" {
			t.Errorf("Expected message 'info message', got '%s'", logs[0].Message)
		}
		if logs[0].Context[0].Key != "key" || logs[0].Context[0].String != "value" {
			t.Errorf("Expected context key 'key' with value 'value'")
		}
	})

	t.Run("Warn", func(t *testing.T) {
		Warn("warn message", zap.String("key", "value"))
		logs := recorded.TakeAll()
		if len(logs) != 1 {
			t.Fatalf("Expected 1 log, got %d", len(logs))
		}
		if logs[0].Message != "warn message" {
			t.Errorf("Expected message 'warn message', got '%s'", logs[0].Message)
		}
	})

	t.Run("Error", func(t *testing.T) {
		Error("error message", zap.String("key", "value"))
		logs := recorded.TakeAll()
		if len(logs) != 1 {
			t.Fatalf("Expected 1 log, got %d", len(logs))
		}
		if logs[0].Message != "error message" {
			t.Errorf("Expected message 'error message', got '%s'", logs[0].Message)
		}
	})

	// We don't test Fatal because it calls os.Exit(1)

	// Test format methods
	t.Run("Infof", func(t *testing.T) {
		Infof("info %s", "message")
		logs := recorded.TakeAll()
		if len(logs) != 1 {
			t.Fatalf("Expected 1 log, got %d", len(logs))
		}
		if logs[0].Message != "info message" {
			t.Errorf("Expected message 'info message', got '%s'", logs[0].Message)
		}
	})

	t.Run("Warnf", func(t *testing.T) {
		Warnf("warn %s", "message")
		logs := recorded.TakeAll()
		if len(logs) != 1 {
			t.Fatalf("Expected 1 log, got %d", len(logs))
		}
		if logs[0].Message != "warn message" {
			t.Errorf("Expected message 'warn message', got '%s'", logs[0].Message)
		}
	})

	t.Run("Errorf", func(t *testing.T) {
		Errorf("error %s", "message")
		logs := recorded.TakeAll()
		if len(logs) != 1 {
			t.Fatalf("Expected 1 log, got %d", len(logs))
		}
		if logs[0].Message != "error message" {
			t.Errorf("Expected message 'error message', got '%s'", logs[0].Message)
		}
	})

	// Test structured logging methods (*w methods)
	t.Run("Debugw", func(t *testing.T) {
		recorded.TakeAll() // Clear previous logs
		Debugw("debug message", "key", "value")
		logs := recorded.TakeAll()
		if len(logs) > 0 {
			t.Error("Debug logs should not be recorded at Info level")
		}
	})

	t.Run("Infow", func(t *testing.T) {
		Infow("info message", "key", "value")
		logs := recorded.TakeAll()
		if len(logs) != 1 {
			t.Fatalf("Expected 1 log, got %d", len(logs))
		}
		if logs[0].Message != "info message" {
			t.Errorf("Expected message 'info message', got '%s'", logs[0].Message)
		}
		// Check field was added
		hasField := false
		for _, field := range logs[0].Context {
			if field.Key == "key" && field.String == "value" {
				hasField = true
				break
			}
		}
		if !hasField {
			t.Error("Expected log to include the field 'key' with value 'value'")
		}
	})

	t.Run("Warnw", func(t *testing.T) {
		Warnw("warn message", "key", "value")
		logs := recorded.TakeAll()
		if len(logs) != 1 {
			t.Fatalf("Expected 1 log, got %d", len(logs))
		}
		if logs[0].Message != "warn message" {
			t.Errorf("Expected message 'warn message', got '%s'", logs[0].Message)
		}
		// Check field was added
		hasField := false
		for _, field := range logs[0].Context {
			if field.Key == "key" && field.String == "value" {
				hasField = true
				break
			}
		}
		if !hasField {
			t.Error("Expected log to include the field 'key' with value 'value'")
		}
	})

	t.Run("Errorw", func(t *testing.T) {
		Errorw("error message", "key", "value")
		logs := recorded.TakeAll()
		if len(logs) != 1 {
			t.Fatalf("Expected 1 log, got %d", len(logs))
		}
		if logs[0].Message != "error message" {
			t.Errorf("Expected message 'error message', got '%s'", logs[0].Message)
		}
		// Check field was added
		hasField := false
		for _, field := range logs[0].Context {
			if field.Key == "key" && field.String == "value" {
				hasField = true
				break
			}
		}
		if !hasField {
			t.Error("Expected log to include the field 'key' with value 'value'")
		}
	})

	// We don't test Fatalw because it calls os.Exit(1)
}

func TestWithContext(t *testing.T) {
	// Create a test logger
	logger := zap.NewExample()

	// Create a context with the logger
	ctx := context.Background()
	ctx = ContextWithLogger(ctx, logger)

	// Get the logger from the context
	loggerFromCtx := WithContext(ctx)

	// Verify that it's the same logger
	if loggerFromCtx != logger {
		t.Error("Expected logger from context to be the same as the original logger")
	}

	// Test with a context that doesn't have a logger
	emptyCtx := context.Background()
	loggerFromEmptyCtx := WithContext(emptyCtx)

	// Verify that it returns the global logger
	if loggerFromEmptyCtx != GetLogger() {
		t.Error("Expected logger from empty context to be the global logger")
	}
}

func TestWithFields(t *testing.T) {
	// Create a test logger with an observer core
	core, recorded := observer.New(zapcore.InfoLevel)
	globalLogger = zap.New(core)
	globalSugared = globalLogger.Sugar()

	// Create a logger with fields
	fields := map[string]interface{}{
		"service": "test",
		"version": "1.0.0",
	}

	logger := WithFields(fields)

	// Log a message
	logger.Info("test message")

	// Check the logs
	logs := recorded.TakeAll()
	if len(logs) != 1 {
		t.Fatalf("Expected 1 log, got %d", len(logs))
	}

	// Check that the fields are included
	hasService := false
	hasVersion := false

	for _, field := range logs[0].Context {
		if field.Key == "service" && field.String == "test" {
			hasService = true
		}
		if field.Key == "version" && field.String == "1.0.0" {
			hasVersion = true
		}
	}

	if !hasService {
		t.Error("Expected log to include service field")
	}
	if !hasVersion {
		t.Error("Expected log to include version field")
	}
}

func TestWith(t *testing.T) {
	// Create a test logger with an observer core
	core, recorded := observer.New(zapcore.InfoLevel)
	globalLogger = zap.New(core)
	globalSugared = globalLogger.Sugar()

	// Create a logger with fields
	logger := With(zap.String("service", "test"), zap.String("version", "1.0.0"))

	// Log a message
	logger.Info("test message")

	// Check the logs
	logs := recorded.TakeAll()
	if len(logs) != 1 {
		t.Fatalf("Expected 1 log, got %d", len(logs))
	}

	// Check that the fields are included
	hasService := false
	hasVersion := false

	for _, field := range logs[0].Context {
		if field.Key == "service" && field.String == "test" {
			hasService = true
		}
		if field.Key == "version" && field.String == "1.0.0" {
			hasVersion = true
		}
	}

	if !hasService {
		t.Error("Expected log to include service field")
	}
	if !hasVersion {
		t.Error("Expected log to include version field")
	}
}

func TestGetLogger(t *testing.T) {
	// Reset global logger
	globalLogger = nil
	globalSugared = nil

	// Get logger should create a default logger if none exists
	logger := GetLogger()
	if logger == nil {
		t.Fatal("Expected GetLogger to return a logger")
	}

	// Get sugared logger should create a default logger if none exists
	sugared := GetSugared()
	if sugared == nil {
		t.Fatal("Expected GetSugared to return a logger")
	}
}
