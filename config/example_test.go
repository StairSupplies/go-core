package config_test

import (
	"fmt"
	"os"

	"github.com/StairSupplies/go-core/config"
)

// ExampleNew demonstrates how to define and load a configuration structure
// from environment variables.
func ExampleNew() {
	// Define a configuration structure
	type AppConfig struct {
		ServerPort  string `mapstructure:"SERVER_PORT"`
		DatabaseURL string `mapstructure:"DATABASE_URL"`
		Debug       bool   `mapstructure:"DEBUG"`
	}

	// Set environment variables for testing
	os.Setenv("SERVER_PORT", "8080")
	os.Setenv("DATABASE_URL", "postgres://user:password@localhost:5432/mydb")
	os.Setenv("DEBUG", "true")

	// Load configuration
	cfg, err := config.New[AppConfig]("")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	// Use the configuration
	fmt.Printf("Server will run on port: %s\n", cfg.ServerPort)
	fmt.Printf("Database URL: %s\n", cfg.DatabaseURL)
	fmt.Printf("Debug mode: %v\n", cfg.Debug)

	// Clean up environment variables
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("DEBUG")

	// Output:
	// Server will run on port: 8080
	// Database URL: postgres://user:password@localhost:5432/mydb
	// Debug mode: true
}

// ExampleNew_withDefault demonstrates how to provide default values in your
// configuration structure.
func ExampleNew_withDefault() {
	// Define a configuration structure with default values
	// Note: Viper doesn't directly support default values via struct tags,
	// so we'd check for zero values after loading the config.
	type ServerConfig struct {
		Port    string `mapstructure:"PORT"`
		Host    string `mapstructure:"HOST"`
		Timeout int    `mapstructure:"TIMEOUT"`
	}

	// Set only some environment variables
	os.Setenv("PORT", "3000")
	// HOST is intentionally not set

	// Load configuration
	cfg, err := config.New[ServerConfig]("")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	// Apply defaults for missing values
	if cfg.Host == "" {
		cfg.Host = "localhost"
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 // Default timeout in seconds
	}

	// Use the configuration with defaults applied
	fmt.Printf("Server address: %s:%s\n", cfg.Host, cfg.Port)
	fmt.Printf("Timeout: %d seconds\n", cfg.Timeout)

	// Clean up environment variables
	os.Unsetenv("PORT")

	// Output:
	// Server address: localhost:3000
	// Timeout: 30 seconds
}

// ExampleNew_withEnvFile demonstrates how you could load config from a .env file.
// Note: In tests, we're simulating with direct environment variables instead
// of actually creating a .env file.
func ExampleNew_withEnvFile() {
	// Define a configuration structure
	type LogConfig struct {
		Level string `mapstructure:"LOG_LEVEL"`
		Path  string `mapstructure:"LOG_PATH"`
	}

	// Normally these would be loaded from a .env file
	// but for test purposes we set them directly
	os.Setenv("LOG_LEVEL", "info")
	os.Setenv("LOG_PATH", "/var/log/app.log")

	// Load configuration (normally this would include a path to .env file)
	cfg, err := config.New[LogConfig]("")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	// Use the configuration
	fmt.Printf("Log level: %s\n", cfg.Level)
	fmt.Printf("Log file: %s\n", cfg.Path)

	// Clean up environment variables
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("LOG_PATH")

	// Output:
	// Log level: info
	// Log file: /var/log/app.log
}