# FreeRTOS for Embedded Systems

This directory contains examples and code for FreeRTOS on Raspberry Pi Pico, Arduino, and other microcontrollers using TinyGo.

## Overview

FreeRTOS-style examples demonstrating real-time embedded systems for robotics, rocketry, and renewable energy monitoring. All examples use TinyGo for cross-platform compatibility across various microcontrollers.

## Directory Structure

### Basic Examples
- **blink_green.go**: Basic LED blinking with multi-tasking
- **blink_red_green_tune.go**: Multi-task LED control with buzzer melody

### [Robotics](robotics/) 🤖
Advanced robotics control systems:
- **motor_control.go**: DC motors, servos, and stepper motors
- **sensor_integration.go**: IMU, GPS, and ultrasonic sensors
- **robot_movement.go**: Autonomous navigation and movement patterns

**Features:**
- Differential drive robot control
- Obstacle avoidance algorithms
- Wall-following navigation
- Multi-sensor fusion
- Supports Raspberry Pi Pico and Arduino

### [Rocketry](rocketry/) 🚀
Rocket launch control and telemetry:
- **launch_control.go**: Launch sequence automation
- **telemetry_monitor.go**: Real-time flight data monitoring

**Features:**
- Multi-stage rocket support
- Automated countdown and ignition
- Stage separation control
- Parachute deployment
- Real-time telemetry (altitude, velocity, GPS)
- Radio transmission
- Data logging

⚠️ **Safety Warning**: Rocketry involves dangerous materials and high-speed flight. Always follow safety regulations and obtain necessary permits.

### [Energy](energy/) 🔋
Renewable energy monitoring systems:
- **wind_generator.go**: Wind turbine monitoring and control
- **solar_monitor.go**: Solar panel MPPT controller
- **hydro_monitor.go**: Hydroelectric turbine monitoring
- **thermo_generator.go**: Thermoelectric generator monitoring

**Features:**
- Real-time power monitoring (voltage, current, power)
- Automatic optimization (MPPT, yaw control)
- Safety protection (over-voltage, over-current)
- Environmental monitoring (wind speed, irradiance, temperature)
- Energy accumulation tracking
- Battery management

## Platform Support

### Raspberry Pi
- Raspberry Pi Pico (RP2040)
- Raspberry Pi Pico W (with WiFi)
- Raspberry Pi Zero

### Arduino
- Arduino Uno
- Arduino Nano
- Arduino Mega 2560
- Arduino Nano 33 IoT

### Other Microcontrollers
Any microcontroller supported by TinyGo can run these examples with pin configuration adjustments.

## Basic LED Examples

### Go Version
Refer to the [Go version](blink_green.go) for basic FreeRTOS-style multi-tasking.

### Advanced Multi-Task Example
See [blink_red_green_tune.go](blink_red_green_tune.go) for concurrent LED control and melody playback.

### MicroPython Version
MicroPython versions available for comparison with embedded Go.

## Building and Flashing

### Prerequisites
Install TinyGo:
```bash
# Linux
wget https://github.com/tinygo-org/tinygo/releases/download/v0.30.0/tinygo_0.30.0_amd64.deb
sudo dpkg -i tinygo_0.30.0_amd64.deb

# macOS
brew install tinygo

# Windows
# Download installer from https://github.com/tinygo-org/tinygo/releases
```

### Build Examples

**Raspberry Pi Pico:**
```bash
# Basic LED blink
tinygo flash -target=pico blink_green.go

# Robotics examples
cd robotics
tinygo flash -target=pico motor_control.go
tinygo flash -target=pico robot_movement.go

# Rocketry examples
cd ../rocketry
tinygo flash -target=pico launch_control.go

# Energy monitoring
cd ../energy
tinygo flash -target=pico wind_generator.go
tinygo flash -target=pico solar_monitor.go
```

**Arduino:**
```bash
# Basic LED blink
tinygo flash -target=arduino blink_green.go

# Motor control on Arduino
cd robotics
tinygo flash -target=arduino motor_control.go
```

## FreeRTOS Task Pattern

All examples follow FreeRTOS-style task patterns using goroutines:

```go
func Task1() {
    for {
        // Task execution
        time.Sleep(taskPeriod)
    }
}

func Task2() {
    for {
        // Task execution
        time.Sleep(taskPeriod)
    }
}

func main() {
    go Task1()  // Launch concurrent task
    go Task2()
    select {}   // Keep main running
}
```

## Key Features

### Real-Time Capabilities
- Concurrent task execution with goroutines
- Precise timing control
- Hardware interrupt support
- Priority-based scheduling

### Hardware Interfaces
- **GPIO**: Digital input/output
- **PWM**: Motor control, servo control
- **ADC**: Analog sensor reading
- **I2C**: Sensor communication (IMU, sensors)
- **SPI**: High-speed peripherals
- **UART**: Serial communication, GPS

### Safety and Reliability
- Watchdog timer support
- Error handling and recovery
- Over-voltage/current protection
- Emergency shutdown procedures
- Data validation and checksums

## Example Applications

### 1. Autonomous Robot
Combine robotics examples:
```go
go MotorControlTask()
go SensorReadTask()
go NavigationTask()
```

### 2. Model Rocket Flight Computer
Integrate rocketry systems:
```go
go TelemetryTask()
go LaunchControlTask()
go RecoveryTask()
```

### 3. Off-Grid Power System
Multi-source energy monitoring:
```go
go WindTurbineTask()
go SolarMonitorTask()
go BatteryManagerTask()
```

## Reference Documentation

### TinyGo and Embedded
- [TinyGo Documentation](https://tinygo.org/docs/)
- [TinyGo Machine Package](https://tinygo.org/docs/reference/machine/)
- [Raspberry Pi Pico Datasheet](https://datasheets.raspberrypi.com/pico/pico-datasheet.pdf)
- [Arduino Reference](https://www.arduino.cc/reference/en/)

### FreeRTOS
- [FreeRTOS Documentation](https://www.freertos.org/Documentation/RTOS_book.html)
- [FreeRTOS API Reference](https://www.freertos.org/a00106.html)

### Domain-Specific
- [Robotics](robotics/README.md) - Motor control, sensors, navigation
- [Rocketry](rocketry/README.md) - Launch control, telemetry
- [Energy](energy/README.md) - Renewable energy monitoring

### General
- [Go Documentation](https://golang.org/doc/)
- [Embedded C Documentation](https://www.embedded.com/embedded-c/)
- [MicroPython Documentation](https://docs.micropython.org/en/latest/)

## Performance Considerations

- **Memory Usage**: Typical 10-30KB RAM per application
- **Flash Usage**: 50-200KB depending on complexity
- **Update Rates**: 10-100 Hz typical
- **Real-Time**: Microsecond-level timing precision
- **Concurrency**: Multiple tasks run efficiently

## Debugging

### Serial Output
```go
machine.Serial.Configure(machine.UARTConfig{BaudRate: 115200})
machine.Serial.Write([]byte("Debug message\n"))
```

### Monitor Output
```bash
tinygo monitor
```

### GDB Debugging
```bash
tinygo gdb -target=pico program.go
```

## Contributing

When adding new FreeRTOS examples:
1. Follow the task pattern structure
2. Include comprehensive documentation
3. Add hardware requirements and wiring diagrams
4. Implement safety features
5. Test on actual hardware
6. Update relevant README files

## Safety Guidelines

⚠️ **Important Safety Notes:**

- **Robotics**: Use current-limited power supplies, implement emergency stops
- **Rocketry**: Follow all local regulations, never test near people or structures
- **Energy Systems**: Use proper fusing, implement over-voltage protection
- **General**: Always test code without hardware first, verify pin configurations

## License
This project is licensed under the MIT License - see the [LICENSE](../../../../LICENSE) file for details.
