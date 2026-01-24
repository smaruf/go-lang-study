# Function Examples

This directory contains comprehensive examples of Go functions, including advanced patterns like functional options, channel operations, and practical implementations.

## Examples

### Functional Options Pattern

The functional options pattern is an elegant way to provide optional parameters to functions and constructors in Go.

#### `option_call.go` - Basic Options Pattern

A simple implementation demonstrating:

- **Type-safe options**: Using function types for configuration
- **Default values**: Setting sensible defaults
- **Extensibility**: Easy to add new options without breaking existing code

**Key Concepts:**
```go
type TokenOption func(t *tokenOpts)

func WithAccount(name string) TokenOption {
    return func(t *tokenOpts) {
        t.account = name
    }
}
```

#### `calloptions_exp.go` - Advanced Options Pattern

A production-ready implementation using the `calloptions` library showing:

- **Type-safe options**: Compile-time safety for method-specific options
- **Interface segregation**: Each method has its own option interface
- **Shared options**: Options that work across multiple methods
- **Compile-time checks**: Prevents using incompatible options

**Features:**
- `WithAccount()` - Can be used with both method A and B
- `WithID()` - Only works with method A
- `WithTurnOn()` - Only works with method B

**Usage:**
```bash
go run calloptions_exp.go
```

**Output:**
```
method A received: main.aCallOptions{accountName:"John Doak", id:3}
method B received: main.bCallOptions{accountName:"David Luyer", turnOn:true}
```

**Key Concepts:**
- Embedded interfaces for type-safe options
- Using the `calloptions` package for validation
- Method-specific option interfaces
- Shared options across multiple methods

#### `call_options.go` - Multi-Method Options

Demonstrates applying options to multiple authentication methods:

- **ByToken()** - Token-based authentication with options
- **ByUserPass()** - Username/password authentication with options
- **Shared configuration**: Account settings across methods

### Channel Operations

#### `channel_101.go` - Basic Channel Usage

A fundamental example of channel send and receive operations:

**Features:**
- Creating unbuffered channels
- Sending data through channels
- Receiving data from channels
- Goroutine communication

**Usage:**
```bash
go run channel_101.go
```

**Output:**
```
start Main method
257
End Main method
```

**Key Concepts:**
- Channel creation with `make(chan int)`
- Sending with `ch <- 23`
- Receiving with `<-ch`
- Goroutines and synchronization

#### `close_channel_101.go` - Channel Closing

Demonstrates proper channel closing patterns:

**Features:**
- Closing channels with `close()`
- Checking if a channel is closed
- Using for-range loops with channels
- Channel state detection

**Usage:**
```bash
go run close_channel_101.go
```

**Key Concepts:**
- The `ok` idiom for checking channel state
- Producer-consumer pattern
- Proper channel closing
- Avoiding deadlocks

### Sorting Algorithms

#### `insertion_vs_merge_sort.go` - Algorithm Implementation

Implementations of two sorting algorithms for comparison:

**Algorithms:**
1. **Insertion Sort**: O(n²) time complexity, good for small datasets
2. **Merge Sort**: O(n log n) time complexity, better for larger datasets

**Usage:**
```bash
go run insertion_vs_merge_sort.go
```

**Key Concepts:**
- In-place sorting vs divide-and-conquer
- Time complexity trade-offs
- Algorithm selection based on data size

#### `brenchmark_test_sorting.go` - Performance Benchmarking

Comprehensive benchmark tests comparing sorting algorithms:

**Features:**
- Multiple input sizes (10, 100, 1000, 10000, 100000)
- Performance comparison between insertion and merge sort
- Benchmark naming conventions
- Random data generation for testing

**Usage:**
```bash
go test -bench=. -benchmem
```

**Expected Results:**
- Insertion sort faster for small arrays (< 100 elements)
- Merge sort faster for larger arrays
- Clear performance crossover point

**Key Concepts:**
- Go benchmark functions (`BenchmarkXxx`)
- `testing.B` for benchmark control
- `b.ResetTimer()` to exclude setup time
- Sub-benchmarks with `b.Run()`

### Sudoku Solver

#### `suduku_solve.go` - Backtracking Algorithm

A complete Sudoku solver implementation using backtracking:

**Features:**
- Backtracking algorithm
- Rule validation (row, column, 3x3 grid)
- Board representation using 2D byte slices
- Pretty-printed output

**Usage:**
```bash
go run suduku_solve.go
```

**Key Concepts:**
- Recursive backtracking
- Constraint satisfaction problems
- Nested functions and closures
- Board state management

## Running the Examples

```bash
# Navigate to the function directory
cd src/function

# Run individual examples
go run calloptions_exp.go
go run channel_101.go
go run close_channel_101.go
go run insertion_vs_merge_sort.go
go run suduku_solve.go

# Run benchmarks
go test -bench=. -benchmem

# Run benchmarks with specific pattern
go test -bench=BenchmarkInsertionSort -benchmem
go test -bench=BenchmarkMergeSort -benchmem
```

## Learning Objectives

After studying these examples, you'll understand:

1. **Functional Options Pattern**: How to implement flexible, extensible APIs
2. **Type Safety**: Using interfaces for compile-time safety
3. **Channel Operations**: Sending, receiving, and closing channels properly
4. **Goroutine Communication**: Coordinating concurrent operations
5. **Algorithm Implementation**: Writing efficient sorting algorithms
6. **Benchmarking**: Measuring and comparing performance
7. **Backtracking**: Solving constraint satisfaction problems
8. **Code Organization**: Structuring Go programs effectively

## Best Practices Demonstrated

### Functional Options
- Use function types for optional parameters
- Provide sensible defaults
- Make options composable and type-safe
- Use interfaces to restrict which options work with which methods

### Channels
- Always close channels when done producing
- Use the `ok` idiom to check if channels are closed
- Prefer unbuffered channels unless you need buffering
- Avoid closing channels from the receiver side

### Benchmarking
- Exclude setup code with `b.ResetTimer()`
- Test with multiple input sizes
- Use `b.N` for iteration count
- Add `-benchmem` to see memory allocations

### Algorithm Design
- Choose algorithms based on data characteristics
- Consider time vs space complexity trade-offs
- Use benchmarks to validate performance assumptions
- Document algorithm complexity

## Common Patterns

1. **Functional Options**: Flexible API design without breaking changes
2. **Channel Communication**: Safe concurrent data exchange
3. **Benchmarking**: Performance measurement and comparison
4. **Backtracking**: Systematic search for solutions
5. **Type Safety**: Using interfaces to enforce constraints at compile time

## Additional Resources

- [Functional Options Pattern](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis) - Dave Cheney's article
- [Go Concurrency Patterns](https://go.dev/blog/pipelines) - Official Go blog
- [Effective Go](https://go.dev/doc/effective_go) - Official Go documentation
- [Go by Example: Channels](https://gobyexample.com/channels) - Practical examples

## Dependencies

Some examples require external packages:

```bash
# For calloptions_exp.go
go get github.com/johnsiilver/calloptions
```

## Testing

```bash
# Run benchmark tests
go test -bench=. -benchmem

# Run with race detection
go test -race

# Run specific benchmark
go test -bench=BenchmarkMergeSort -benchmem
```
