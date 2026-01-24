# 3-Weeks Go Learning Plan

A comprehensive 3-week structured learning plan for mastering Go through building a real-world Task Management REST API. This project demonstrates production-ready Go application architecture, best practices, and modern patterns.

## 🎯 Project Overview

This is a complete **Task Management REST API** built with Go that demonstrates:
- Standard Go project layout
- RESTful API design
- Database integration (SQLite)
- Comprehensive testing
- Clean architecture
- Environment-based configuration
- Middleware patterns
- Error handling strategies

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or later
- SQLite3

### Installation & Running

1. **Navigate to the project directory:**
   ```bash
   cd src/3-weeks-plan
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Setup environment:**
   ```bash
   cp .env.example .env
   # Edit .env if needed
   ```

4. **Run the application:**
   ```bash
   go run cmd/api/main.go
   ```

5. **The server will start on http://localhost:8080**

### Quick Test

Create a task:
```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Learn Go","description":"Complete 3-weeks plan","priority":"high"}'
```

List all tasks:
```bash
curl http://localhost:8080/api/tasks
```

## 📁 Project Structure

```
3-weeks-plan/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── pkg/
│   ├── models/
│   │   └── task.go              # Data models
│   ├── database/
│   │   └── database.go          # Database layer
│   └── handlers/
│       └── handlers.go          # HTTP handlers
├── internal/
│   └── config/
│       └── config.go            # Configuration management
├── web/
│   ├── static/                  # Static files (future)
│   └── templates/               # HTML templates (future)
├── db/
│   └── migrations/              # Database migrations (future)
├── tests/
│   └── handlers_test.go         # Integration tests
├── .env.example                 # Environment template
├── .gitignore                   # Git ignore rules
├── go.mod                       # Go module definition
└── README.md                    # This file
```

### Directory Explanation

- **cmd/**: Application entry points (main packages)
- **pkg/**: Reusable library code that can be imported by other projects
- **internal/**: Private application code (cannot be imported by other projects)
- **web/**: Web assets (static files, templates)
- **db/**: Database-related code (migrations, seeds)
- **tests/**: End-to-end and integration tests

## 📚 3-Week Learning Plan

### Week 1: Foundations & Setup
**Goal**: Understand Go basics and project structure

#### Day 1-2: Go Fundamentals
- [ ] Variables, types, and constants
- [ ] Functions and methods
- [ ] Structs and interfaces
- [ ] Pointers and memory management

**Practice**: Study `pkg/models/task.go`
- Understand struct tags for JSON serialization
- Learn about pointer vs value receivers
- Explore time.Time handling

#### Day 3-4: Project Layout
- [ ] Standard Go project layout
- [ ] Package organization
- [ ] Import paths and modules
- [ ] Public vs private (exported vs unexported)

**Practice**: Analyze the project structure
- Why is config in `internal/`?
- What makes `pkg/` reusable?
- How does `cmd/api/main.go` tie everything together?

#### Day 5-7: HTTP & REST APIs
- [ ] net/http package basics
- [ ] HTTP methods (GET, POST, PUT, DELETE)
- [ ] Request/response handling
- [ ] JSON encoding/decoding
- [ ] Router patterns with gorilla/mux

**Practice**: Study `pkg/handlers/handlers.go`
- Implement a new endpoint
- Add request validation
- Understand middleware patterns

### Week 2: Database & Business Logic
**Goal**: Master data persistence and application logic

#### Day 8-10: Database Integration
- [ ] SQL basics and SQLite
- [ ] database/sql package
- [ ] Prepared statements
- [ ] Connection management
- [ ] CRUD operations

**Practice**: Study `pkg/database/database.go`
- Add a new field to tasks
- Implement filtering by date range
- Create an index for optimization

#### Day 11-12: Error Handling
- [ ] Error interface
- [ ] Error wrapping with fmt.Errorf
- [ ] Custom error types
- [ ] Panic and recover
- [ ] Error propagation patterns

**Practice**: Enhance error handling
- Add validation errors
- Create custom error types
- Implement error logging

#### Day 13-14: Testing
- [ ] Unit testing basics
- [ ] Table-driven tests
- [ ] Test fixtures and setup
- [ ] Mocking dependencies
- [ ] Integration testing

**Practice**: Study `tests/handlers_test.go`
- Add tests for edge cases
- Test error scenarios
- Measure test coverage

### Week 3: Advanced Patterns & Production
**Goal**: Production-ready code and best practices

#### Day 15-16: Configuration & Environment
- [ ] Environment variables
- [ ] Configuration management
- [ ] Secrets handling
- [ ] Multi-environment setup

**Practice**: Study `internal/config/config.go`
- Add new configuration options
- Implement environment-specific settings
- Add validation for required config

#### Day 17-18: Middleware & Logging
- [ ] Middleware pattern
- [ ] Request logging
- [ ] Authentication/authorization
- [ ] Rate limiting
- [ ] CORS handling

**Practice**: Enhance `cmd/api/main.go`
- Add authentication middleware
- Implement rate limiting
- Add structured logging

#### Day 19-20: Optimization & Best Practices
- [ ] Benchmarking
- [ ] Profiling (CPU, memory)
- [ ] Concurrent patterns
- [ ] Context usage
- [ ] Graceful shutdown

**Practice**: Optimize the application
- Add benchmarks for handlers
- Implement connection pooling
- Add context timeouts
- Profile and optimize slow queries

#### Day 21: Deployment & Documentation
- [ ] Binary compilation
- [ ] Docker containerization
- [ ] API documentation
- [ ] README best practices
- [ ] Version management

**Practice**: Prepare for production
- Create a Dockerfile
- Write API documentation
- Add deployment scripts
- Document setup process

## 🔌 API Endpoints

### Tasks

#### Create Task
```bash
POST /api/tasks
Content-Type: application/json

{
  "title": "Task title",
  "description": "Task description",
  "priority": "high|medium|low",
  "due_date": "2024-12-31T23:59:59Z"  # Optional
}

Response: 201 Created
{
  "success": true,
  "data": {
    "id": 1,
    "title": "Task title",
    "description": "Task description",
    "status": "pending",
    "priority": "high",
    "created_at": "2024-01-24T10:00:00Z",
    "updated_at": "2024-01-24T10:00:00Z",
    "due_date": "2024-12-31T23:59:59Z"
  }
}
```

#### List Tasks
```bash
GET /api/tasks?status=pending&priority=high

Response: 200 OK
{
  "success": true,
  "data": [
    {
      "id": 1,
      "title": "Task title",
      ...
    }
  ]
}
```

#### Get Task
```bash
GET /api/tasks/{id}

Response: 200 OK
{
  "success": true,
  "data": {
    "id": 1,
    "title": "Task title",
    ...
  }
}
```

#### Update Task
```bash
PUT /api/tasks/{id}
Content-Type: application/json

{
  "title": "Updated title",           # Optional
  "description": "Updated desc",      # Optional
  "status": "completed",              # Optional: pending|in_progress|completed
  "priority": "low"                   # Optional: low|medium|high
}

Response: 200 OK
{
  "success": true,
  "data": {
    "id": 1,
    "title": "Updated title",
    ...
  }
}
```

#### Delete Task
```bash
DELETE /api/tasks/{id}

Response: 200 OK
{
  "success": true,
  "data": {
    "message": "Task deleted successfully"
  }
}
```

#### Get Statistics
```bash
GET /api/stats

Response: 200 OK
{
  "success": true,
  "data": {
    "total": 10,
    "by_status": {
      "pending": 5,
      "in_progress": 3,
      "completed": 2
    },
    "by_priority": {
      "high": 4,
      "medium": 3,
      "low": 3
    }
  }
}
```

### Health Check
```bash
GET /health

Response: 200 OK
{
  "success": true,
  "data": {
    "status": "healthy"
  }
}
```

## 🧪 Testing

### Run All Tests
```bash
go test ./...
```

### Run Tests with Coverage
```bash
go test -cover ./...
```

### Run Tests Verbosely
```bash
go test -v ./tests/...
```

### Run Specific Test
```bash
go test -run TestCreateTask ./tests/...
```

### Generate Coverage Report
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🔧 Configuration

Environment variables (set in `.env` file):

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `HOST` | `localhost` | Server host |
| `DB_PATH` | `./data/tasks.db` | SQLite database path |
| `ENV` | `development` | Environment (development/production) |

## 🏗️ Architecture & Design Patterns

### Layered Architecture

```
HTTP Request → Handler → Database → Response
     ↓            ↓          ↓
Validation   Business    SQL
              Logic     Queries
```

### Key Patterns Used

1. **Repository Pattern**: Database access abstracted in `pkg/database`
2. **Handler Pattern**: HTTP handlers separated from business logic
3. **Dependency Injection**: Dependencies passed via constructors
4. **Middleware**: Request processing pipeline
5. **Environment Configuration**: External configuration management
6. **Error Wrapping**: Context-rich error messages

### Clean Code Principles

- **Single Responsibility**: Each package has one clear purpose
- **DRY (Don't Repeat Yourself)**: Shared logic in reusable functions
- **KISS (Keep It Simple)**: Simple, readable code over clever code
- **Error Handling**: Explicit error checking and propagation
- **Testing**: Comprehensive test coverage

## 💡 Key Concepts Demonstrated

### 1. Standard Project Layout
Following the community-standard Go project structure for maintainability and clarity.

### 2. RESTful API Design
Proper HTTP methods, status codes, and resource-oriented endpoints.

### 3. Database Integration
SQLite integration with proper connection management and migrations.

### 4. JSON Handling
Request/response serialization with struct tags and custom marshaling.

### 5. Error Handling
Comprehensive error handling with context-rich error messages.

### 6. Testing
Unit and integration tests with table-driven test patterns.

### 7. Configuration Management
Environment-based configuration with sensible defaults.

### 8. Middleware
Request logging and processing pipeline.

### 9. Graceful Shutdown
Proper signal handling for clean server shutdown.

### 10. Package Organization
Clear separation between public API (`pkg/`) and internal code (`internal/`).

## 🎓 Learning Resources

### Official Documentation
- [Go Official Documentation](https://go.dev/doc/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### Project Layout
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Go Project Structure Best Practices](https://go.dev/doc/modules/layout)

### Libraries Used
- [Gorilla Mux](https://github.com/gorilla/mux) - HTTP router
- [Go-SQLite3](https://github.com/mattn/go-sqlite3) - SQLite driver
- [GoDotEnv](https://github.com/joho/godotenv) - Environment loader

### Best Practices
- [Go Best Practices](https://go.dev/doc/effective_go)
- [Go Code Organization](https://go.dev/blog/organizing-go-code)
- [Testing Best Practices](https://go.dev/doc/tutorial/add-a-test)

## 🚀 Next Steps & Enhancements

Once you've completed the 3-week plan, try these enhancements:

### Week 4+: Advanced Features
1. **Authentication & Authorization**
   - Add JWT authentication
   - Implement role-based access control
   - Add user management

2. **Advanced Database**
   - Switch to PostgreSQL
   - Implement migrations
   - Add database seeding
   - Implement soft deletes

3. **Performance**
   - Add caching (Redis)
   - Implement pagination
   - Add database connection pooling
   - Profile and optimize

4. **API Enhancements**
   - Add filtering and sorting
   - Implement search
   - Add bulk operations
   - Version the API (v1, v2)

5. **DevOps**
   - Create Dockerfile
   - Add Docker Compose
   - CI/CD pipeline
   - Kubernetes deployment

6. **Monitoring**
   - Add Prometheus metrics
   - Implement distributed tracing
   - Add health checks
   - Error tracking (Sentry)

7. **Documentation**
   - Generate OpenAPI/Swagger docs
   - Add code examples
   - Create Postman collection
   - Write architecture diagrams

## 🤝 Contributing

This is a learning project. Feel free to:
1. Add new features
2. Improve documentation
3. Fix bugs
4. Add more tests
5. Suggest improvements

## 📄 License

This project is part of the [go-lang-study](https://github.com/smaruf/go-lang-study) repository and follows the same license.

---

## 📌 Tips for Success

1. **Code Daily**: Even 30 minutes a day is better than cramming
2. **Type the Code**: Don't copy-paste; type it yourself to build muscle memory
3. **Experiment**: Break things and fix them; that's how you learn
4. **Read Documentation**: Go's standard library docs are excellent
5. **Test Everything**: Write tests as you go, not after
6. **Ask Questions**: Use Go forums, Discord, and Stack Overflow
7. **Build Projects**: Apply what you learn to your own ideas
8. **Review Code**: Read other people's Go code on GitHub

**Happy Learning! 🎉**
