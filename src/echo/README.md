# Echo Framework Example

A comprehensive Echo web server demonstrating modern Go web development patterns, middleware usage, and API design.

## Features

- **Echo Framework**: High-performance, minimalist web framework
- **Middleware Stack**: Logger, Recovery, CORS, Security, Rate Limiting
- **Structured Logging**: JSON logging with Logrus and request tracking
- **RESTful API**: Complete CRUD operations with proper error handling
- **Environment Configuration**: Environment-based settings
- **Graceful Shutdown**: Proper server lifecycle management
- **Security Headers**: Built-in security middleware
- **Rate Limiting**: Memory-based rate limiting
- **Request ID Tracking**: Request tracing across the application

## Quick Start

1. **Install dependencies:**
   ```bash
   go mod download
   ```

2. **Setup environment:**
   ```bash
   cp .env.example .env
   # Edit .env if needed
   ```

3. **Run the server:**
   ```bash
   go run main.go
   ```

4. **Visit the application:**
   - Web interface: http://localhost:1323
   - Health check: http://localhost:1323/health
   - API: http://localhost:1323/api/v1/users

## API Endpoints

### Core Routes
- `GET /` - Interactive web interface with documentation
- `GET /health` - Service health check

### User Management API
- `GET /api/v1/users` - Get all users
- `POST /api/v1/users` - Create a new user
- `GET /api/v1/users/:id` - Get user by ID
- `PUT /api/v1/users/:id` - Update user by ID
- `DELETE /api/v1/users/:id` - Delete user by ID

### Echo Test API
- `GET /api/v1/echo/:message` - Echo URL parameter
- `POST /api/v1/echo` - Echo POST JSON data

## Configuration

Environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `1323` | Server port |
| `LOG_LEVEL` | `info` | Logging level (debug, info, warn, error) |
| `ENV` | `development` | Environment (development, production) |

## Testing the API

### Using curl

```bash
# Health check
curl http://localhost:1323/health

# Get all users
curl http://localhost:1323/api/v1/users

# Create user
curl -X POST http://localhost:1323/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Johnson","email":"alice@example.com"}'

# Echo message
curl http://localhost:1323/api/v1/echo/HelloWorld

# Echo POST data
curl -X POST http://localhost:1323/api/v1/echo \
  -H "Content-Type: application/json" \
  -d '{"message":"Hello","timestamp":"2024-01-01"}'
```

### Response Format

All API responses follow a consistent format:

```json
{
  "success": true,
  "data": {...},
  "message": "Optional message",
  "error": "Error message if success is false"
}
```

## Architecture

### Middleware Stack
1. **Logger**: Request/response logging
2. **Recover**: Panic recovery
3. **CORS**: Cross-origin resource sharing
4. **Secure**: Security headers
5. **RequestID**: Request tracking
6. **RateLimiter**: Rate limiting (10 req/sec)
7. **StructuredLogging**: Custom structured logging

### Project Structure
```
echo/
├── main.go              # Main application
├── go.mod              # Dependencies
├── go.sum              # Dependency checksums
├── .env.example        # Environment template
├── README.md           # This file
└── simple_echo.go      # Original simple example
```

### Key Components

1. **Server Struct**: Wraps Echo instance with dependencies
2. **Structured Logging**: Request tracking with Logrus
3. **Error Handling**: Consistent error responses
4. **Environment Config**: Environment-based configuration
5. **Graceful Shutdown**: Proper server lifecycle

## Echo Framework Features Demonstrated

### Middleware Usage
```go
s.echo.Use(middleware.Logger())
s.echo.Use(middleware.Recover())
s.echo.Use(middleware.CORS())
s.echo.Use(middleware.Secure())
s.echo.Use(middleware.RequestID())
s.echo.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(10)))
```

### Route Grouping
```go
api := s.echo.Group("/api/v1")
users := api.Group("/users")
users.GET("", s.getUsersHandler)
users.POST("", s.createUserHandler)
```

### Custom Middleware
```go
func (s *Server) structuredLoggingMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            // Custom logging logic
        }
    }
}
```

### Context Usage
```go
func (s *Server) getUserHandler(c echo.Context) error {
    id, err := strconv.Atoi(c.Param("id"))
    // Handle request
    return c.JSON(http.StatusOK, response)
}
```

## Key Concepts Demonstrated

1. **Echo Framework**: High-performance web framework usage
2. **Middleware**: Building and using middleware stack
3. **Structured Logging**: JSON logging with contextual information
4. **RESTful Design**: Proper REST API patterns
5. **Error Handling**: Consistent error response patterns
6. **Environment Config**: Configuration management
7. **Graceful Shutdown**: Server lifecycle management
8. **Security**: Security headers and rate limiting
9. **Request Tracking**: Request ID generation and tracking
10. **JSON Processing**: Request/response JSON handling

## Comparison with Standard Library

Echo provides several advantages over the standard `net/http`:
- **Performance**: Optimized router and middleware
- **Middleware**: Rich built-in middleware ecosystem
- **Binding**: Automatic request binding to structs
- **Validation**: Built-in validation support
- **Error Handling**: Centralized error handling
- **Route Groups**: Organized route management

## Production Considerations

For production deployment:
- Add authentication/authorization middleware
- Implement proper database integration
- Add metrics and monitoring
- Configure TLS/HTTPS
- Set up load balancing
- Add caching layer
- Implement circuit breaker patterns
- Add input validation and sanitization

## Learning Objectives

This example teaches:
- Echo framework fundamentals
- Middleware development and usage
- Structured logging implementation
- RESTful API design patterns
- Environment-based configuration
- Graceful server shutdown
- Security best practices
- Request/response handling
- Error management patterns

## Related Examples

- [Web Server](../web-server/) - Comparison with Gorilla Mux
- [CLI Tool](../cli-tool/) - Command-line applications
- [Concurrency](../concurrency/) - Concurrent programming patterns
