# Modern Go Web Server

A comprehensive example of a modern Go web server implementing REST API patterns, middleware, structured logging, and graceful shutdown.

## Features

- **RESTful API**: Complete CRUD operations for user management
- **Middleware Chain**: Logging, CORS, and content-type middleware
- **Structured Logging**: JSON-formatted logs with contextual information
- **Graceful Shutdown**: Proper server shutdown handling
- **Environment Configuration**: Environment-based configuration management
- **Error Handling**: Consistent error response format
- **Request Validation**: Input validation and sanitization
- **Static File Serving**: Serving static assets
- **Health Checks**: Service health monitoring endpoint

## Project Structure

```
web-server/
├── main.go              # Main application with server implementation
├── go.mod              # Go module dependencies
├── .env.example        # Environment configuration example
├── README.md           # This file
└── static/             # Static files directory (optional)
```

## Quick Start

1. **Setup environment:**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Run the server:**
   ```bash
   go run main.go
   ```

4. **Visit the application:**
   - Web interface: http://localhost:8080
   - API health check: http://localhost:8080/api/v1/health
   - Get users: http://localhost:8080/api/v1/users

## API Endpoints

### Health Check
```
GET /api/v1/health
```
Returns server health status and version information.

### User Management

#### Get All Users
```
GET /api/v1/users
```

#### Get User by ID
```
GET /api/v1/users/{id}
```

#### Create User  
```
POST /api/v1/users
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com"
}
```

#### Update User
```
PUT /api/v1/users/{id}
Content-Type: application/json

{
  "name": "John Smith",
  "email": "johnsmith@example.com"
}
```

#### Delete User
```
DELETE /api/v1/users/{id}
```

## Configuration

Set these environment variables in your `.env` file:

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `HOST` | `localhost` | Server host |
| `LOG_LEVEL` | `info` | Logging level (debug, info, warn, error) |
| `ENV` | `development` | Environment (development, production) |

## Testing the API

### Using curl

```bash
# Health check
curl http://localhost:8080/api/v1/health

# Get all users
curl http://localhost:8080/api/v1/users

# Get specific user
curl http://localhost:8080/api/v1/users/1

# Create a new user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Bob Wilson","email":"bob@example.com"}'

# Update user
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"John Updated","email":"john.updated@example.com"}'

# Delete user
curl -X DELETE http://localhost:8080/api/v1/users/1
```

### Using HTTPie

```bash
# Health check
http GET :8080/api/v1/health

# Get all users
http GET :8080/api/v1/users

# Create user
http POST :8080/api/v1/users name="Alice Johnson" email="alice@example.com"

# Update user
http PUT :8080/api/v1/users/1 name="John Updated" email="john.updated@example.com"
```

## Architecture Patterns

### Middleware Chain
The server implements a middleware chain with:
- **Logging Middleware**: Logs all HTTP requests with timing and metadata
- **CORS Middleware**: Handles Cross-Origin Resource Sharing
- **Content-Type Middleware**: Sets default response content type

### Structured Logging
Uses `logrus` for structured JSON logging with contextual fields:
```go
s.logger.WithFields(logrus.Fields{
    "user_id": user.ID,
    "name":    user.Name,
    "email":   user.Email,
}).Info("User created")
```

### Graceful Shutdown
Implements proper graceful shutdown:
- Listens for OS signals (SIGINT, SIGTERM)
- Completes in-flight requests
- Times out after 30 seconds
- Proper resource cleanup

### Error Handling
Consistent error response format:
```json
{
  "success": false,
  "error": "User not found",
  "data": null,
  "message": null
}
```

### Configuration Management
Environment-based configuration with sensible defaults:
```go
config := &Config{
    Port:         getEnv("PORT", "8080"),
    Host:         getEnv("HOST", "localhost"),
    LogLevel:     getEnv("LOG_LEVEL", "info"),
}
```

## Key Concepts Demonstrated

1. **HTTP Server Patterns**: Router setup, handler functions, middleware
2. **Context Usage**: Request context, graceful shutdown context
3. **Structured Logging**: JSON logging with contextual information
4. **Environment Configuration**: Environment variables with defaults
5. **Error Handling**: Consistent error responses and logging
6. **JSON Processing**: Request/response JSON marshaling
7. **Validation**: Input validation and sanitization
8. **Resource Management**: Proper cleanup and graceful shutdown
9. **Observability**: Request logging, metrics, health checks
10. **Security**: CORS headers, input validation

## Production Considerations

For production deployment, consider adding:
- Authentication and authorization middleware
- Rate limiting
- Request ID tracking
- Metrics collection (Prometheus)
- Database integration
- Caching layer
- Load balancing
- TLS/HTTPS configuration
- Container deployment (Docker)
- Monitoring and alerting

## Learning Objectives

After studying this example, you'll understand:
- How to build production-ready HTTP servers in Go
- Middleware pattern implementation
- Structured logging and observability
- Environment-based configuration
- Graceful shutdown patterns
- RESTful API design
- Error handling best practices
- JSON API development