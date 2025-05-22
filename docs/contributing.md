# Contributing Guide

Thank you for your interest in contributing to the Go Caching Proxy project! This document provides guidelines and instructions for contributing to the project.

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Git

### Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/your-username/go-caching-proxy.git
   cd go-caching-proxy
   ```
3. Add the original repository as a remote:
   ```bash
   git remote add upstream https://github.com/user/go-caching-proxy.git
   ```
4. Install dependencies:
   ```bash
   go mod download
   ```

## Development Workflow

1. Create a new branch for your feature or bugfix:
   ```bash
   git checkout -b feature/your-feature-name
   ```
   or
   ```bash
   git checkout -b fix/your-bugfix-name
   ```

2. Make your changes

3. Run tests:
   ```bash
   go test ./...
   ```

4. Run the linter:
   ```bash
   # If using golangci-lint
   golangci-lint run
   ```

5. Commit your changes with a descriptive commit message:
   ```bash
   git commit -m "Add feature: your feature description"
   ```

6. Push your branch to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

7. Create a pull request on GitHub

## Code Style Guidelines

- Follow the standard Go code style and conventions as described in [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` or `goimports` to format your code before committing
- Write meaningful comments and documentation
- Include unit tests for new features

## Pull Request Process

1. Update the documentation with details of changes to the interface, if applicable
2. Update the README.md with details of changes to the command-line interface, if applicable
3. The PR must pass all CI tests before it will be reviewed
4. A maintainer will review your PR and may request changes
5. Once approved, a maintainer will merge your PR

## Adding New Features

When adding new features, please follow these guidelines:

1. Discuss major features in an issue before implementing
2. Keep the code modular and maintainable
3. Ensure backwards compatibility when possible
4. Update relevant documentation
5. Add appropriate tests

## Reporting Bugs

When reporting bugs, please include:

- A clear and descriptive title
- Steps to reproduce the issue
- Expected behavior
- Actual behavior
- Go version and operating system

## License

By contributing to Go Caching Proxy, you agree that your contributions will be licensed under the project's MIT license.
