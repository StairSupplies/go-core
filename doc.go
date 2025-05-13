/*
Go-Core is a shared utility library for common Go functions and patterns.

This library provides reusable, non-business-specific utility functions across
multiple domains, organized into several packages:

# API Package

Package api provides utilities for handling HTTP responses, API error handling, and
middleware patterns.

	import "github.com/StairSupplies/go-core/api"

# Router Package

Package router provides an opinionated chi-based HTTP router with built-in middleware
for logging, error handling, and request tracing.

	import "github.com/StairSupplies/go-core/router"

# REST Package

Package rest provides a fluent REST client for making HTTP requests with built-in
JSON serialization/deserialization.

	import "github.com/StairSupplies/go-core/rest"

# JSON Utils Package

Package jsonutils provides enhanced JSON utilities for encoding and decoding with
better error handling.

	import "github.com/StairSupplies/go-core/jsonutils"

# Logger Package

Package logger provides structured logging using Uber's Zap with context-aware logging
and both structured and formatted logging options.

	import "github.com/StairSupplies/go-core/logger"

# Config Package

Package config provides utilities for loading application configuration from environment
variables with type-safe retrieval using generics.

	import "github.com/StairSupplies/go-core/config"

See the individual package documentation for more details and examples.
*/
package core