# Go Core

A collection of core Go packages for building robust and maintainable applications.

## Packages

- **config**: Type-safe configuration management with environment variable support
- **httputils**: HTTP response helpers and error handling
- **jsonutils**: JSON serialization and deserialization utilities
- **logger**: Structured logging based on zap
- **rest**: REST client for API interactions
- **str**: String formatting and manipulation utilities
- **timeutils**: Time and date utilities
- **validate**: Data validation helpers

## Installation

```bash
go get github.com/StairSupplies/go-core
```

## Usage

Import the packages you need:

```go
import (
    "github.com/StairSupplies/go-core/config"
    "github.com/StairSupplies/go-core/logger"
    // ... other packages as needed
)
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
