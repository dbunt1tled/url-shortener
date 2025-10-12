# URL Shortener

[![Go Version](https://img.shields.io/github/go-mod/go-version/dbunt1tled/url-shortener)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/dbunt1tled/url-shortener)](https://goreportcard.com/report/github.com/dbunt1tled/url-shortener)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/dbunt1tled/url-shortener)](https://github.com/dbunt1tled/url-shortener/releases)
[![Build Status](https://github.com/dbunt1tled/url-shortener/actions/workflows/release.yml/badge.svg)](https://github.com/dbunt1tled/url-shortener/actions/workflows/release.yml)

A high-performance URL shortening service built with Go, using CloudWeGo's Hertz framework and MySQL for storage.

## Features

- Shorten long URLs to compact, shareable links
- Custom alias support for URLs
- API-first design with RESTful endpoints
- JWT-based authentication
- Database migrations support
- Configurable through environment variables
- Internationalization (i18n) support
- Rate limiting and request validation
- Structured logging
- Graceful shutdown

## Tech Stack

- **Framework**: [CloudWeGo Hertz](https://github.com/cloudwego/hertz)
- **Database**: MySQL
- **ORM**: SQLx
- **Configuration**: Viper
- **Logging**: Custom structured logger
- **Internationalization**: go-i18n
- **Build**: Go 1.25.2+
- **Code Quality**: golangci-lint
- **Release**: GoReleaser

## Getting Started

### Prerequisites

- Go 1.25.2 or higher
- MySQL 5.7 or higher
- Make (optional, for using Makefile commands)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/dbunt1tled/url-shortener.git
   cd url-shortener
   ```

2. Copy the example environment file and update with your configuration:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. Install dependencies:
   ```bash
   go mod download
   ```

4. Run database migrations:
   ```bash
   make migrate-up
   ```

### Running the Application

```bash
# Development mode with hot reload
make dev

# Production build
make build
./bin/url-shortener
```

## Configuration

Configuration is done through environment variables. See `.env.example` for all available options.

### Required Environment Variables

- `DATABASE_DSN`: MySQL connection string (e.g., `user:password@tcp(host:port)/database`)
- `BASE_URL`: Base URL for the shortened links (e.g., `https://short.ly`)
- `JWT_PUBLIC_KEY`: Base64 encoded JWT public key
- `JWT_PRIVATE_KEY`: Base64 encoded JWT private key

## API Documentation

### Shorten URL

```http
POST {{host}}/url
Accept: application/json
Content-Type: application/json
Accept: application/json
Authorization: Bearer {{auth_token}}

{
  "url": "https://www.google.com"
}
```

### Redirect to Original URL

```http
GET /{short_code}
```

## Development

### Building

```bash
make build
```

### Testing

```bash
make test
```

### Linting

```bash
make lint
```

### Running Migrations

```bash
# Create new migration
make migration-create name=migration_name

# Apply migrations
make migrate-up

# Rollback last migration
make migrate-down
```

## Deployment

### Using Docker

```bash
docker-compose up -d
```

### Production

1. Build the production binary:
   ```bash
   make build-release
   ```

2. Deploy the binary and required files to your server.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
