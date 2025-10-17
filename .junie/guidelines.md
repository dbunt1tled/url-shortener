# URL Shortener Development Guidelines

This document contains essential development information for the URL Shortener project.

## Project Overview

This is a Go-based URL shortener service built with the Hertz web framework, MySQL database, and dependency injection architecture. The application uses RSA/ECDSA JWT tokens for authentication and provides a REST API for URL shortening operations.

## Build and Configuration

### Prerequisites

- Go 1.25.2 or later
- MySQL database
- Goose migration tool for database migrations

### Environment Configuration

The application requires a `.env` file in the project root. Use `.env.example` as a template:

**Required Environment Variables:**
- `APP_NAME` - Application name
- `BASE_URL` - Base URL for shortened links (required for link generation)
- `DATABASE_DSN` - MySQL connection string: `user:password@tcp(host:port)/database?charset=utf8mb4&parseTime=True`
- `JWT_PRIVATE_KEY` - Base64 encoded RSA/ECDSA private key for JWT signing
- `JWT_PUBLIC_KEY` - Base64 encoded RSA/ECDSA public key for JWT verification
- `JWT_TOKEN_ALGORITHM` - Supported: `RS256`, `RS512`, `ES256`, `ES512`

**Optional Environment Variables:**
- `HTTP_SERVER_ADDRESS` - Server address (default: `localhost:8083`)
- `LOG_LEVEL` - Logging level: -4=Debug, -3=Info, -2=Warn, -1=Error, 0=Fatal
- CORS settings for cross-origin requests

### Database Setup

1. Create MySQL database
2. Configure `DATABASE_DSN` in `.env` file
3. Run migrations using Make commands:
   ```bash
   # Check migration status
   make migrate_status
   
   # Run migrations
   make migrate_up
   
   # Rollback migrations (if needed)
   make migrate_down
   ```

### Building and Running

```bash
# Download dependencies
go mod download

# Build the application
go build -o bin/urlshortener cmd/urlshortener/main.go

# Or run directly using Make
make run api
```

## Architecture Notes

### Key Components

- **Hertz Framework**: CloudWeGo's high-performance HTTP framework
- **Dependency Injection**: Uses uber/dig for clean dependency management
- **Database**: MySQL with sqlx for enhanced SQL operations
- **JWT Authentication**: RSA/ECDSA signed tokens (no HMAC support)
- **Internationalization**: Support for multiple languages (ru/en)
- **Graceful Shutdown**: Proper resource cleanup on application termination

### Important Implementation Details

- The hasher package uses base62 encoding (not hex) despite method naming
- JWT implementation only supports RSA and ECDSA algorithms, not HMAC
- Database migrations are managed via goose with MySQL driver
- The application uses middleware for authentication, CORS, logging, and error handling

## Testing

### Running Tests

```bash
# Run tests for specific package
go test ./internal/lib/hasher -v

# Run all tests
go test ./... -v

# Run tests with coverage
go test -cover ./...
```

### Test Example

A working test example is available in `internal/lib/hasher/hasher_test.go`:

```go
// TestRandomBytes - Tests random byte generation
func TestRandomBytes(t *testing.T) {
    hasher := &Hasher{}
    result, err := hasher.RandomBytes(10)
    if err != nil {
        t.Errorf("RandomBytes() error = %v", err)
    }
    if len(result) != 10 {
        t.Errorf("Expected length 10, got %d", len(result))
    }
}
```

### Adding New Tests

1. Create `*_test.go` files in the same package as the code being tested
2. Use table-driven tests for multiple test cases
3. Test both success and error conditions
4. For components requiring JWT setup, use appropriate RSA/ECDSA keys
5. Keep tests simple and focused on individual functionality

### Test Structure Best Practices

- Use descriptive test names: `TestFunctionName_Condition_ExpectedResult`
- Use subtests with `t.Run()` for related test cases
- Always check for both error conditions and expected results
- Use `t.Fatal()` for setup failures, `t.Error()` for test failures

## Development Best Practices

### Code Organization

- Follow Go project layout standards
- Use dependency injection for clean architecture
- Separate concerns: handlers, services, repositories
- Keep configuration centralized in the config package

### Database Operations

- Use sqlx for enhanced SQL operations with struct mapping
- Always use prepared statements for SQL queries
- Handle database connections gracefully with proper cleanup
- Use migrations for all schema changes

### Error Handling

- Use structured error handling with pkg/errors
- Provide meaningful error messages for debugging
- Log errors appropriately based on severity level
- Handle database connection errors gracefully

### Security Considerations

- JWT tokens require proper RSA/ECDSA key management
- Use base64 encoding for storing keys in environment variables
- Validate all input data through the custom validator
- Implement proper CORS policies for production environments

### Logging

- Use structured logging with appropriate log levels
- Debug level (-4) for development, higher levels for production
- Log important operations like database connections and shutdowns
- Configure logger level status to control HTTP request logging

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

---

*Last updated: 2025-10-17*