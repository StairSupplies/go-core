/*
Package api provides utilities for handling HTTP API responses and standardized error handling.

This package offers a consistent approach to HTTP response formatting, error handling,
and middleware for RESTful API development. It's designed to simplify common HTTP
operations while providing a structured response format.

# Response Formats

The package provides a standardized response structure for both successful and error responses:

Successful responses use the SuccessResponse struct:

	{
	  "status_code": 200,
	  "data": { ... },
	  "meta": { ... } // optional
	}

Error responses use the Error struct wrapped in an Envelope:

	{
	  "error": {
	    "status_code": 400,
	    "message": "Bad request: invalid user ID"
	  }
	}

# Basic Usage

Writing successful responses:

	// Return data with status 200
	users := []User{...}
	meta := map[string]int{"total": len(users)}
	api.WriteSuccess(w, users, meta)

	// Custom response with Envelope
	envelope := api.Envelope{"users": users, "count": len(users)}
	api.WriteJSON(w, http.StatusOK, envelope, nil)

# Error Handling

The package provides error constructors for common HTTP error codes:

	// Return a 404 error
	if user == nil {
	    return api.NotFoundError(errors.New("user not found"))
	}

	// Return a 400 error
	if !validator.Valid() {
	    return api.BadRequestError(errors.New("validation failed"))
	}

	// Return a 401 error
	if !authenticated {
	    return api.UnauthorizedError(errors.New("invalid credentials"))
	}

# Handler Functions

The package defines the HandlerFunc type that returns an error instead of directly
handling it, allowing for cleaner controller logic:

	func getUserHandler(w http.ResponseWriter, r *http.Request) error {
	    user, err := db.GetUser(id)
	    if err != nil {
	        return api.NotFoundError(errors.New("user not found"))
	    }
	    return api.WriteSuccess(w, user)
	}

	// In your router setup, wrap the handler:
	router.Get("/api/users/{id}", api.WrapHandler(getUserHandler))

The WrapHandler function automatically logs errors and writes appropriate error responses.

# Integration with Router

This package works seamlessly with the router package, which provides additional
middleware and routing capabilities built on the chi router.
*/
package api