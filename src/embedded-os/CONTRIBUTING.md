# Contributing to Embedded OS Projects

Thank you for your interest in contributing to our embedded systems projects! This guide will help you get started.

## Table of Contents

1. [Code of Conduct](#code-of-conduct)
2. [Getting Started](#getting-started)
3. [Development Setup](#development-setup)
4. [Contribution Guidelines](#contribution-guidelines)
5. [Testing Requirements](#testing-requirements)
6. [Documentation Standards](#documentation-standards)
7. [Pull Request Process](#pull-request-process)
8. [Hardware Testing](#hardware-testing)

## Code of Conduct

- Be respectful and inclusive
- Provide constructive feedback
- Focus on the code, not the person
- Help others learn and grow

## Getting Started

### Prerequisites

Before contributing, ensure you have:

- Go 1.21 or later
- TinyGo 0.30.0 or later
- Git for version control
- Access to target hardware (for hardware-specific contributions)

### Fork and Clone

```bash
# Fork the repository on GitHub
# Then clone your fork
git clone https://github.com/YOUR_USERNAME/go-lang-study.git
cd go-lang-study

# Add upstream remote
git remote add upstream https://github.com/smaruf/go-lang-study.git
```

## Development Setup

### Environment Configuration

```bash
# Install dependencies
go mod download

# Install TinyGo
# Linux:
wget https://github.com/tinygo-org/tinygo/releases/download/v0.30.0/tinygo_0.30.0_amd64.deb
sudo dpkg -i tinygo_0.30.0_amd64.deb

# macOS:
brew install tinygo

# Verify installation
go version
tinygo version
```

### Editor Setup

Recommended VS Code extensions:
- Go (official)
- TinyGo
- GitLens
- Markdown All in One

## Contribution Guidelines

### Types of Contributions

We welcome:

1. **Bug Fixes**
   - Fix incorrect behavior
   - Improve error handling
   - Enhance stability

2. **New Features**
   - New hardware support
   - Additional sensors/actuators
   - Enhanced algorithms

3. **Documentation**
   - Improve README files
   - Add code comments
   - Create tutorials

4. **Performance Improvements**
   - Optimize memory usage
   - Reduce power consumption
   - Improve execution speed

### Coding Standards

#### Go Code Style

Follow standard Go conventions:

```go
// Use meaningful names
func calculateMotorSpeed(targetRPM int) int {
    // Implementation
}

// Comment exported functions
// CalculateDistance computes the distance between two GPS coordinates
func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
    // Implementation
}

// Use error handling
if err := sensor.Read(); err != nil {
    return fmt.Errorf("failed to read sensor: %w", err)
}
```

#### File Organization

```
project/
├── main.go              # Main entry point
├── sensor.go            # Sensor interface
├── motor.go             # Motor control
├── navigation.go        # Navigation logic
├── sensor_test.go       # Sensor tests
├── motor_test.go        # Motor tests
└── README.md            # Project documentation
```

#### Hardware-Specific Code

Document hardware requirements clearly:

```go
// Pin Configuration for Raspberry Pi Pico
const (
    // Motor control pins
    MotorLeftPWM  = machine.GP0  // Left motor PWM (Pin 1)
    MotorRightPWM = machine.GP1  // Right motor PWM (Pin 2)
    
    // Sensor pins
    UltrasonicTrig = machine.GP2 // Ultrasonic trigger (Pin 4)
    UltrasonicEcho = machine.GP3 // Ultrasonic echo (Pin 5)
)
```

### Commit Message Format

Use semantic commit messages:

```
type(scope): subject

body (optional)

footer (optional)
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Build process or auxiliary tool changes

Examples:

```
feat(robotics): add obstacle detection algorithm

Implement ultrasonic-based obstacle detection with adjustable
sensitivity. Tested on Raspberry Pi Pico with HC-SR04 sensor.

Closes #123
```

```
fix(energy): correct solar panel voltage calculation

Fixed integer overflow in voltage calculation that occurred
at high irradiance levels.

Fixes #456
```

## Testing Requirements

### Unit Tests

All new code must include unit tests:

```go
func TestMotorSpeed(t *testing.T) {
    tests := []struct {
        name     string
        input    int
        expected int
    }{
        {"Zero speed", 0, 0},
        {"Half speed", 50, 128},
        {"Full speed", 100, 255},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := calculatePWM(tt.input)
            if result != tt.expected {
                t.Errorf("got %d, want %d", result, tt.expected)
            }
        })
    }
}
```

### Hardware Tests

For hardware-specific code:

```go
// +build hardware

// TestMotorHardware tests motor control on actual hardware
// Run with: go test -tags=hardware
func TestMotorHardware(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping hardware test")
    }
    
    // Hardware test implementation
}
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run with race detection
go test -race ./...

# Run hardware tests
go test -tags=hardware ./...
```

## Documentation Standards

### README Structure

Each project should have a README with:

1. **Title and Description**
2. **Features**
3. **Hardware Requirements**
4. **Wiring Diagram** (if applicable)
5. **Installation**
6. **Usage**
7. **Configuration**
8. **Troubleshooting**
9. **Contributing**
10. **License**

### Code Comments

```go
// Package motor provides motor control functionality for differential
// drive robots using PWM-controlled DC motors.
package motor

// Motor represents a DC motor with PWM speed control
type Motor struct {
    pin      machine.Pin // PWM output pin
    maxSpeed uint16      // Maximum PWM duty cycle (0-65535)
}

// NewMotor creates a new Motor instance
//
// Parameters:
//   - pin: GPIO pin connected to motor controller
//   - maxSpeed: Maximum PWM value (typically 65535)
//
// Returns:
//   - *Motor: Configured motor instance
func NewMotor(pin machine.Pin, maxSpeed uint16) *Motor {
    // Implementation
}
```

### Inline Documentation

Use inline comments for complex logic:

```go
// Calculate PID correction with anti-windup
error := targetSpeed - currentSpeed

// Proportional term
pTerm := kp * error

// Integral term with anti-windup (clamping)
integral += error * dt
if integral > maxIntegral {
    integral = maxIntegral
} else if integral < -maxIntegral {
    integral = -maxIntegral
}
iTerm := ki * integral

// Derivative term with filtering
dTerm := kd * (error - lastError) / dt

// PID output
output := pTerm + iTerm + dTerm
```

## Pull Request Process

### Before Submitting

1. **Update from upstream**
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Run all tests**
   ```bash
   go test ./...
   go test -race ./...
   ```

3. **Run linter**
   ```bash
   golangci-lint run
   ```

4. **Update documentation**
   - Update README if needed
   - Add/update code comments
   - Update CHANGELOG

### Submitting Pull Request

1. Push your changes to your fork
   ```bash
   git push origin feature-branch
   ```

2. Create pull request on GitHub

3. Fill out PR template:
   - Clear description of changes
   - Reference related issues
   - List of changes
   - Testing performed
   - Hardware tested (if applicable)

### PR Review Process

1. **Automated Checks**
   - CI/CD pipeline must pass
   - Code coverage must not decrease
   - Linter must pass

2. **Code Review**
   - At least one approval required
   - Address all review comments
   - Update PR as needed

3. **Merge**
   - Squash and merge for single commits
   - Merge commit for feature branches
   - Delete branch after merge

## Hardware Testing

### Testing Checklist

Before submitting hardware-related PRs:

- [ ] Code compiles for target hardware
- [ ] Tested on actual hardware
- [ ] Pin assignments documented
- [ ] Power requirements specified
- [ ] Wiring diagram included (if new)
- [ ] Safety considerations documented
- [ ] Performance measurements included

### Test Documentation

Include test results in PR:

```markdown
## Hardware Test Results

**Board**: Raspberry Pi Pico
**Test Date**: 2026-01-29

### Tests Performed
- [x] Motor control (forward/reverse)
- [x] Speed control (0-100%)
- [x] Emergency stop
- [x] Sensor reading accuracy
- [x] Power consumption measurement

### Results
- Motor response time: 50ms
- Speed accuracy: ±2%
- Power consumption: 150mA @ 5V
- All safety features working

### Issues Found
- None

### Photos/Videos
[Link to test video]
```

## Safety Guidelines

When contributing to safety-critical projects (rocketry, robotics):

1. **Always Include Safety Features**
   - Emergency stop mechanisms
   - Watchdog timers
   - Input validation
   - Bounds checking

2. **Test Thoroughly**
   - Bench test before field test
   - Use current-limited power supplies
   - Test in safe environment
   - Have emergency shutdown ready

3. **Document Hazards**
   - Identify potential dangers
   - Specify safety procedures
   - Include warnings in documentation

## Questions and Support

- **Issues**: Use GitHub Issues for bugs and feature requests
- **Discussions**: Use GitHub Discussions for questions
- **Email**: Contact maintainers for security issues

## Recognition

Contributors will be:
- Listed in CONTRIBUTORS.md
- Mentioned in release notes
- Credited in documentation

Thank you for contributing to making embedded systems more accessible!

---

**Last Updated**: 2026-01-29
**Version**: 1.0.0
