# SR-71 Simulator - Quick Start Guide

## Getting Started in 5 Minutes

### Prerequisites
- Go 1.21 or later
- Git

### Installation

```bash
# Clone the repository
git clone https://github.com/smaruf/go-lang-study.git
cd go-lang-study/src/embedded-os/tiny/sr71sim

# Initialize Go module (if not already done)
go mod init github.com/smaruf/go-lang-study/src/embedded-os/tiny/sr71sim

# Install dependencies
go get github.com/go-echarts/go-echarts/v2
go get github.com/gin-gonic/gin
go get github.com/spf13/cobra

# Build the simulator
go build -o sr71sim

# Run the simulator
./sr71sim
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run specific module tests
go test ./avionics
go test ./engine
go test ./flying

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Usage Examples

### 1. Console Mode (Interactive)

```bash
# Start interactive console
./sr71sim console

# Or use the CLI commands
./sr71sim start --mission reconnaissance --duration 30m
./sr71sim status
./sr71sim stop
```

### 2. Web Server Mode

```bash
# Start web server
./sr71sim serve --port 8080

# Access in browser
# http://localhost:8080

# API Examples:
curl http://localhost:8080/api/v1/simulation/status
curl -X POST http://localhost:8080/api/v1/simulation/start
```

### 3. Programmatic Usage

```go
package main

import (
    "fmt"
    "github.com/smaruf/go-lang-study/src/embedded-os/tiny/sr71sim/engine"
    "github.com/smaruf/go-lang-study/src/embedded-os/tiny/sr71sim/avionics"
    "github.com/smaruf/go-lang-study/src/embedded-os/tiny/sr71sim/flying"
)

func main() {
    // Create new SR-71 aircraft
    aircraft := flying.NewSR71()
    
    // Start simulation
    aircraft.FlyAtHeight(80000) // Fly at 80,000 ft
    aircraft.AdjustVelocityForMission("reconnaissance")
    
    // Get engine state
    eng := engine.New()
    eng.SetSpeed(2200) // Mach 3+
    fmt.Printf("Engine Mode: %s\n", eng.Mode())
    
    // Check avionics
    avio := avionics.New()
    state := avio.GetState()
    fmt.Printf("Altitude: %.0f ft, Speed: %.0f mph\n", 
        state.Altitude, state.Speed)
}
```

## Configuration

### Environment Variables

Create a `.env` file:

```bash
# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# Database
DB_TYPE=sqlite
DB_PATH=./sr71sim.db

# Simulation
SIM_TICK_RATE=60
SIM_MAX_CONCURRENT=10

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

### Config File

Create `config.yaml`:

```yaml
server:
  host: "0.0.0.0"
  port: 8080

simulation:
  tick_rate: 60
  default_altitude: 75000
  default_speed: 2000

logging:
  level: "info"
  format: "json"
```

## Platform-Specific Setup

### Web Platform

```bash
# Build web assets
cd web
npm install
npm run build

# Start server
cd ..
./sr71sim serve
```

### Raspberry Pi

```bash
# Cross-compile for RPI
GOOS=linux GOARCH=arm GOARM=7 go build -o sr71sim-rpi

# Copy to RPI
scp sr71sim-rpi pi@raspberrypi.local:~/
scp config.yaml pi@raspberrypi.local:~/

# SSH into RPI and run
ssh pi@raspberrypi.local
./sr71sim-rpi console
```

### Arduino

```bash
# Flash Arduino firmware
cd firmware/arduino
arduino-cli compile --fqbn arduino:avr:mega sr71_client
arduino-cli upload -p /dev/ttyACM0 --fqbn arduino:avr:mega sr71_client

# Connect Arduino to host
./sr71sim --arduino /dev/ttyACM0
```

## Docker Deployment

```bash
# Build Docker image
docker build -t sr71sim:latest .

# Run container
docker run -p 8080:8080 sr71sim:latest

# Or use docker-compose
docker-compose up -d
```

## Development

### Project Structure

```
sr71sim/
├── avionics/          # Avionics module
├── engine/            # Engine simulation
├── flying/            # Flight control
├── fueling/           # Fuel management
├── landing/           # Landing operations
├── texing/            # Taxi operations
├── api/               # REST API handlers (new)
├── web/               # Web UI (new)
├── cli/               # CLI commands (new)
├── firmware/          # Arduino/embedded (new)
│   └── arduino/
├── config/            # Configuration files
├── docs/              # Documentation
├── main.go            # Entry point
├── config.yaml        # Configuration
├── Dockerfile         # Docker build
├── docker-compose.yml # Multi-container setup
└── README.md
```

### Adding New Features

1. Create new module in appropriate directory
2. Add tests in `*_test.go` files
3. Update API if needed
4. Update documentation
5. Run tests: `go test ./...`
6. Build: `go build`

### Code Style

```bash
# Format code
go fmt ./...

# Lint code
golangci-lint run

# Vet code
go vet ./...
```

## Troubleshooting

### Common Issues

**1. Import errors**
```bash
# Solution: Initialize module
go mod init github.com/smaruf/go-lang-study/src/embedded-os/tiny/sr71sim
go mod tidy
```

**2. Port already in use**
```bash
# Solution: Use different port
./sr71sim serve --port 8081
```

**3. Permission denied (RPI/Arduino)**
```bash
# Solution: Add user to dialout group
sudo usermod -a -G dialout $USER
# Logout and login again
```

**4. Database locked**
```bash
# Solution: Close other instances
pkill sr71sim
rm sr71sim.db-journal
```

## Performance Tuning

### Optimization Tips

1. **Increase tick rate for smoother simulation**
   ```yaml
   simulation:
     tick_rate: 120  # Hz
   ```

2. **Enable caching**
   ```yaml
   cache:
     enabled: true
     ttl: 300  # seconds
   ```

3. **Reduce telemetry frequency**
   ```yaml
   telemetry:
     sample_rate: 10  # Hz instead of 60
   ```

4. **Use binary protocol for Arduino**
   ```yaml
   arduino:
     protocol: binary  # instead of text
   ```

## Next Steps

1. **Read Documentation**
   - [Architecture Guide](ARCHITECTURE.md)
   - [Production Guide](PRODUCTION_GUIDE.md)
   - [API Documentation](docs/API.md)

2. **Try Examples**
   - Basic simulation
   - Web interface
   - CLI commands
   - Hardware integration

3. **Explore Advanced Features**
   - Mission creation
   - Multi-user support
   - Real-time visualization
   - Hardware control

4. **Contribute**
   - Report issues
   - Submit PRs
   - Improve documentation
   - Share your setup

## Support

- **Issues:** https://github.com/smaruf/go-lang-study/issues
- **Discussions:** https://github.com/smaruf/go-lang-study/discussions
- **Email:** [Add contact email]

## License

[Add license information]

---

**Ready to fly? Start your first simulation:**

```bash
./sr71sim console
> start
> status
> help
```

Happy flying! ✈️🚀
