/*
Package config provides utilities for loading application configuration from environment variables and .env files.

# Overview

The config package simplifies the process of loading and managing configuration settings
in Go applications. It uses environment variables as the primary source of configuration,
with optional support for loading from .env files.

# Features

  - Type-safe configuration using generics
  - Automatic binding of environment variables to struct fields
  - Support for loading from .env files using godotenv
  - Environment constants for standard deployment environments

# Usage

Define a configuration struct with mapstructure tags to specify which environment variables
to bind to each field:

	type AppConfig struct {
		AppName  string `mapstructure:"APP_NAME"`
		Port     int    `mapstructure:"APP_PORT"`
		LogLevel string `mapstructure:"LOG_LEVEL"`
		Debug    bool   `mapstructure:"APP_DEBUG"`
	}

Then load the configuration:

	// Load from environment variables only
	cfg, err := config.New[AppConfig]("")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Or load from .env file with fallback to environment variables
	cfg, err := config.New[AppConfig](".env")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

# Environment Management

The package provides constants for standard deployment environments:

	- config.EnvDevelopment - Local development environment
	- config.EnvStaging - Staging/test environment
	- config.EnvProduction - Production environment

Helper functions for environment management:

	// Check if an environment is valid
	isValid := config.IsValidEnvironment(env)

	// Get the default environment (development)
	defaultEnv := config.GetDefaultEnvironment()

# Priority Order

When loading configuration, environment variables take precedence over values defined
in .env files. This allows for easy overriding of configuration values in different
deployment environments.

# Dependencies

This package uses:
  - github.com/joho/godotenv for loading .env files
  - github.com/spf13/viper for binding environment variables to struct fields
*/
package config