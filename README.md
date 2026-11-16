# Semantic Scholar Go

[![Go Reference](https://pkg.go.dev/badge/github.com/deblinux123/semantic-scholar-go.svg)](https://pkg.go.dev/github.com/deblinux123/semantic-scholar-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A lightweight, reusable, and idiomatic Go SDK for the [Semantic Scholar Academic Graph API](https://api.semanticscholar.org/api-docs/).

## Overview

Semantic Scholar Go provides a simple Go interface for searching academic papers, discovering authors, and retrieving research metadata from Semantic Scholar.

The SDK is designed to be reusable in CLI applications, backend services, research tools, and AI-powered RAG systems.

## Features

- [x] Paper search
- [x] Paper details
- [x] Author search
- [x] Author details
- [x] Context-aware HTTP requests
- [x] Configurable HTTP client
- [x] Custom API base URL
- [x] API error handling
- [ ] Citation and reference retrieval
- [ ] Rate limiting and retries
- [ ] Redis caching

## Installation

```bash
go get github.com/deblinux123/semantic-scholar-go
```

## Package

```go
import "github.com/deblinux123/semantic-scholar-go/semanticscholar"
```

## Quick Start

See the examples directory for runnable examples.

## Configuration

The SDK supports functional options:

- `WithAPIKey` — configure your Semantic Scholar API key.
- `WithTimeout` — configure the HTTP timeout.
- `WithBaseURL` — configure a custom API base URL.
- `WithHTTPClient` — inject a custom `http.Client`.

An API key is optional, but Semantic Scholar may apply stricter rate limits to unauthenticated requests.

## Project Structure

```text
semantic-scholar-go/
├── cmd/
├── docs/
├── examples/
├── semanticscholar/
├── go.mod
├── LICENSE
└── README.md
```

## API Documentation

- [Semantic Scholar API Documentation](https://api.semanticscholar.org/api-docs/)
- [Go Package Documentation](https://pkg.go.dev/github.com/deblinux123/semantic-scholar-go/semanticscholar)

## Contributing

Contributions, bug reports, and feature requests are welcome.

Please open an issue before submitting major changes.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.