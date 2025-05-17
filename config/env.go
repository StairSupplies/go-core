package config

// Environment represents the environment in which the application is running.
type Environment string

// Environment constants define the standard environments for applications.
const (
	// EnvDevelopment represents a local development environment.
	EnvDevelopment Environment = "development"

	// EnvStaging represents a staging/test environment.
	EnvStaging Environment = "staging"

	// EnvProduction represents a production environment.
	EnvProduction Environment = "production"
)

// IsValidEnvironment checks if the provided environment string is one of the standard environments.
func IsValidEnvironment(env Environment) bool {
	switch env {
	case EnvDevelopment, EnvStaging, EnvProduction:
		return true
	default:
		return false
	}
}

// GetDefaultEnvironment returns the default environment (development).
func GetDefaultEnvironment() Environment {
	return EnvDevelopment
}
