/*
Package httputils provides utilities for handling HTTP responses and API error handling.

This package offers a standardized way to handle HTTP responses, API errors,
and includes middleware for consistent error handling in web applications.

# Envelope Type

The Envelope type is a map for wrapping JSON responses in a consistent structure.

	envelope := httputils.Envelope{"users": users}
	httputils.WriteJSON(w, http.StatusOK, envelope, nil)

# Error Handling

The package provides several pre-defined error types that implement Go's error interface:

	// Return a 404 error
	err := httputils.NotFoundError(errors.New("user not found"))
	return err

	// Return a 400 error
	if !validator.Valid() {
	    return httputils.BadRequestError(errors.New("validation failed"))
	}

# Response Writing

The package offers helpers for writing JSON responses:

	// Write a success response
	users := []User{...}
	meta := map[string]int{"total": len(users)}
	err := httputils.WriteSuccess(w, users, meta)

	// Write an error response
	httputils.WriteError(w, err)

# Error Middleware

The WrapHandler function simplifies error handling in HTTP handlers:

	func getUserHandler(w http.ResponseWriter, r *http.Request) error {
	    user, err := db.GetUser(id)
	    if err != nil {
	        return httputils.NotFoundError(errors.New("user not found"))
	    }
	    return httputils.WriteSuccess(w, user)
	}

	// In your router setup:
	router.Get("/api/users/:id", httputils.WrapHandler(getUserHandler))
*/
package httputils
