package validate

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Validator provides methods for validating form data
type Validator struct {
	Errors map[string]string
}

// New creates a new validator
func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

// Valid returns true if there are no errors
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds an error message for a given field
func (v *Validator) AddError(field, message string) {
	if _, exists := v.Errors[field]; !exists {
		v.Errors[field] = message
	}
}

// Check adds an error message for a given field if the condition is false
func (v *Validator) Check(ok bool, field, message string) {
	if !ok {
		v.AddError(field, message)
	}
}

// NotBlank checks that a value is not empty
func (v *Validator) NotBlank(value string, field string) {
	if strings.TrimSpace(value) == "" {
		v.AddError(field, "This field cannot be blank")
	}
}

// MaxLength checks that a value is not longer than n characters
func (v *Validator) MaxLength(value string, n int, field string) {
	if len(value) > n {
		v.AddError(field, fmt.Sprintf("This field cannot be more than %d characters", n))
	}
}

// MinLength checks that a value is at least n characters
func (v *Validator) MinLength(value string, n int, field string) {
	if len(value) < n {
		v.AddError(field, fmt.Sprintf("This field must be at least %d characters", n))
	}
}

// Matches checks that a value matches a regular expression pattern
func (v *Validator) Matches(value string, pattern *regexp.Regexp, field, message string) {
	if !pattern.MatchString(value) {
		v.AddError(field, message)
	}
}

// In checks that a value is in a list of permitted values
func (v *Validator) In(value string, list []string, field string) {
	for _, item := range list {
		if value == item {
			return
		}
	}
	v.AddError(field, "This field contains an invalid value")
}

// IsEmail checks that a value is a valid email address
var EmailRX = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func (v *Validator) IsEmail(value string, field string) {
	if !EmailRX.MatchString(value) {
		v.AddError(field, "This field must be a valid email address")
	}
}

// ValidationError represents an error from validation
type ValidationError struct {
	Errors map[string]string
}

func (ve ValidationError) Error() string {
	var sb strings.Builder
	sb.WriteString("validation failed: ")
	
	for field, err := range ve.Errors {
		sb.WriteString(fmt.Sprintf("%s: %s; ", field, err))
	}
	
	return strings.TrimSuffix(sb.String(), "; ")
}

// AsValidationError converts a validator into an error if it has errors
func (v *Validator) AsValidationError() error {
	if v.Valid() {
		return nil
	}
	
	return ValidationError{Errors: v.Errors}
}

// IsValidationError checks if an error is a ValidationError
func IsValidationError(err error) bool {
	var ve ValidationError
	return errors.As(err, &ve)
}

// GetValidationErrors extracts the validation errors from an error
func GetValidationErrors(err error) map[string]string {
	var ve ValidationError
	if errors.As(err, &ve) {
		return ve.Errors
	}
	return nil
}