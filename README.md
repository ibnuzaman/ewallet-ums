# E-Wallet UMS (User Management Service)

A microservice for user management in an e-wallet system built with Go.

## 🏗️ Architecture

This project follows Clean Architecture principles with the following layers:

- **cmd/**: Application entry points and server setup
- **internal/**: Business logic (API handlers, services, repositories)
- **helpers/**: Shared utilities (config, logger, response)
- **models/**: Data models
- **constants/**: Application constants

## 🚀 Getting Started

### Prerequisites

- Go 1.23.3 or higher
- Make (optional)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/ibnuzaman/ewallet-ums.git
cd ewallet-ums
```

2. Install dependencies:
```bash
go mod download
```

3. Copy environment file:
```bash
cp .env.example .env
```

4. Run the application:
```bash
go run main.go
```

The server will start on port 8080 (or the port specified in your `.env` file).

## 📝 API Endpoints

### Health Check
- **GET** `/healthcheck` - Check service health status

## 🔧 Configuration

Configuration is managed through environment variables. See `.env.example` for available options.

### Key Configuration Variables:

- `PORT`: Server port (default: 8080)
- `ENVIRONMENT`: Environment mode (development/production)

## 🏃 Running in Development

```bash
go run main.go
```

## 🏗️ Building

```bash
go build -o bin/ewallet-ums main.go
```

## 📦 Project Structure

```
.
├── cmd/
│   ├── http.go              # HTTP server setup with graceful shutdown
│   └── proto/               # gRPC server (future)
├── helpers/
│   ├── config.go            # Configuration management
│   ├── logger.go            # Logging setup & middleware
│   └── response.go          # Standard API responses
├── internal/
│   ├── api/                 # HTTP handlers
│   ├── services/            # Business logic
│   ├── repository/          # Data access layer
│   ├── models/              # Data models
│   ├── constants/           # Constants
│   └── interfaces/          # Interface definitions
├── main.go                  # Application entry point
├── go.mod
└── README.md
```

## ✨ Features

- ✅ Clean Architecture
- ✅ Dependency Injection
- ✅ Structured Logging (Logrus)
- ✅ Graceful Shutdown
- ✅ Request Logging Middleware
- ✅ Panic Recovery Middleware
- ✅ Request Timeout
- ✅ Environment-based Configuration
- ✅ Standardized API Responses
- ✅ Request ID Tracking

## 🛠️ Tech Stack

- **Framework**: Chi Router
- **Logger**: Logrus
- **Config**: godotenv
- **Language**: Go 1.23.3

## 📄 License

MIT

## 👥 Author

ibnuzaman
