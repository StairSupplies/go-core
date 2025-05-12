package config

import (
	"os"
	"testing"
)

type TestConfig struct {
	Port        string `mapstructure:"PORT"`
	DatabaseURL string `mapstructure:"DATABASE_URL"`
	Debug       bool   `mapstructure:"DEBUG"`
}

func TestNew(t *testing.T) {
	// Setup - set environment variables
	os.Setenv("PORT", "8080")
	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/testdb")
	os.Setenv("DEBUG", "true")

	// Clean up after test
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("DEBUG")
	}()

	// Test case 1: Load config from environment variables
	cfg, err := New[TestConfig]("")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Validate values were loaded correctly
	if cfg.Port != "8080" {
		t.Errorf("Expected Port to be '8080', got '%s'", cfg.Port)
	}
	if cfg.DatabaseURL != "postgres://user:pass@localhost:5432/testdb" {
		t.Errorf("Expected DatabaseURL to be 'postgres://user:pass@localhost:5432/testdb', got '%s'", cfg.DatabaseURL)
	}
	if !cfg.Debug {
		t.Errorf("Expected Debug to be true")
	}
}

func TestNewWithEnvFile(t *testing.T) {
	// Create temporary .env file
	envContent := `PORT=9090
DATABASE_URL=mysql://root:password@localhost:3306/testdb
DEBUG=false
`
	tmpFile := "./test.env"
	err := os.WriteFile(tmpFile, []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test .env file: %v", err)
	}

	// Clean up after test
	defer os.Remove(tmpFile)

	// Test case 2: Load config from .env file
	cfg, err := New[TestConfig](tmpFile)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Validate values were loaded correctly
	if cfg.Port != "9090" {
		t.Errorf("Expected Port to be '9090', got '%s'", cfg.Port)
	}
	if cfg.DatabaseURL != "mysql://root:password@localhost:3306/testdb" {
		t.Errorf("Expected DatabaseURL to be 'mysql://root:password@localhost:3306/testdb', got '%s'", cfg.DatabaseURL)
	}
	if cfg.Debug {
		t.Errorf("Expected Debug to be false")
	}
}

func TestNonStructType(t *testing.T) {
	// Test case 3: Non-struct type should return error
	_, err := New[string]("")
	if err == nil {
		t.Fatalf("Expected error for non-struct type, got nil")
	}
}

func TestStructPointer(t *testing.T) {
	// Setup
	os.Setenv("PORT", "8080")
	defer os.Unsetenv("PORT")

	// Test case 4: Struct pointer should work
	type PointerTestConfig struct {
		Port string `mapstructure:"PORT"`
	}

	cfg, err := New[PointerTestConfig]("")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.Port != "8080" {
		t.Errorf("Expected Port to be '8080', got '%s'", cfg.Port)
	}
}

func TestInvalidEnvFile(t *testing.T) {
	// Test case 5: Invalid .env file path
	_, err := New[TestConfig]("/non/existent/path")

	// Should still work, just log error and continue with env vars
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
