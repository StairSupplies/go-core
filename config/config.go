package config

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// New creates a new configuration instance of type T by loading environment variables
// and/or a .env file. The configuration struct should use mapstructure tags to define
// which environment variables to bind to each field.
//
// Example:
//
//	type AppConfig struct {
//	    Port        string `mapstructure:"PORT"`
//	    DatabaseURL string `mapstructure:"DATABASE_URL"`
//	    Debug       bool   `mapstructure:"DEBUG"`
//	}
//
//	cfg, err := config.New[AppConfig](".env")
func New[T any](path string) (*T, error) {
	// Load .env file if available
	err := godotenv.Load(path)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Error loading .env file: %v", err)
	}

	// Configure viper
	viper.AutomaticEnv()

	// Create an instance of the type to inspect its fields
	var cfg T
	t := reflect.TypeOf(cfg)

	// Handle both struct types and pointers to struct types
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Ensure we're working with a struct
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("config type must be a struct or pointer to struct")
	}

	// Bind each field with a mapstructure tag to its environment variable
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		envVar := field.Tag.Get("mapstructure")
		if envVar != "" {
			err = viper.BindEnv(envVar, envVar)
			if err != nil {
				return nil, fmt.Errorf("failed to bind environment variable %s: %w", envVar, err)
			}
		}
	}

	// Unmarshal the configuration
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	return &cfg, nil
}
