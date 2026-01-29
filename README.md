# Go Language Study Repository

A comprehensive collection of Go programming examples, from basic concepts to advanced patterns, designed for learning and reference.

## 🎯 Goal

Exploring Go features through practical implementations, modern best practices, and real-world application patterns.

## 📚 Table of Contents

- [Quick Start](#-quick-start)
- [Environment Setup](#-environment-setup)
- [Project Structure](#-project-structure)
- [Learning Path](#-learning-path)
- [Projects Index](#-projects-index)
- [Core Concepts](#-core-concepts)
- [Advanced Topics](#-advanced-topics)
- [Best Practices](#-best-practices)
- [Contributing](#-contributing)

## 🚀 Quick Start

1. **Clone the repository:**
   ```bash
   git clone https://github.com/smaruf/go-lang-study.git
   cd go-lang-study
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Set up environment:**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Run a sample project:**
   ```bash
   cd src/echo
   go run simple_echo.go
   ```

## 🔧 Environment Setup

### Prerequisites
- Go 1.21 or later ([Installation Guide](https://golang.org/doc/install))
- Git
- Code editor (VS Code, GoLand, or similar)

### Environment Configuration
1. Copy `.env.example` to `.env` and configure your settings
2. Install Go tools for development:
   ```bash
   go install golang.org/x/tools/cmd/goimports@latest
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

## 📁 Project Structure

```
go-lang-study/
├── src/                     # Source code examples
│   ├── basic/              # Basic Go concepts
│   ├── web/                # Web applications
│   ├── concurrency/        # Concurrency patterns
│   ├── cli/                # Command-line tools
│   ├── data-processing/    # Data manipulation examples
│   └── projects/           # Full-featured projects
├── docs/                   # Documentation
├── .env.example           # Environment configuration template
├── go.mod                 # Go module definition
└── README.md              # This file
```

## 🎓 Learning Path

### Beginner (Start Here)
1. **Basic Syntax** - Variables, functions, control structures
2. **Data Structures** - Arrays, slices, maps, structs
3. **Error Handling** - Go's error handling patterns
4. **File I/O** - Reading and writing files

### Intermediate
1. **Concurrency** - Goroutines and channels
2. **HTTP/Web** - Building web servers and APIs
3. **Testing** - Unit testing and benchmarking
4. **Database** - SQL and NoSQL interactions

### Advanced
1. **Microservices** - Distributed systems patterns
2. **Performance** - Profiling and optimization
3. **Cloud Integration** - AWS, GCP services
4. **System Programming** - Low-level Go programming

## 🗂️ Projects Index

### 🌟 Full-Featured Projects
| Project | Description | Concepts Covered |
|---------|-------------|------------------|
| [✈️ Remote Aircraft](src/remote-aircraft/) | Parametric aircraft design system for FPV drones & fixed-wing | Engineering calculations, CLI design, JSON export |
| [📈 NASDAQ Trading Simulator](src/nasdaq-cse/) | Complete trading system with WebSocket, AI bot | REST API, WebSocket, GORM, Testing |
| [🤖 GoBot Collection](src/gobot/) | IoT and bot integrations | Hardware control, Discord API, Environment variables |
| [💰 Wallet Service](src/wallet/) | gRPC-based wallet system | gRPC, Protocol Buffers, Microservices |
| [🎨 UI Applications](src/ui-app/) | Desktop GUI applications with Fyne | GUI, Fyne framework, Event handling |
| [🔧 Embedded OS](src/embedded-os/) | Minimal OS for Raspberry Pi & Arduino | TinyGo, Embedded systems, GPIO, Sensors |
| [⚙️ FreeRTOS Systems](src/tiny/freeRTOS/) | Robotics, rocketry, and renewable energy monitoring | Real-time OS, Motor control, Telemetry, MPPT |
| [🚀 Advanced Patterns](src/advanced-patterns/) | Advanced Go patterns and best practices | Design patterns, Concurrency, Testing |

### 🌐 Web Development
| Example | Description | Framework/Library |
|---------|-------------|-------------------|
| [Echo Server](src/echo/) | Simple HTTP server | Echo framework |
| [Gin API](src/gin/) | REST API with middleware | Gin framework |
| [HTTP/2 Server](src/http2_server.go) | HTTP/2 implementation | Standard library |
| [Web Server](src/web-server/) | Basic web server examples | Standard library |
| [WebSocket](src/ws/) | WebSocket communication | WebSocket library |
| [Broadcast](src/broadcast/) | Broadcasting patterns | Standard library |

### ⚡ Concurrency & Performance
| Example | Description | Concepts |
|---------|-------------|----------|
| [Goroutine Workers](src/goroutine/) | Worker pool patterns | Goroutines, Channels |
| [Concurrency Patterns](src/concurrency/) | Advanced concurrency patterns | Channels, Select, Context |
| [Atomic Operations](src/atomic_worker.go) | Thread-safe counters | sync/atomic |
| [Atomic Increment](src/atomic_increament.go) | Atomic increment operations | sync/atomic |
| [Mutex Examples](src/mutex-count.go) | Synchronization patterns | Mutexes, RWMutex |
| [RWMutex](src/rw_mutex.go) | Read-Write mutex patterns | sync.RWMutex |
| [Goroutine with Channel](src/goroutin_with_channel.go) | Basic goroutine communication | Goroutines, Channels |
| [Sample Channel](src/sample_channel.go) | Channel usage examples | Channels |
| [Sample Goroutine](src/sample_go_routine.go) | Basic goroutine patterns | Goroutines |

### 🛠️ System & Tools
| Example | Description | Use Case |
|---------|-------------|----------|
| [CLI Tool](src/cli-tool/) | Command-line application | CLI development |
| [Context Demo](src/context_demo.go) | Context usage patterns | Cancellation, Timeouts |
| [Custom Middleware](src/custom_middleware.go) | HTTP middleware | Request/Response processing |
| [Retry Mechanism](src/retry_mech.go) | Fault tolerance | Error handling |
| [Retry Main](src/retry-main.go) | Retry pattern implementation | Error handling |
| [Panic Recover](src/panic_recover_with_error.go) | Error recovery patterns | Panic/Recover |
| [Reverse Proxy](src/reverse_proxy_server.go) | Reverse proxy server | HTTP proxy |
| [Upstream Server](src/upstream_server.go) | Backend server example | HTTP server |

### 🧮 Algorithms & Data Structures
| Example | Description | Algorithm |
|---------|-------------|-----------|
| [GCD Calculator](src/calculate_gcd.go) | Greatest Common Divisor | Euclidean algorithm |
| [Palindrome Check](src/palindrome_str.go) | String manipulation | String algorithms |
| [Radix Sort](src/radix_sort.go) | Sorting algorithm | Non-comparison sorting |
| [Generate Combination](src/generate_combination.go) | Combination generation | Combinatorics |
| [Calculator App](src/calculator_app.go) | Expression calculator | Parsing, Evaluation |
| [Rubik Cube Solver](src/rubik_cube_solver.go) | Rubik's cube solving | Graph search |

### 🎮 Game & Graphics
| Example | Description | Technology |
|---------|-------------|------------|
| [Game Engine](src/game-engine/) | 3D game engine with G3N | G3N framework |
| [Ping Pong Ball](src/ping_pong_ball.go) | Game simulation | Goroutines |
| [Ping Pong 2 Min](src/ping_pong_ball_2_min.go) | Timed game simulation | Goroutines, Time |

### 🔌 IoT & Embedded Systems
| Example | Description | Platform |
|---------|-------------|----------|
| [TinyGo Projects](src/tiny/) | Embedded systems examples | TinyGo, RPi, Arduino |
| [FreeRTOS Examples](src/tiny/freeRTOS/) | Real-time OS with robotics, rocketry & energy | FreeRTOS, TinyGo |
| [🤖 Robotics Systems](src/tiny/freeRTOS/robotics/) | Motor control, sensors, autonomous navigation | RPi Pico, Arduino |
| [🚀 Rocketry Control](src/tiny/freeRTOS/rocketry/) | Launch control, telemetry, flight computer | RPi Pico, Arduino |
| [🔋 Energy Monitoring](src/tiny/freeRTOS/energy/) | Wind, solar, hydro, thermoelectric generators | RPi Pico, Arduino |
| [SR-71 Simulator](src/tiny/sr71sim/) | Aircraft simulation | TinyGo |
| [TinyGo Blinky](src/tinygo_blinky.go) | LED blink example | TinyGo |
| [TinyGo PWM](src/tinygo_pwm.go) | PWM control example | TinyGo |
| [GoBot Hello](src/gobot_hellow.go) | IoT framework intro | Gobot |
| [GoBot Collection](src/gobot/) | IoT and bot integrations | Gobot framework |

### ☁️ Cloud & Infrastructure
| Example | Description | Service |
|---------|-------------|---------|
| [AWS Examples](src/aws/) | AWS services integration | AWS SDK |
| [S3 Bucket](src/aws/bucket/) | S3 operations | AWS S3 |
| [Lambda Functions](src/aws/lambda/) | Serverless functions | AWS Lambda |
| [CloudFormation](src/aws/cloudFormation/) | Infrastructure as code | CloudFormation |
| [Docker Examples](src/aws/docker/) | Containerization | Docker |
| [SAM](src/aws/sam/) | Serverless application model | AWS SAM |

### 🗄️ Database & Storage
| Example | Description | Database |
|---------|-------------|----------|
| [GORM SQLite](src/gorm-sqllite.go) | ORM with SQLite | GORM |
| [CockroachDB](src/cockroachdb/) | Distributed SQL | CockroachDB |
| [Data Processing](src/data-processing/) | Data manipulation | Standard library |
| [Processing](src/processing/) | Data processing patterns | Standard library |
| [Simple CRUD](src/simple_crud.go) | Basic CRUD operations | Standard library |

### 🏗️ Architecture & Patterns
| Example | Description | Pattern |
|---------|-------------|---------|
| [Microservices](src/microservices/) | Microservice architecture | Clean architecture |
| [User Service](src/microservices/user-service/) | User management service | Microservices |
| [Interface Sample](src/Interface_sample.go) | Interface patterns | Interfaces |
| [Pointer Struct Interface](src/pointer_struct_interface.go) | Advanced interface usage | Interfaces, Pointers |
| [Pointer Struct Interface Closure](src/pointer_struct_interface_closure.go) | Closures with interfaces | Closures |
| [Struct Method](src/struct_method.go) | Method declarations | Methods |
| [Sample Pointer](src/sample_pointer.go) | Pointer usage | Pointers |

### 🔐 Security & Authentication
| Example | Description | Technology |
|---------|-------------|------------|
| [Simple JWT](src/simple_jwt.go) | JWT authentication | JWT |
| [Jaeger Tracer](src/jagger_tracer.go) | Distributed tracing | Jaeger |

### 🌐 Networking
| Example | Description | Protocol |
|---------|-------------|----------|
| [Socket Mesh Net](src/socket_mesh_net.go) | Mesh networking | TCP/UDP |
| [Go Socket Client](src/go_socket_client.go) | Socket client | TCP |
| [Sample Subdomain](src/sample_subdomain.go) | Subdomain handling | HTTP |

### 📚 Learning Resources
| Example | Description | Topic |
|---------|-------------|-------|
| [Cheat Sheet](src/chitsheet/) | Go language cheat sheet | Quick reference |
| [3 Weeks Plan](src/3-weeks-plan/) | Structured learning plan | Study guide |
| [Function Examples](src/function/) | Function patterns | Functions |
| [Named Return Values](src/name_return_valued_func.go) | Named returns | Functions |
| [Nested Functions](src/nested_function_op.go) | Function composition | Functions |

## 💡 Core Concepts

### Language Fundamentals
- **Type System**: Static typing, interfaces, type assertions
- **Memory Management**: Pointers, garbage collection
- **Error Handling**: Error interface, panic/recover
- **Concurrency**: Goroutines, channels, select statements

### Standard Library Usage
- **net/http**: HTTP client/server programming
- **encoding/json**: JSON marshaling/unmarshaling
- **database/sql**: Database interactions
- **context**: Request-scoped values and cancellation

### Modern Go Practices
- **Modules**: Dependency management with go.mod
- **Testing**: Unit tests, benchmarks, table-driven tests
- **Logging**: Structured logging with logrus
- **Configuration**: Environment-based configuration

## 🏗️ Advanced Topics

### Architectural Patterns
- **Clean Architecture**: Separation of concerns
- **Hexagonal Architecture**: Ports and adapters
- **Microservices**: Distributed system design
- **Event-Driven**: Pub/sub patterns

### Performance & Optimization
- **Profiling**: CPU and memory profiling
- **Benchmarking**: Performance measurement
- **Optimization**: Memory allocation, garbage collection
- **Concurrency Patterns**: Worker pools, fan-in/fan-out

### Cloud & Infrastructure
- **Container Integration**: Docker, Kubernetes
- **Cloud Services**: AWS Lambda, GCP functions
- **Monitoring**: Metrics, tracing, observability
- **CI/CD**: GitHub Actions, automated testing

## ✅ Best Practices

### Code Quality
- Use `gofmt` for consistent formatting
- Follow effective Go guidelines
- Write clear, self-documenting code
- Implement comprehensive error handling

### Testing Strategy
- Unit tests for business logic
- Integration tests for external dependencies
- Benchmarks for performance-critical code
- Table-driven tests for multiple scenarios

### Security Considerations
- Input validation and sanitization
- Secure credential management
- HTTPS/TLS configuration
- SQL injection prevention

## 🧪 Testing

Run all tests across the repository:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

Run benchmarks:
```bash
go test -bench=. ./...
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Follow Go conventions and best practices
4. Add tests for new functionality
5. Ensure all tests pass (`go test ./...`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Development Guidelines
- Follow the existing project structure
- Include comprehensive documentation
- Add tests for new features
- Use environment variables for configuration
- Format code with `gofmt`

## 📄 License

This project is open source and available under the [MIT License](LICENSE).

---

**Start Learning**: Choose a project from the index above and dive into Go! 🚀

For questions or suggestions, please open an issue or contribute to the discussion.
