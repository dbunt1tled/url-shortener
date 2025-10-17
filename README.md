# URL Shortener

[![Go Version](https://img.shields.io/github/go-mod/go-version/dbunt1tled/url-shortener)](https://golang.org/)
[![Go Reference](https://pkg.go.dev/badge/github.com/dbunt1tled/url-shortener.svg)](https://pkg.go.dev/github.com/dbunt1tled/url-shortener)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/dbunt1tled/url-shortener)](https://goreportcard.com/report/github.com/dbunt1tled/url-shortener)

A high-performance URL shortening service built with Go, using CloudWeGo's Hertz framework, MySQL database, and dependency injection architecture.

## Features

- **URL Shortening**: Convert long URLs into compact, shareable links using base62 encoding
- **JWT Authentication**: RSA/ECDSA signed tokens for secure API access (no HMAC support)
- **Database Migrations**: Automated database schema management with Goose
- **Internationalization**: Multi-language support (Russian/English) with go-i18n
- **Structured Logging**: Configurable log levels with pretty formatting
- **Graceful Shutdown**: Proper resource cleanup on application termination
- **CORS Support**: Configurable cross-origin request handling
- **Input Validation**: Custom validation middleware for request data
- **Dependency Injection**: Clean architecture using uber/dig

## Tech Stack

- **Framework**: [CloudWeGo Hertz](https://github.com/cloudwego/hertz) - High-performance HTTP framework
- **Database**: MySQL with [SQLx](https://github.com/jmoiron/sqlx) for enhanced SQL operations
- **Configuration**: [Viper](https://github.com/spf13/viper) with environment variable support
- **JWT**: [golang-jwt/jwt](https://github.com/golang-jwt/jwt) with RSA/ECDSA algorithms
- **Migrations**: [Goose](https://github.com/pressly/goose) database migration tool
- **Internationalization**: [go-i18n](https://github.com/nicksnyder/go-i18n)
- **ID Generation**: [base62](https://github.com/jxskiss/base62) encoding for short URLs
- **Build**: Go 1.25.3+

## Getting Started

### Prerequisites

- Go 1.25.2 or later
- MySQL database
- Goose migration tool for database migrations

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/dbunt1tled/url-shortener.git
   cd url-shortener
   ```

2. Copy the example environment file and configure it:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration (see Configuration section below)
   ```

3. Install dependencies:
   ```bash
   go mod download
   ```

4. Set up database and run migrations:
   ```bash
   # Check migration status
   make migrate_status
   
   # Run migrations
   make migrate_up
   ```

### Running the Application

```bash
# Run directly using Make
make run api

# Or build and run manually
go build -o bin/urlshortener cmd/urlshortener/main.go
./bin/urlshortener
```

## Configuration

The application requires a `.env` file in the project root. Use `.env.example` as a template.

### Required Environment Variables

- `APP_NAME` - Application name
- `BASE_URL` - Base URL for shortened links (required for link generation)
- `DATABASE_DSN` - MySQL connection string: `user:password@tcp(host:port)/database?charset=utf8mb4&parseTime=True`
- `JWT_PRIVATE_KEY` - Base64 encoded RSA/ECDSA private key for JWT signing
- `JWT_PUBLIC_KEY` - Base64 encoded RSA/ECDSA public key for JWT verification
- `JWT_TOKEN_ALGORITHM` - Supported algorithms: `RS256`, `RS512`, `ES256`, `ES512`

### Optional Environment Variables

- `HTTP_SERVER_ADDRESS` - Server address (default: `localhost:8083`)
- `LOG_LEVEL` - Logging level: `-4`=Debug, `-3`=Info, `-2`=Warn, `-1`=Error, `0`=Fatal
- `LOGGER_LEVEL_STATUS` - HTTP status code threshold for request logging
- CORS settings: `ACCESS_CONTROL_ALLOW_HEADERS`, `ACCESS_CONTROL_ALLOW_METHODS`, etc.

## API Documentation

### Shorten URL

Create a shortened URL:

```http
POST /url
Content-Type: application/json
Authorization: Bearer <your_jwt_token>
Accept-Language: en

{
  "url": "https://www.example.com/very/long/url/path"
}
```

**Response:**
```json
{
  "short_url": "http://localhost:8083/abc123",
  "code": "abc123"
}
```

### Access Shortened URL

Redirect to the original URL:

```http
GET /{short_code}
```

**Response:** HTTP 302 redirect to the original URL

## Development

### Available Make Commands

```bash
# Run the application
make run api

# Database migrations
make migrate_up      # Apply all pending migrations
make migrate_down    # Rollback last migration
make migrate_status  # Check migration status

# Create new migrations
make migration_sql MIGRATION_NAME=your_migration_name
make migration_go MIGRATION_NAME=your_migration_name
```

### Testing

```bash
# Run all tests
go test ./... -v

# Run tests for specific package
go test ./internal/lib/hasher -v

# Run tests with coverage
go test -cover ./...
```

### Project Structure

```
├── cmd/urlshortener/           # Application entry point
├── internal/
│   ├── app/                    # Application layer
│   │   ├── container/          # Dependency injection setup
│   │   └── shorturl/           # URL shortening domain
│   ├── config/                 # Configuration and logging
│   ├── domain/                 # Domain entities and enums
│   ├── lib/                    # Shared libraries
│   │   ├── hasher/            # JWT and ID generation
│   │   ├── http-server/       # HTTP middleware and utilities
│   │   └── validator/         # Input validation
└── storage/
    ├── migrations/            # Database migration files
    └── mysql/                 # MySQL connection setup
```

### Important Implementation Notes

- The hasher package uses base62 encoding despite method naming conventions
- JWT implementation only supports RSA and ECDSA algorithms, not HMAC
- Database migrations use Goose with MySQL driver
- The application uses middleware for authentication, CORS, logging, and error handling
- Graceful shutdown is implemented with proper resource cleanup

## Troubleshooting

### Common Issues

1. **JWT Algorithm Errors**: Ensure you're using supported algorithms (RS256, RS512, ES256, ES512)
2. **Database Connection**: Verify DATABASE_DSN format and database accessibility
3. **Migration Issues**: Check GOOSE_DRIVER and GOOSE_DBSTRING environment variables
4. **Build Failures**: Ensure Go version compatibility (1.25.2+)

### Debugging Tips

- Use debug log level (-4) for verbose logging
- Check database migrations status before running the application
- Verify all required environment variables are set
- Test database connectivity independently before starting the application

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

*Last updated: 2025-10-17*
