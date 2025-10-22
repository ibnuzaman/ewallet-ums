# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial project setup with Clean Architecture
- HTTP server with Chi router
- Graceful shutdown implementation
- Structured logging with Logrus
- Configuration management with environment variables
- Standardized API response format
- Request ID tracking
- Health check endpoint
- Middleware stack:
  - Request logging
  - Panic recovery
  - Request timeout (60s)
  - Request ID injection
- Dependency injection pattern
- Interface-based architecture
- Unit tests for API and Services layers (100% coverage)
- Docker support with multi-stage builds
- Docker Compose configuration
- GitHub Actions CI/CD pipeline
- Makefile for common tasks
- GolangCI-lint configuration
- API documentation
- Example environment file

### Security
- Non-root user in Docker container
- Health check in Docker container
- Proper error handling without exposing sensitive information

## [0.1.0] - 2025-10-22

### Added
- Initial boilerplate release
