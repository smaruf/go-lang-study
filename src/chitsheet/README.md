# Cheatsheet

This directory contains Go code examples, cheatsheets, and quick reference materials for various Go concepts and patterns.

## Contents

### Documentation
- [CLI Commands](CLI.md) - Common Go CLI commands and usage

### Code Examples

#### Basic Concepts
- [Arrays](arrays.go) - Array manipulation and operations
- [Maps](maps.go) - Map usage and operations
- [Functions](functions.go) - Function examples and patterns
- [Sum Arrays](sum_arrays.go) - Array summation examples

#### Concurrency
- [Concurrency](concurrency.go) - Concurrency patterns and examples
- [Routine Leaks Memory](routine_leaks_mem.go) - Examples of goroutine memory leaks

#### Memory Management
- [GC GoDEBUG Memory](gc_godebug_mem.go) - Garbage collection and memory debugging
- [Memory GC Check](mem_gc_check.go) - Memory and GC checking utilities
- [Simulate Memory Leak Defer](simulate_memory_leak_defer.go) - Memory leak simulations with defer

#### WebSocket
- [Client WebSocket](client_websocket.go) - WebSocket client implementation
- [Router WebSocket](router_websocket.go) - WebSocket router implementation

#### Testing
- [Testing](testing.go) - Testing examples and patterns
- [Simple List Test](simple_list_test.go) - Simple list testing examples
- [Table Test](table_test.go) - Table-driven test examples
- [Test Implementation Function](test_impl_func.go) - Test implementation patterns

#### Other
- [Connect OpenTelemetry](connect_open_telimetry.go) - OpenTelemetry integration example
- [Minimum Remove Parenthesis](minimum_remove_parenthesis.go) - Algorithm example

## Getting Started

You can run any of the Go files directly:

```bash
go run <filename>.go
```

For test files:

```bash
go test <filename>_test.go
```

## Useful Links

- [Go Programming Language](https://golang.org)
- [Go Documentation](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go)
- [GitHub Repository](https://github.com/smaruf/go-lang-study)

## Repository Structure

Browse the complete cheatsheet directory on GitHub:
- [Cheatsheet Directory](https://github.com/smaruf/go-lang-study/tree/master/src/chitsheet)

## Contributions

Feel free to fork this repository and submit pull requests with additional examples or improvements.
