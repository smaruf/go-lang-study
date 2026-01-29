# FreeRTOS Robotics Examples

This directory contains robotics control examples for FreeRTOS running on Raspberry Pi Pico and Arduino platforms using TinyGo.

## Overview

These examples demonstrate real-time robotics applications including motor control, sensor integration, and autonomous navigation patterns suitable for embedded platforms.

## Files

### Motor Control (`motor_control.go`)
Comprehensive motor control system supporting:
- **DC Motors**: PWM speed control with direction
- **Servo Motors**: Precise angle positioning (0-180°)
- **Stepper Motors**: Step-by-step positioning with half-step sequence

**Hardware Requirements:**
- DC Motor with H-bridge (L298N or similar)
- Servo motor (SG90 or similar)
- Stepper motor (28BYJ-48 or NEMA17)
- Motor driver boards

**Wiring:**
```
DC Motor:
- GPIO2 → PWM to motor driver
- GPIO3 → Direction to motor driver

Servo Motor:
- GPIO4 → Servo signal pin

Stepper Motor:
- GPIO5-8 → Stepper driver inputs
```

### Sensor Integration (`sensor_integration.go`)
Multi-sensor integration for robotics:
- **Ultrasonic Sensors**: Distance measurement (HC-SR04)
- **IMU**: 6-axis accelerometer/gyroscope (MPU6050)
- **GPS**: Position and navigation data

**Hardware Requirements:**
- HC-SR04 ultrasonic sensor
- MPU6050 IMU module (I2C)
- GPS module (UART)

**Wiring:**
```
Ultrasonic Sensor:
- GPIO10 → TRIG
- GPIO11 → ECHO

IMU (I2C):
- GPIO12 → SDA
- GPIO13 → SCL

GPS (UART):
- GPIO14 → TX
- GPIO15 → RX
```

### Robot Movement (`robot_movement.go`)
Autonomous robot navigation patterns:
- **Differential Drive**: Two-wheel robot control
- **Obstacle Avoidance**: Dynamic path planning
- **Wall Following**: Edge-following behavior
- **Pattern Movement**: Square, circle, and figure-8 patterns

**Features:**
- Real-time obstacle detection and avoidance
- Proportional control for wall following
- Pre-programmed movement patterns
- Multi-sensor fusion for navigation

## Building and Flashing

### For Raspberry Pi Pico:
```bash
# Motor control example
tinygo flash -target=pico motor_control.go

# Sensor integration
tinygo flash -target=pico sensor_integration.go

# Robot movement patterns
tinygo flash -target=pico robot_movement.go
```

### For Arduino:
```bash
# Motor control example
tinygo flash -target=arduino motor_control.go

# Sensor integration
tinygo flash -target=arduino-nano sensor_integration.go
```

## Configuration

Each example uses GPIO pin constants that can be modified for your hardware setup:

```go
const (
    MOTOR_PWM_PIN = machine.GPIO2  // Adjust as needed
    MOTOR_DIR_PIN = machine.GPIO3
    // ... other pins
)
```

## FreeRTOS Task Structure

All examples follow FreeRTOS task patterns:
```go
func RobotTask() {
    // Initialization
    
    for {
        // Periodic task execution
        // Sensor reading
        // Control logic
        // Actuator commands
        
        time.Sleep(taskPeriod)
    }
}

func main() {
    go RobotTask()  // Launch as FreeRTOS-style task
    select {}       // Keep main running
}
```

## Application Examples

### 1. Line-Following Robot
Combine sensor integration with differential drive:
```go
// Use ultrasonic or IR sensors
// Adjust motor speeds based on sensor readings
```

### 2. Obstacle-Avoiding Robot
Use the navigation controller with multiple sensors:
```go
nav := NewNavigationController(robot)
nav.AvoidObstacles()
```

### 3. Autonomous Maze Solver
Implement wall-following algorithm:
```go
nav.FollowWall(1) // Follow right wall
```

## Performance Notes

- **Update Frequency**: 10-100 Hz typical for robot control
- **Sensor Latency**: Ultrasonic ~25ms, IMU ~5ms, GPS ~1Hz
- **Motor Response**: PWM frequency 1-20 kHz
- **Memory Usage**: ~10-20KB RAM per task

## Safety Considerations

1. **Motor Current Limits**: Use appropriate motor drivers with current limiting
2. **Battery Protection**: Monitor voltage to prevent over-discharge
3. **Emergency Stop**: Implement hardware emergency stop button
4. **Obstacle Detection**: Always include failsafe for obstacle detection
5. **Sensor Validation**: Validate sensor readings before acting

## Advanced Features

### Multi-Tasking
Run multiple control loops concurrently:
```go
go MotorControlTask()
go SensorReadTask()
go NavigationTask()
```

### Priority Management
Critical tasks (safety) should run at higher frequency:
```go
// Safety check: 100 Hz
// Navigation: 10 Hz
// Telemetry: 1 Hz
```

## Debugging

Enable serial output for debugging:
```go
machine.Serial.Configure(machine.UARTConfig{BaudRate: 115200})
machine.Serial.Write([]byte("Debug message\n"))
```

Monitor with:
```bash
tinygo monitor
```

## Resources

- [TinyGo Machine Package](https://tinygo.org/docs/reference/machine/)
- [Raspberry Pi Pico Pinout](https://datasheets.raspberrypi.com/pico/Pico-R3-A4-Pinout.pdf)
- [Arduino Pin Reference](https://www.arduino.cc/reference/en/)
- [MPU6050 Datasheet](https://invensense.tdk.com/products/motion-tracking/6-axis/mpu-6050/)

## Contributing

When adding new robotics examples:
1. Follow the existing task structure
2. Include hardware requirements and wiring diagrams
3. Add safety checks and error handling
4. Test on actual hardware before committing
5. Document pin configurations clearly

## License

MIT License - See repository root for details
