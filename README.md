# Go Core

A collection of core Go packages for building robust and maintainable applications.

## Packages

- **api**: HTTP API response helpers and error handling
- **config**: Type-safe configuration management with environment variable support
- **jsonutils**: JSON serialization and deserialization utilities
- **logger**: Structured logging based on zap
- **rest**: REST client for API interactions
- **router**: Opinionated chi-based HTTP router with middleware

## Installation

```bash
go get github.com/StairSupplies/go-core
```

## Usage

### API Package

The `api` package provides a standardized way to handle HTTP responses and errors in RESTful APIs:

```go
import "github.com/StairSupplies/go-core/api"

// In your HTTP handler
func getUserHandler(w http.ResponseWriter, r *http.Request) error {
    user, err := userService.GetUser(id)
    if err != nil {
        return api.NotFoundError(fmt.Errorf("user not found: %w", err))
    }
    
    return api.WriteSuccess(w, user)
}

// Wrap the handler for automatic error handling
http.HandleFunc("/api/users/{id}", api.WrapHandler(getUserHandler))
```

### Router Package

The `router` package provides an opinionated HTTP router based on chi with built-in middleware:

```go
import (
    "github.com/StairSupplies/go-core/api"
    "github.com/StairSupplies/go-core/router"
)

// Create a new router with default middleware
r := router.New()

// Add routes
r.Get("/health", healthCheckHandler)
r.Get("/api/users/{id}", router.WithErrorHandler(getUserHandler))

// Create a protected API group
protectedAPI := r.Group(func(r chi.Router) {
    r.Use(authMiddleware)
    r.Get("/profile", router.WithErrorHandler(getProfileHandler))
})

// Mount the protected API
r.Mount("/api/auth", protectedAPI)

// Start the server
http.ListenAndServe(":8080", r)
```

### Config Package

The `config` package provides type-safe configuration management with environment variable support:

```go
import "github.com/StairSupplies/go-core/config"

// Define your configuration structure
type AppConfig struct {
    Server struct {
        Port    int    `mapstructure:"port"`
        Host    string `mapstructure:"host"`
        Timeout int    `mapstructure:"timeout"`
    } `mapstructure:"server"`
    Database struct {
        DSN      string `mapstructure:"dsn"`
        MaxConns int    `mapstructure:"max_conns"`
    } `mapstructure:"database"`
}

// Load configuration
var cfg AppConfig
err := config.Load("config", &cfg)
```

### Logger Package

The `logger` package provides structured logging based on zap:

```go
import "github.com/StairSupplies/go-core/logger"

// Initialize the logger
err := logger.Init(logger.Config{
    Level:       "info",
    Development: true,
    ServiceName: "my-service",
})

// Log with structured fields
logger.Info("Server started", 
    zap.Int("port", 8080),
    zap.String("environment", "development"),
)
```

### JSON Utilities

The `jsonutils` package provides utilities for JSON handling:

```go
import "github.com/StairSupplies/go-core/jsonutils"

// Parse JSON with error handling
data, err := jsonutils.Parse(jsonString)

// Format JSON with indentation
formattedJSON, err := jsonutils.Format(jsonObj)
```

### REST Client

The `rest` package provides a client for API interactions:

```go
import "github.com/StairSupplies/go-core/rest"

// Create a new client
client := rest.NewClient().
    WithBaseURL("https://api.example.com").
    WithHeader("Authorization", "Bearer token").
    WithTimeout(30 * time.Second)

// Make a request
var response MyResponse
err := client.Get(ctx, "/users/123", &response)
```

## Versioning

This project follows [Semantic Versioning](https://semver.org/). Releases are automatically created when changes are merged to the main branch.

- Major version increase (X.y.z) - incompatible API changes
- Minor version increase (x.Y.z) - backwards-compatible functionality
- Patch version increase (x.y.Z) - backwards-compatible bug fixes

We use conventional commits to determine version bumps automatically.

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.