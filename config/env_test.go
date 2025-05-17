package config

import (
	"testing"
)

func TestIsValidEnvironment(t *testing.T) {
	tests := []struct {
		name string
		env  Environment
		want bool
	}{
		{"DevelopmentValid", EnvDevelopment, true},
		{"StagingValid", EnvStaging, true},
		{"ProductionValid", EnvProduction, true},
		{"EmptyInvalid", "", false},
		{"UnknownInvalid", "test", false},
		{"CaseSensitive", "PRODUCTION", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidEnvironment(tt.env)
			if got != tt.want {
				t.Errorf("IsValidEnvironment(%q) = %v, want %v", tt.env, got, tt.want)
			}
		})
	}
}

func TestGetDefaultEnvironment(t *testing.T) {
	if got := GetDefaultEnvironment(); got != EnvDevelopment {
		t.Errorf("GetDefaultEnvironment() = %v, want %v", got, EnvDevelopment)
	}
}

func TestEnvironmentConstantValues(t *testing.T) {
	// Ensure environment constants have the expected string values
	testCases := []struct {
		env      Environment
		expected string
	}{
		{EnvDevelopment, "development"},
		{EnvStaging, "staging"},
		{EnvProduction, "production"},
	}

	for _, tc := range testCases {
		t.Run(string(tc.env), func(t *testing.T) {
			if string(tc.env) != tc.expected {
				t.Errorf("Environment %v has value %q, want %q", tc.env, tc.env, tc.expected)
			}
		})
	}
}

func TestEnvironmentString(t *testing.T) {
	// Ensure Environment type correctly converts to string
	testCases := []struct {
		env      Environment
		expected string
	}{
		{EnvDevelopment, "development"},
		{EnvStaging, "staging"},
		{EnvProduction, "production"},
	}

	for _, tc := range testCases {
		t.Run(string(tc.env), func(t *testing.T) {
			if got := string(tc.env); got != tc.expected {
				t.Errorf("string(%v) = %q, want %q", tc.env, got, tc.expected)
			}
		})
	}
}