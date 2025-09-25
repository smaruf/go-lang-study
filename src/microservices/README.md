# Microservices Architecture Example

A complete microservices architecture demonstration using Go, showcasing service communication patterns, API gateway, and modern distributed system practices.

## Architecture Overview

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────────┐
│   API Gateway   │───▶│  User Service   │    │   Order Service     │
│    Port: 8080   │    │   Port: 8081    │◀───│    Port: 8082      │
│                 │    │                 │    │                     │
│ - Load Balancing│    │ - User CRUD     │    │ - Order Management  │
│ - Authentication│    │ - Validation    │    │ - User Integration  │
│ - Rate Limiting │    │ - JWT tokens    │    │ - Event Publishing  │
│ - Circuit Breaker│   │ - Health Check  │    │ - Health Check      │
└─────────────────┘    └─────────────────┘    └─────────────────────┘
```

## Services

### 1. API Gateway (`api-gateway/`)
- **Purpose**: Single entry point for all client requests
- **Features**: 
  - Request routing to appropriate services
  - Load balancing and circuit breaker patterns
  - Authentication and authorization
  - Rate limiting and request throttling
  - Request/response logging and monitoring
  - CORS handling

### 2. User Service (`user-service/`)
- **Purpose**: User management and authentication
- **Features**:
  - User registration and profile management
  - JWT token generation and validation
  - Password hashing and security
  - User data persistence (in-memory for demo)
  - Health check endpoints

### 3. Order Service (`order-service/`)
- **Purpose**: Order processing and management
- **Features**:
  - Order creation and management
  - User validation via User Service communication
  - Order status tracking
  - Service-to-service communication
  - Error handling and retry mechanisms

## Quick Start

### Option 1: Docker Compose (Recommended)

```bash
# Build and start all services
docker-compose up --build

# Test the services
curl http://localhost:8080/health
curl http://localhost:8080/api/users
curl http://localhost:8080/api/orders
```

### Option 2: Manual Start

```bash
# Terminal 1 - User Service
cd user-service
go run main.go

# Terminal 2 - Order Service  
cd order-service
go run main.go

# Terminal 3 - API Gateway
cd api-gateway
go run main.go
```

## API Endpoints

### Through API Gateway (Port 8080)

#### User Management
- `POST /api/users/register` - Register new user
- `POST /api/users/login` - User login
- `GET /api/users` - Get all users (requires auth)
- `GET /api/users/:id` - Get user by ID (requires auth)
- `PUT /api/users/:id` - Update user (requires auth)

#### Order Management
- `GET /api/orders` - Get user orders (requires auth)
- `POST /api/orders` - Create new order (requires auth)
- `GET /api/orders/:id` - Get order by ID (requires auth)
- `PUT /api/orders/:id/status` - Update order status

#### System
- `GET /health` - Overall system health
- `GET /metrics` - System metrics (basic)

### Direct Service Access (Development)

#### User Service (Port 8081)
- `GET /health` - Service health
- `POST /register` - Register user
- `POST /login` - User login
- `GET /users` - List users
- `GET /users/:id` - Get user

#### Order Service (Port 8082)  
- `GET /health` - Service health
- `GET /orders` - List orders
- `POST /orders` - Create order
- `GET /orders/:id` - Get order

## Example Usage

### 1. Register a User
```bash
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

### 2. Login
```bash
curl -X POST http://localhost:8080/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### 3. Create Order (with JWT token)
```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "product": "Go Programming Book",
    "quantity": 2,
    "price": 29.99
  }'
```

## Configuration

Each service supports environment-based configuration:

### Common Variables
- `PORT` - Service port (default varies by service)
- `LOG_LEVEL` - Logging level (debug, info, warn, error)
- `SERVICE_NAME` - Service identifier

### API Gateway Specific
- `USER_SERVICE_URL` - User service endpoint
- `ORDER_SERVICE_URL` - Order service endpoint
- `JWT_SECRET` - JWT signing secret
- `RATE_LIMIT` - Requests per minute limit

### Order Service Specific
- `USER_SERVICE_URL` - User service for validation

## Key Patterns Demonstrated

### 1. Service Discovery
- Environment-based service URLs
- Health check endpoints
- Service registration patterns

### 2. Inter-Service Communication
- HTTP-based communication
- Error handling and retries
- Circuit breaker patterns
- Timeout management

### 3. Authentication & Authorization
- JWT token-based authentication
- Token validation middleware
- User context propagation

### 4. Error Handling
- Graceful degradation
- Error propagation
- Consistent error responses
- Logging and monitoring

### 5. Observability
- Structured logging
- Health check endpoints
- Request tracing
- Metrics collection

## Architecture Patterns

### API Gateway Pattern
- Single entry point for clients
- Request routing and load balancing
- Cross-cutting concerns (auth, logging, rate limiting)

### Service Mesh Concepts
- Service-to-service communication
- Load balancing and discovery
- Security and observability

### Circuit Breaker Pattern
- Prevents cascade failures
- Fallback mechanisms
- Automatic recovery

### Event-Driven Architecture
- Asynchronous communication
- Event publishing and consumption
- Loose coupling between services

## Technology Stack

- **Go**: Primary programming language
- **Gin**: HTTP web framework
- **JWT**: Authentication tokens  
- **Docker**: Containerization
- **Docker Compose**: Multi-service orchestration
- **Logrus**: Structured logging

## Development

### Project Structure
```
microservices/
├── api-gateway/
│   ├── main.go
│   ├── go.mod
│   ├── Dockerfile
│   └── README.md
├── user-service/
│   ├── main.go
│   ├── go.mod
│   ├── Dockerfile
│   └── README.md
├── order-service/
│   ├── main.go
│   ├── go.mod
│   ├── Dockerfile
│   └── README.md
├── docker-compose.yml
└── README.md
```

### Adding New Services

1. Create service directory
2. Implement service with health endpoints
3. Add to docker-compose.yml
4. Update API gateway routing
5. Add service discovery logic

### Testing

```bash
# Unit tests for each service
cd user-service && go test ./...
cd order-service && go test ./...
cd api-gateway && go test ./...

# Integration testing
docker-compose up -d
./run-integration-tests.sh
docker-compose down
```

## Production Considerations

### Security
- Use HTTPS/TLS for all communication
- Implement proper authentication and authorization
- Secure service-to-service communication
- Input validation and sanitization
- Secrets management

### Scalability
- Horizontal service scaling
- Database per service pattern
- Caching strategies
- Load balancing
- Auto-scaling policies

### Observability
- Distributed tracing (Jaeger, Zipkin)
- Metrics collection (Prometheus)
- Centralized logging (ELK stack)
- Health monitoring and alerts

### Reliability
- Circuit breaker implementation
- Retry mechanisms with backoff
- Graceful degradation
- Database replication
- Disaster recovery planning

### Infrastructure
- Kubernetes deployment
- Service mesh (Istio, Linkerd)
- API management platforms
- Container orchestration
- Infrastructure as Code

## Learning Objectives

This example teaches:

1. **Microservices Architecture**: Service decomposition and boundaries
2. **API Gateway Pattern**: Centralized request handling
3. **Service Communication**: HTTP-based inter-service calls
4. **Authentication**: JWT-based security
5. **Error Handling**: Distributed system error management
6. **Containerization**: Docker and Docker Compose
7. **Configuration Management**: Environment-based config
8. **Health Monitoring**: Service health checks
9. **Load Balancing**: Request distribution
10. **Circuit Breaker**: Failure isolation

## Next Steps

- Add database persistence (PostgreSQL, MongoDB)
- Implement message queue (RabbitMQ, Apache Kafka)
- Add distributed tracing
- Implement caching layer (Redis)
- Add monitoring and metrics
- Deploy to Kubernetes
- Implement CI/CD pipeline
- Add comprehensive testing suite