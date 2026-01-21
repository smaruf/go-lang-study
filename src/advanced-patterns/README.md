# Advanced Go Patterns

This directory contains advanced Go programming patterns and best practices.

## Concurrency Patterns

### Worker Pool
- **File**: `worker_pool.go`
- **Description**: Demonstrates efficient worker pool pattern for processing tasks concurrently
- **Concepts**: Goroutines, Channels, WaitGroups, Context

### Pipeline Pattern
- **File**: `pipeline.go`
- **Description**: Shows how to build concurrent pipelines for data processing
- **Concepts**: Channels, Fan-out/Fan-in, Stages

### Rate Limiter
- **File**: `rate_limiter.go`
- **Description**: Token bucket rate limiter implementation
- **Concepts**: Channels, Time, Concurrency control

## Design Patterns

### Factory Pattern
- **File**: `factory.go`
- **Description**: Factory pattern for object creation
- **Concepts**: Interfaces, Encapsulation

### Singleton Pattern
- **File**: `singleton.go`
- **Description**: Thread-safe singleton implementation
- **Concepts**: sync.Once, Initialization

### Observer Pattern
- **File**: `observer.go`
- **Description**: Event notification pattern
- **Concepts**: Interfaces, Event handling

### Strategy Pattern
- **File**: `strategy.go`
- **Description**: Runtime algorithm selection
- **Concepts**: Interfaces, Polymorphism

## Testing Patterns

### Table-Driven Tests
- **File**: `table_driven_test.go`
- **Description**: Comprehensive table-driven testing examples
- **Concepts**: Testing, Sub-tests

### Mocking
- **File**: `mocking_test.go`
- **Description**: Interface-based mocking patterns
- **Concepts**: Testing, Interfaces, Dependency injection

## Running Examples

Each file can be run independently:

```bash
go run worker_pool.go
go run pipeline.go
go run rate_limiter.go
go run factory.go
go run singleton.go
go run observer.go
go run strategy.go
```

## Running Tests

```bash
go test -v .
```

## Key Concepts Covered

- **Concurrency**: Advanced goroutine patterns, synchronization
- **Design Patterns**: Common GoF patterns adapted for Go
- **Testing**: Best practices for testing Go applications
- **Performance**: Efficient resource usage and optimization
