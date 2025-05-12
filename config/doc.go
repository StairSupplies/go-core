/*
Package config provides utilities for loading application configuration from environment variables.

The package uses Viper and godotenv to create a flexible configuration system that supports
environment variables and .env files.

# Basic Usage

Define a struct with fields tagged using mapstructure:

	type AppConfig struct {
		Port        string `mapstructure:"PORT"`
		DatabaseURL string `mapstructure:"DATABASE_URL"`
		Debug       bool   `mapstructure:"DEBUG"`
	}

Then load configuration:

	// Load from .env file and environment variables
	cfg, err := config.New[AppConfig](".env")
	
	// Or simply from environment variables
	cfg, err := config.New[AppConfig]("")

# Features

- Type-safe configuration with Go generics
- Automatic binding of struct fields to environment variables
- Support for .env file loading
- Compatible with various data types including strings, integers, booleans, etc.
- Error handling for improperly formatted environment values

# Advanced Use Cases

If you have nested configuration:

	type DatabaseConfig struct {
		URL      string `mapstructure:"DATABASE_URL"`
		MaxConns int    `mapstructure:"DATABASE_MAX_CONNECTIONS"`
	}

	type AppConfig struct {
		Port     string         `mapstructure:"PORT"`
		Database DatabaseConfig
	}

Viper will attempt to map nested fields using dot notation.
*/
package config