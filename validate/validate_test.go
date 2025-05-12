package validate

import (
	"errors"
	"regexp"
	"strings"
	"testing"
)

func TestValidator_Valid(t *testing.T) {
	t.Run("empty validator is valid", func(t *testing.T) {
		v := New()
		if !v.Valid() {
			t.Error("Expected empty validator to be valid")
		}
	})

	t.Run("validator with errors is invalid", func(t *testing.T) {
		v := New()
		v.AddError("field", "error message")
		if v.Valid() {
			t.Error("Expected validator with errors to be invalid")
		}
	})
}

func TestValidator_AddError(t *testing.T) {
	t.Run("add error to empty validator", func(t *testing.T) {
		v := New()
		v.AddError("field", "error message")
		
		if len(v.Errors) != 1 {
			t.Errorf("Expected 1 error, got %d", len(v.Errors))
		}
		
		if v.Errors["field"] != "error message" {
			t.Errorf("Expected error message 'error message', got '%s'", v.Errors["field"])
		}
	})
	
	t.Run("adding duplicate field does not overwrite", func(t *testing.T) {
		v := New()
		v.AddError("field", "first error")
		v.AddError("field", "second error")
		
		if len(v.Errors) != 1 {
			t.Errorf("Expected 1 error, got %d", len(v.Errors))
		}
		
		if v.Errors["field"] != "first error" {
			t.Errorf("Expected error message 'first error', got '%s'", v.Errors["field"])
		}
	})
}

func TestValidator_Check(t *testing.T) {
	t.Run("check passes", func(t *testing.T) {
		v := New()
		v.Check(true, "field", "error message")
		
		if len(v.Errors) != 0 {
			t.Errorf("Expected 0 errors, got %d", len(v.Errors))
		}
	})
	
	t.Run("check fails", func(t *testing.T) {
		v := New()
		v.Check(false, "field", "error message")
		
		if len(v.Errors) != 1 {
			t.Errorf("Expected 1 error, got %d", len(v.Errors))
		}
		
		if v.Errors["field"] != "error message" {
			t.Errorf("Expected error message 'error message', got '%s'", v.Errors["field"])
		}
	})
}

func TestValidator_NotBlank(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"empty string", "", false},
		{"whitespace", "   ", false},
		{"tabs and newlines", "\t\n", false},
		{"valid string", "hello", true},
		{"string with spaces", "  hello  ", true},
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := New()
			v.NotBlank(tc.value, "field")
			
			if tc.expected && len(v.Errors) != 0 {
				t.Errorf("Expected no errors for '%s', got %v", tc.value, v.Errors)
			}
			
			if !tc.expected && len(v.Errors) != 1 {
				t.Errorf("Expected 1 error for '%s', got %d", tc.value, len(v.Errors))
			}
		})
	}
}

func TestValidator_MaxLength(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		maxLen   int
		expected bool
	}{
		{"empty string", "", 10, true},
		{"string at max length", "1234567890", 10, true},
		{"string under max length", "12345", 10, true},
		{"string over max length", "12345678901", 10, false},
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := New()
			v.MaxLength(tc.value, tc.maxLen, "field")
			
			if tc.expected && len(v.Errors) != 0 {
				t.Errorf("Expected no errors for '%s' with max length %d, got %v", tc.value, tc.maxLen, v.Errors)
			}
			
			if !tc.expected && len(v.Errors) != 1 {
				t.Errorf("Expected 1 error for '%s' with max length %d, got %d", tc.value, tc.maxLen, len(v.Errors))
			}
		})
	}
}

func TestValidator_MinLength(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		minLen   int
		expected bool
	}{
		{"empty string", "", 1, false},
		{"string at min length", "12345", 5, true},
		{"string over min length", "123456", 5, true},
		{"string under min length", "1234", 5, false},
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := New()
			v.MinLength(tc.value, tc.minLen, "field")
			
			if tc.expected && len(v.Errors) != 0 {
				t.Errorf("Expected no errors for '%s' with min length %d, got %v", tc.value, tc.minLen, v.Errors)
			}
			
			if !tc.expected && len(v.Errors) != 1 {
				t.Errorf("Expected 1 error for '%s' with min length %d, got %d", tc.value, tc.minLen, len(v.Errors))
			}
		})
	}
}

func TestValidator_Matches(t *testing.T) {
	alphaPattern := regexp.MustCompile(`^[a-zA-Z]+$`)
	
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"empty string", "", false},
		{"alphabetic string", "hello", true},
		{"alphanumeric string", "hello123", false},
		{"string with symbols", "hello!", false},
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := New()
			v.Matches(tc.value, alphaPattern, "field", "must contain only letters")
			
			if tc.expected && len(v.Errors) != 0 {
				t.Errorf("Expected no errors for '%s', got %v", tc.value, v.Errors)
			}
			
			if !tc.expected && len(v.Errors) != 1 {
				t.Errorf("Expected 1 error for '%s', got %d", tc.value, len(v.Errors))
			}
		})
	}
}

func TestValidator_In(t *testing.T) {
	allowed := []string{"apple", "banana", "cherry"}
	
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"empty string", "", false},
		{"allowed value", "apple", true},
		{"not allowed value", "orange", false},
		{"case sensitive", "Apple", false},
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := New()
			v.In(tc.value, allowed, "field")
			
			if tc.expected && len(v.Errors) != 0 {
				t.Errorf("Expected no errors for '%s', got %v", tc.value, v.Errors)
			}
			
			if !tc.expected && len(v.Errors) != 1 {
				t.Errorf("Expected 1 error for '%s', got %d", tc.value, len(v.Errors))
			}
		})
	}
}

func TestValidator_IsEmail(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"empty string", "", false},
		{"valid email", "test@example.com", true},
		{"email with subdomain", "test@sub.example.com", true},
		{"email with plus", "test+tag@example.com", true},
		{"email with dots", "first.last@example.com", true},
		{"no @", "testexample.com", false},
		{"no domain", "test@", false},
		{"invalid tld", "test@example", false},
		{"multiple @", "test@@example.com", false},
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := New()
			v.IsEmail(tc.value, "email")
			
			if tc.expected && len(v.Errors) != 0 {
				t.Errorf("Expected no errors for '%s', got %v", tc.value, v.Errors)
			}
			
			if !tc.expected && len(v.Errors) != 1 {
				t.Errorf("Expected 1 error for '%s', got %d", tc.value, len(v.Errors))
			}
		})
	}
}

func TestValidationError_Error(t *testing.T) {
	tests := []struct {
		name     string
		errors   map[string]string
		expected string
	}{
		{
			name:     "single error",
			errors:   map[string]string{"field": "error message"},
			expected: "validation failed: field: error message",
		},
		{
			name: "multiple errors",
			errors: map[string]string{
				"field1": "error message 1",
				"field2": "error message 2",
			},
			expected: "validation failed: field1: error message 1; field2: error message 2",
		},
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ve := ValidationError{Errors: tc.errors}
			
			errorStr := ve.Error()
			
			// Since the order of map iteration is not guaranteed, we need to check
			// that the error string contains all the expected parts, rather than
			// checking for an exact match.
			
			if !strings.Contains(errorStr, "validation failed") {
				t.Errorf("Error string does not contain 'validation failed': %s", errorStr)
			}
			
			for field, msg := range tc.errors {
				if !strings.Contains(errorStr, field+": "+msg) {
					t.Errorf("Error string does not contain '%s: %s': %s", field, msg, errorStr)
				}
			}
		})
	}
}

func TestValidator_AsValidationError(t *testing.T) {
	t.Run("valid validator returns nil", func(t *testing.T) {
		v := New()
		err := v.AsValidationError()
		
		if err != nil {
			t.Errorf("Expected nil error for valid validator, got %v", err)
		}
	})
	
	t.Run("invalid validator returns ValidationError", func(t *testing.T) {
		v := New()
		v.AddError("field", "error message")
		err := v.AsValidationError()
		
		if err == nil {
			t.Fatal("Expected ValidationError for invalid validator, got nil")
		}
		
		ve, ok := err.(ValidationError)
		if !ok {
			t.Fatalf("Expected ValidationError, got %T", err)
		}
		
		if len(ve.Errors) != 1 {
			t.Errorf("Expected 1 error in ValidationError, got %d", len(ve.Errors))
		}
		
		if ve.Errors["field"] != "error message" {
			t.Errorf("Expected error message 'error message', got '%s'", ve.Errors["field"])
		}
	})
}

func TestIsValidationError(t *testing.T) {
	t.Run("ValidationError", func(t *testing.T) {
		ve := ValidationError{Errors: map[string]string{"field": "error"}}
		
		if !IsValidationError(ve) {
			t.Errorf("Expected IsValidationError to return true for ValidationError")
		}
	})
	
	t.Run("error wrapping ValidationError", func(t *testing.T) {
		ve := ValidationError{Errors: map[string]string{"field": "error"}}
		wrappedErr := errors.New("wrapped: " + ve.Error())
		
		// This should fail because we're not using errors.Wrap
		if IsValidationError(wrappedErr) {
			t.Errorf("Expected IsValidationError to return false for string-wrapped error")
		}
	})
	
	t.Run("other error", func(t *testing.T) {
		err := errors.New("some other error")
		
		if IsValidationError(err) {
			t.Errorf("Expected IsValidationError to return false for non-ValidationError")
		}
	})
}

func TestGetValidationErrors(t *testing.T) {
	t.Run("ValidationError", func(t *testing.T) {
		ve := ValidationError{Errors: map[string]string{"field": "error"}}
		
		errors := GetValidationErrors(ve)
		
		if len(errors) != 1 {
			t.Errorf("Expected 1 error in ValidationError, got %d", len(errors))
		}
		
		if errors["field"] != "error" {
			t.Errorf("Expected error message 'error', got '%s'", errors["field"])
		}
	})
	
	t.Run("other error", func(t *testing.T) {
		err := errors.New("some other error")
		
		errors := GetValidationErrors(err)
		
		if errors != nil {
			t.Errorf("Expected nil for non-ValidationError, got %v", errors)
		}
	})
}