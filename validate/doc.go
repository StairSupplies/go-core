/*
Package validate provides input validation utilities.

It offers a fluent validation API with common validation functions
and the ability to collect multiple validation errors.

# Basic Validation

Create a validator and add validations:

    v := validate.New()
    
    // Check that a field is not empty
    v.NotBlank(user.Name, "name")
    
    // Check that a field is within length limits
    v.MaxLength(user.Bio, 200, "bio")
    v.MinLength(user.Password, 8, "password")
    
    // Check that a field matches a pattern
    v.Matches(
        user.Username, 
        regexp.MustCompile(`^[a-zA-Z0-9]+$`), 
        "username", 
        "Username can only contain letters and numbers"
    )
    
    // Check that a field is in a list of allowed values
    v.In(user.Role, []string{"admin", "user", "guest"}, "role")
    
    // Validate email format
    v.IsEmail(user.Email, "email")
    
    // Check if the validation passed
    if v.Valid() {
        // Proceed with valid data
    } else {
        // Handle validation errors
    }

# Conditional Validation

Add validation errors based on conditions:

    v.Check(len(password) >= 8, "password", "Password must be at least 8 characters")

# Custom Error Messages

Add custom error messages for fields:

    v.AddError("email", "This email is already registered")

# Converting Validation to Errors

Convert validation errors to a standard error:

    v := validate.New()
    v.NotBlank(user.Name, "name")
    // ... more validations
    
    if err := v.AsValidationError(); err != nil {
        return err
    }

# Working with Validation Errors

Check if an error is a validation error:

    if validate.IsValidationError(err) {
        // Handle validation errors specially
    }

Get validation errors from a standard error:

    if errs := validate.GetValidationErrors(err); errs != nil {
        for field, msg := range errs {
            fmt.Printf("%s: %s\n", field, msg)
        }
    }
*/
package validate