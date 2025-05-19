package config_test

import (
	"fmt"
	"os"

	"github.com/StairSupplies/go-core/config"
)

func Example() {
	// Define your configuration structure
	type AppConfig struct {
		AppName string `mapstructure:"APP_NAME"`
		Port    int    `mapstructure:"APP_PORT"`
		Debug   bool   `mapstructure:"APP_DEBUG"`
		Env     string `mapstructure:"APP_ENV"`
	}

	// Set environment variables for the example
	os.Setenv("APP_NAME", "ExampleApp")
	os.Setenv("APP_PORT", "3000")
	os.Setenv("APP_DEBUG", "true")
	os.Setenv("APP_ENV", "development")

	// Load the configuration
	cfg, err := config.New[AppConfig]("")
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		return
	}

	// Use the configuration values
	fmt.Printf("App Name: %s\n", cfg.AppName)
	fmt.Printf("Port: %d\n", cfg.Port)
	fmt.Printf("Debug Mode: %v\n", cfg.Debug)
	fmt.Printf("Environment: %s\n", cfg.Env)

	// Cleanup
	os.Unsetenv("APP_NAME")
	os.Unsetenv("APP_PORT")
	os.Unsetenv("APP_DEBUG")
	os.Unsetenv("APP_ENV")

	// Output:
	// App Name: ExampleApp
	// Port: 3000
	// Debug Mode: true
	// Environment: development
}

func ExampleIsValidEnvironment() {
	// Check if environments are valid
	fmt.Printf("Is 'development' valid? %v\n", config.IsValidEnvironment(config.EnvDevelopment))
	fmt.Printf("Is 'staging' valid? %v\n", config.IsValidEnvironment(config.EnvStaging))
	fmt.Printf("Is 'production' valid? %v\n", config.IsValidEnvironment(config.EnvProduction))
	fmt.Printf("Is 'test' valid? %v\n", config.IsValidEnvironment("test"))

	// Output:
	// Is 'development' valid? true
	// Is 'staging' valid? true
	// Is 'production' valid? true
	// Is 'test' valid? false
}

func ExampleGetDefaultEnvironment() {
	// Get the default environment
	defaultEnv := config.GetDefaultEnvironment()
	fmt.Printf("Default environment: %s\n", defaultEnv)

	// Output:
	// Default environment: development
}

func ExampleNew_withEnvFile() {
	// This example demonstrates loading configuration from a .env file
	// For testing purposes, we'll simulate this with environment variables

	// Define a configuration struct
	type DatabaseConfig struct {
		Host     string `mapstructure:"DB_HOST"`
		Port     int    `mapstructure:"DB_PORT"`
		Name     string `mapstructure:"DB_NAME"`
		User     string `mapstructure:"DB_USER"`
		Password string `mapstructure:"DB_PASSWORD"`
	}

	// In a real application, these would be defined in a .env file
	// DB_HOST=localhost
	// DB_PORT=5432
	// DB_NAME=myapp
	// DB_USER=dbuser
	// DB_PASSWORD=secret

	// For the example, we'll set environment variables directly
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "myapp")
	os.Setenv("DB_USER", "dbuser")
	os.Setenv("DB_PASSWORD", "secret")

	// Load configuration
	// In a real app, you'd provide the path to your .env file
	dbConfig, err := config.New[DatabaseConfig]("")
	if err != nil {
		fmt.Printf("Error loading database configuration: %v\n", err)
		return
	}

	// Use the configuration
	fmt.Printf("Database connection string: %s:%d/%s\n", 
		dbConfig.Host, dbConfig.Port, dbConfig.Name)

	// Cleanup
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")

	// Output:
	// Database connection string: localhost:5432/myapp
}