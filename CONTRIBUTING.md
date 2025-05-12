# Contributing to go-core

Thank you for considering contributing to go-core! This document outlines how to contribute to the project, including our commit message conventions which are critical for our automated release and versioning system.

## Commit Message Guidelines

We follow the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) specification for our commit messages. This enables automatic versioning and changelog generation.

### Commit Message Format

Each commit message consists of a **header**, a **body** and a **footer**. The header has a special format that includes a **type**, a **scope** and a **subject**:

```
<type>(<scope>): <subject>
<BLANK LINE>
<body>
<BLANK LINE>
<footer>
```

The **header** is mandatory and the **scope** of the header is optional.

### Type

The type must be one of the following:

- **feat**: A new feature
- **fix**: A bug fix
- **docs**: Documentation only changes
- **style**: Changes that do not affect the meaning of the code (white-space, formatting, etc)
- **refactor**: A code change that neither fixes a bug nor adds a feature
- **perf**: A code change that improves performance
- **test**: Adding missing tests or correcting existing tests
- **build**: Changes that affect the build system or external dependencies
- **ci**: Changes to our CI configuration files and scripts
- **chore**: Other changes that don't modify src or test files

### Scope

The scope should be the name of the package affected (as perceived by the person reading the changelog generated from commit messages).

### Subject

The subject contains a succinct description of the change:

- use the imperative, present tense: "change" not "changed" nor "changes"
- don't capitalize the first letter
- no dot (.) at the end

### Body

The body should include the motivation for the change and contrast this with previous behavior.

### Footer

The footer should contain any information about **Breaking Changes** and is also the place to reference GitHub issues that this commit **Closes**.

**Breaking Changes** should start with the word `BREAKING CHANGE:` with a space or two newlines. The rest of the commit message is then used for this.

### Examples

```
feat(config): add support for loading nested config

Add the ability to load nested configuration from environment variables
using dot notation.

Closes #123
```

```
fix(logger): ensure log files are flushed before program exit

Fixes an issue where log entries could be lost when the program exits
abruptly.
```

```
BREAKING CHANGE: remove deprecated ValidateConfig function

The ValidateConfig function has been removed in favor of the new
Validator interface. To migrate, replace ValidateConfig calls with
the new Validate method on Config struct.

Closes #456
```

## Pull Request Process

1. Ensure any install or build dependencies are removed before the end of the layer when doing a build.
2. Update the README.md or documentation with details of changes to the interface, this includes new environment variables, exposed ports, useful file locations and container parameters.
3. Increase the version numbers in any examples files and the README.md to the new version that this Pull Request would represent. The versioning scheme we use is [SemVer](http://semver.org/).
4. The Pull Request will be merged once you have the sign-off of at least one maintainer, or if you do not have permission to do that, you may request the reviewer to merge it for you.

## Branch Naming Convention

We follow this branch naming convention to ensure our automatic labeler can correctly categorize your PR:

- `feat/*` or `feature/*` - for features
- `fix/*`, `bugfix/*` or `bug/*` - for bug fixes
- `chore/*` or `maintenance/*` - for maintenance tasks
- `docs/*` or `documentation/*` - for documentation
- `style/*` - for style/formatting changes
- `refactor/*` - for code refactoring
- `test/*` or `testing/*` - for test-related changes
- `build/*` - for build system changes
- `ci/*` - for CI-related changes
- `perf/*` or `performance/*` - for performance improvements