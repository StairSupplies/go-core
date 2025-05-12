/*
Go-Core is a shared utility library for common Go functions and patterns.

This library provides reusable, non-business-specific utility functions across
multiple domains, organized into several packages:

# HTTP Package

Package http provides utilities for handling HTTP responses, API error handling, and
middleware patterns.

	import "github.com/StairSupplies/go-core/http"

# REST Package

Package rest provides a fluent REST client for making HTTP requests with built-in
JSON serialization/deserialization.

	import "github.com/StairSupplies/go-core/rest"

# JSON Utils Package

Package jsonutils provides enhanced JSON utilities for encoding and decoding with
better error handling.

	import "github.com/StairSupplies/go-core/jsonutils"

# Log Package

Package log provides structured logging using Uber's Zap with context-aware logging
and both structured and formatted logging options.

	import "github.com/StairSupplies/go-core/log"

# Config Package

Package config provides utilities for loading application configuration from environment
variables with type-safe retrieval using generics.

	import "github.com/StairSupplies/go-core/config"

# Validate Package

Package validate provides input validation utilities with a fluent validation API.

	import "github.com/StairSupplies/go-core/validate"

# String Package

Package str provides string formatting and manipulation utilities.

	import "github.com/StairSupplies/go-core/str"

# TimeUtil Package

Package timeutil provides date and time utility functions, format constants, and
calculations.

	import "github.com/StairSupplies/go-core/timeutil"

See the individual package documentation for more details and examples.
*/
package core
