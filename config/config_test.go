package config

import (
	"os"
	"testing"
)

// TestConfig is a test configuration struct
type TestConfig struct {
	AppName   string `mapstructure:"APP_NAME"`
	AppPort   int    `mapstructure:"APP_PORT"`
	LogLevel  string `mapstructure:"LOG_LEVEL"`
	DBUrl     string `mapstructure:"DB_URL"`
	EnableSSL bool   `mapstructure:"ENABLE_SSL"`
}

func TestNew(t *testing.T) {
	// Create a temporary .env file for testing
	envContent := `
APP_NAME=TestApp
APP_PORT=8080
LOG_LEVEL=debug
DB_URL=postgres://localhost:5432/testdb
ENABLE_SSL=true
`
	tmpEnvFile := "test.env"
	err := os.WriteFile(tmpEnvFile, []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test env file: %v", err)
	}
	defer os.Remove(tmpEnvFile)

	// Test with env file
	t.Run("With .env file", func(t *testing.T) {
		cfg, err := New[TestConfig](tmpEnvFile)
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}
		if cfg == nil {
			t.Fatal("New() returned nil config")
		}
		
		// Validate loaded values
		if cfg.AppName != "TestApp" {
			t.Errorf("Expected AppName = 'TestApp', got '%s'", cfg.AppName)
		}
		if cfg.AppPort != 8080 {
			t.Errorf("Expected AppPort = 8080, got %d", cfg.AppPort)
		}
		if cfg.LogLevel != "debug" {
			t.Errorf("Expected LogLevel = 'debug', got '%s'", cfg.LogLevel)
		}
		if cfg.DBUrl != "postgres://localhost:5432/testdb" {
			t.Errorf("Expected DBUrl = 'postgres://localhost:5432/testdb', got '%s'", cfg.DBUrl)
		}
		if !cfg.EnableSSL {
			t.Errorf("Expected EnableSSL = true, got %t", cfg.EnableSSL)
		}
	})

	// Test with environment variables
	t.Run("With environment variables", func(t *testing.T) {
		// Set environment variables
		os.Setenv("APP_NAME", "EnvApp")
		os.Setenv("APP_PORT", "9090")
		os.Setenv("LOG_LEVEL", "info")
		os.Setenv("DB_URL", "mysql://localhost:3306/envdb")
		os.Setenv("ENABLE_SSL", "false")
		defer func() {
			os.Unsetenv("APP_NAME")
			os.Unsetenv("APP_PORT")
			os.Unsetenv("LOG_LEVEL")
			os.Unsetenv("DB_URL")
			os.Unsetenv("ENABLE_SSL")
		}()

		// Non-existent env file, should use environment variables
		cfg, err := New[TestConfig]("non-existent.env")
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}
		if cfg == nil {
			t.Fatal("New() returned nil config")
		}
		
		// Validate loaded values
		if cfg.AppName != "EnvApp" {
			t.Errorf("Expected AppName = 'EnvApp', got '%s'", cfg.AppName)
		}
		if cfg.AppPort != 9090 {
			t.Errorf("Expected AppPort = 9090, got %d", cfg.AppPort)
		}
		if cfg.LogLevel != "info" {
			t.Errorf("Expected LogLevel = 'info', got '%s'", cfg.LogLevel)
		}
		if cfg.DBUrl != "mysql://localhost:3306/envdb" {
			t.Errorf("Expected DBUrl = 'mysql://localhost:3306/envdb', got '%s'", cfg.DBUrl)
		}
		if cfg.EnableSSL {
			t.Errorf("Expected EnableSSL = false, got %t", cfg.EnableSSL)
		}
	})

	// Test with invalid configuration type
	t.Run("With invalid config type", func(t *testing.T) {
		_, err := New[string]("")
		if err == nil {
			t.Error("Expected error for invalid config type, got nil")
		}
	})

	// Test with empty struct
	t.Run("With empty struct", func(t *testing.T) {
		type EmptyConfig struct{}
		cfg, err := New[EmptyConfig]("")
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}
		if cfg == nil {
			t.Fatal("New() returned nil config")
		}
	})

	// Test struct pointer
	t.Run("With struct pointer", func(t *testing.T) {
		type PointerConfig struct {
			Field string `mapstructure:"FIELD"`
		}
		os.Setenv("FIELD", "test_value")
		defer os.Unsetenv("FIELD")

		cfg, err := New[*PointerConfig]("")
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}
		if cfg == nil || *cfg == nil {
			t.Fatal("New() returned nil config")
		}

		if (*cfg).Field != "test_value" {
			t.Errorf("Expected Field = 'test_value', got '%s'", (*cfg).Field)
		}
	})

	// Test environment variables have priority over .env file
	t.Run("Environment variables have priority over .env file", func(t *testing.T) {
		// Set environment variables with different values
		os.Setenv("APP_NAME", "EnvAppOverride")
		defer os.Unsetenv("APP_NAME")

		cfg, err := New[TestConfig](tmpEnvFile)
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}
		
		// Environment variables should take precedence
		if cfg.AppName != "EnvAppOverride" {
			t.Errorf("Expected AppName = 'EnvAppOverride' from environment, got '%s'", cfg.AppName)
		}
	})

	// Test missing environment variables
	t.Run("With missing environment variables", func(t *testing.T) {
		type RequiredConfig struct {
			Required string `mapstructure:"REQUIRED"`
			Optional string `mapstructure:"OPTIONAL"`
		}

		// No environment variables set
		cfg, err := New[RequiredConfig]("")
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}

		// Empty values are allowed
		if cfg.Required != "" {
			t.Errorf("Expected Required to be empty, got '%s'", cfg.Required)
		}
		if cfg.Optional != "" {
			t.Errorf("Expected Optional to be empty, got '%s'", cfg.Optional)
		}
	})
}