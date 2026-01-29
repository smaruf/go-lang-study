# Embedded OS for Raspberry Pi and Arduino

This directory contains minimal embedded operating system examples and bare-metal programs for Raspberry Pi and Arduino using TinyGo.

## Overview

TinyGo is a Go compiler for small places. It can compile Go programs for microcontrollers, WebAssembly, and command-line tools. This makes it perfect for embedded systems development.

## Supported Boards

### Raspberry Pi
- Raspberry Pi Pico (RP2040)
- Raspberry Pi Pico W (RP2040 with WiFi)
- Raspberry Pi Zero
- Raspberry Pi 3/4 (bare metal)

### Arduino
- Arduino Uno (ATmega328p)
- Arduino Nano
- Arduino Mega 2560
- Arduino Nano 33 IoT

## Prerequisites

### Install TinyGo

**Linux:**
```bash
wget https://github.com/tinygo-org/tinygo/releases/download/v0.30.0/tinygo_0.30.0_amd64.deb
sudo dpkg -i tinygo_0.30.0_amd64.deb
```

**macOS:**
```bash
brew tap tinygo-org/tools
brew install tinygo
```

**Windows:**
Download and install from: https://github.com/tinygo-org/tinygo/releases

### Verify Installation
```bash
tinygo version
```

## Project Structure

```
embedded-os/
├── README.md                 # This file
├── raspberry-pi/            # Raspberry Pi examples
│   ├── blinky.go           # LED blink example
│   ├── button.go           # Button input example
│   ├── pwm_led.go          # PWM LED dimming
│   ├── temperature.go      # Temperature sensor
│   └── wifi_server.go      # WiFi web server (Pico W)
├── arduino/                 # Arduino examples
│   ├── blinky.go           # LED blink example
│   ├── button.go           # Button input
│   ├── serial.go           # Serial communication
│   ├── servo.go            # Servo motor control
│   └── ultrasonic.go       # Ultrasonic distance sensor
├── tiny/                    # TinyGo projects
│   ├── hello.go            # Hello world TinyGo
│   ├── pwm_blinky.go       # PWM blinky example
│   └── sr71sim/            # SR-71 flight simulator
└── freeRTOS/               # FreeRTOS examples
    ├── robotics/           # Robot control systems
    ├── rocketry/           # Rocket launch control
    └── energy/             # Energy monitoring systems
```

## Examples

### Raspberry Pi Pico

#### 1. Blinky (LED Blink)
```bash
cd raspberry-pi
tinygo flash -target=pico blinky.go
```

#### 2. Button Input
```bash
tinygo flash -target=pico button.go
```

#### 3. PWM LED
```bash
tinygo flash -target=pico pwm_led.go
```

#### 4. Temperature Sensor
```bash
tinygo flash -target=pico temperature.go
```

#### 5. WiFi Server (Pico W only)
```bash
tinygo flash -target=pico-w wifi_server.go
```

### Arduino Uno

#### 1. Blinky
```bash
cd arduino
tinygo flash -target=arduino blinky.go
```

#### 2. Button Input
```bash
tinygo flash -target=arduino button.go
```

#### 3. Serial Communication
```bash
tinygo flash -target=arduino serial.go
# Monitor serial output:
tinygo monitor
```

#### 4. Servo Control
```bash
tinygo flash -target=arduino servo.go
```

#### 5. Ultrasonic Sensor
```bash
tinygo flash -target=arduino ultrasonic.go
```

## Building Without Flashing

Generate hex file without flashing:
```bash
tinygo build -o output.hex -target=pico blinky.go
```

Generate UF2 file for Pico:
```bash
tinygo build -o output.uf2 -target=pico blinky.go
```

## Hardware Setup

### Raspberry Pi Pico - LED Blink
- LED: Connect to GPIO 25 (built-in LED)
- No external components needed for basic blink

### Raspberry Pi Pico - Button
- Button: GPIO 15
- Pull-up resistor: 10kΩ to 3.3V
- Button to GND

### Arduino Uno - LED Blink
- LED: Pin 13 (built-in LED)
- No external components needed

### Arduino Uno - Servo
- Servo signal: Pin 9
- Servo power: 5V
- Servo ground: GND

## Debugging

### Serial Monitor
```bash
# Flash and monitor
tinygo flash -target=pico -monitor blinky.go

# Or just monitor
tinygo monitor
```

### GDB Debugging (Advanced)
```bash
# Start GDB server
tinygo gdb -target=pico blinky.go

# In another terminal, connect with GDB
arm-none-eabi-gdb
```

## Key Concepts

### GPIO (General Purpose Input/Output)
- Digital input/output pins
- Control LEDs, buttons, sensors
- Configure as input or output

### PWM (Pulse Width Modulation)
- Control LED brightness
- Control servo motors
- Generate audio tones

### ADC (Analog to Digital Converter)
- Read analog sensors
- Temperature, light, potentiometers
- Convert voltage to digital value

### Interrupts
- Respond to events immediately
- Button presses, timer events
- Efficient event handling

### I2C and SPI
- Communication protocols
- Talk to sensors and displays
- I2C: 2-wire protocol (SDA, SCL)
- SPI: 4-wire protocol (MOSI, MISO, SCK, CS)

## Advanced Features

### Multi-tasking
TinyGo supports goroutines on microcontrollers:
```go
go func() {
    // Run in background
}()
```

### WiFi (Pico W)
Connect to WiFi and create web servers:
```go
import "machine"

func main() {
    // Configure WiFi
    // Create HTTP server
}
```

### Bluetooth (Some boards)
Some boards support Bluetooth LE:
```go
import "tinygo.org/x/bluetooth"
```

## Performance Tips

1. **Use const where possible** - Reduces memory usage
2. **Avoid allocations** - Minimize heap allocations
3. **Use buffered channels** - Better concurrency performance
4. **Disable GC if needed** - For real-time applications
5. **Optimize binary size** - Use `-opt=2` flag

## Troubleshooting

### Board not detected
```bash
# Linux: Check permissions
sudo usermod -a -G dialout $USER
# Logout and login again

# Check USB connection
lsusb
```

### Flash fails
```bash
# Put board in bootloader mode
# Pico: Hold BOOTSEL button while connecting
# Arduino: Press reset twice quickly
```

### Out of memory
```bash
# Try optimizing
tinygo flash -opt=2 -target=pico blinky.go

# Check memory usage
tinygo build -size=short -target=pico blinky.go
```

## Resources

- [TinyGo Documentation](https://tinygo.org/docs/)
- [TinyGo Supported Boards](https://tinygo.org/docs/reference/microcontrollers/)
- [Machine Package Reference](https://tinygo.org/docs/reference/machine/)
- [Raspberry Pi Pico Datasheet](https://datasheets.raspberrypi.com/pico/pico-datasheet.pdf)
- [Arduino Reference](https://www.arduino.cc/reference/en/)

## Examples from Other Projects

Check these directories for more embedded examples:
- `tiny/` - Basic TinyGo examples (SR-71 simulator)
- `freeRTOS/` - FreeRTOS examples (robotics, rocketry, energy)
- `../gobot/` - Robotics framework examples

## Contributing

When adding new examples:
1. Include hardware requirements in comments
2. Add wiring diagrams if needed
3. Test on actual hardware
4. Document expected behavior
5. Add to this README
