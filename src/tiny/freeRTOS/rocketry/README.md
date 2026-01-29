# FreeRTOS Rocketry Examples

This directory contains rocket launch control and telemetry systems for FreeRTOS running on embedded platforms using TinyGo.

## Overview

These examples demonstrate real-time rocket control systems including launch sequencing, stage separation, telemetry monitoring, and flight control suitable for model rockets and amateur rocketry.

## ⚠️ Safety Warning

**DANGER**: Rocketry involves high temperatures, pressurized gases, explosive materials, and high-speed flight. These examples are for educational purposes. Always:
- Follow local laws and regulations
- Obtain necessary permits and clearances
- Use proper safety equipment and procedures
- Never test near people, buildings, or airports
- Consult with experienced rocketry clubs (NAR, TRA)

## Files

### Launch Control (`launch_control.go`)
Complete launch sequence automation:
- **Pre-launch Checks**: System verification
- **Countdown Sequence**: Automated T-minus countdown
- **Ignition Control**: Precise ignition timing
- **Stage Separation**: Multi-stage rocket support
- **Parachute Deployment**: Recovery system activation
- **State Machine**: Flight phase management

**Hardware Requirements:**
- Ignition relay/solid-state relay (rated for igniter current)
- Separation mechanism actuators
- Parachute deployment servos or pyrotechnic charges
- Status LEDs
- Safety key switch

**Flight States:**
```
Prelaunch → Ignition → Powered Flight → 
Stage 1 Separation → Stage 2 Ignition → 
Coasting → Apogee → Parachute Deploy → 
Landing → Recovered
```

**Wiring:**
```
Ignition System:
- GPIO16 → Ignition relay (through safety key)

Stage Separation:
- GPIO17 → Stage 1 separation actuator
- GPIO18 → Stage 2 ignition/separation

Recovery:
- GPIO19 → Parachute deployment servo

Status LEDs:
- GPIO24 → Green (Ready)
- GPIO25 → Yellow (Armed)
- LED    → Red (Abort)
```

### Telemetry Monitor (`telemetry_monitor.go`)
Real-time flight data collection and transmission:
- **Altitude Tracking**: Barometric altitude
- **Velocity Calculation**: Speed and acceleration
- **GPS Position**: Latitude, longitude, altitude
- **Temperature Monitoring**: Avionics temperature
- **Battery Voltage**: Power system monitoring
- **Radio Transmission**: Ground station communication
- **Data Logging**: Onboard flash storage
- **Alert System**: Critical condition warnings

**Hardware Requirements:**
- Barometric pressure sensor (BMP280/MS5611)
- GPS module (NEO-6M or similar)
- Temperature sensor (TMP36)
- Radio transmitter (LoRa, 433MHz, or 915MHz)
- SD card for data logging
- Battery voltage divider

**Telemetry Data Packet:**
```
{
  Timestamp:    Mission time (ms)
  Altitude:     Meters above launch site
  Velocity:     m/s
  Acceleration: m/s²
  Latitude:     GPS coordinate
  Longitude:    GPS coordinate
  Temperature:  °C
  Pressure:     kPa
  BatteryVolt:  Volts
  State:        Flight phase
}
```

## Building and Flashing

### For Raspberry Pi Pico:
```bash
# Launch control system
tinygo flash -target=pico launch_control.go

# Telemetry monitoring
tinygo flash -target=pico telemetry_monitor.go
```

### For Flight Computer Integration:
```bash
# Combined launch control + telemetry
tinygo flash -target=pico -opt=2 flight_computer.go
```

## Configuration

### Launch Parameters
Modify constants for your rocket:
```go
const (
    STAGE1_BURN_TIME    = 5 * time.Second
    STAGE1_SEP_ALTITUDE = 10000.0  // meters
    STAGE2_BURN_TIME    = 3 * time.Second
    PARACHUTE_DEPLOY_DELAY = 2 * time.Second
)
```

### Telemetry Settings
```go
const (
    TELEMETRY_RATE = 10  // Hz
    RADIO_POWER    = 20  // dBm
    LOG_INTERVAL   = 100 // ms
)
```

## Flight Phases

### 1. Pre-Launch (-T to T-0)
- System checks
- Sensor calibration
- Radio link verification
- Safety arm/disarm
- Countdown sequence

### 2. Launch (T+0 to Burnout)
- Ignition
- Thrust phase
- Powered flight telemetry
- Guidance (if active)

### 3. Stage Separation
- Burnout detection
- Stage separation command
- Inter-stage delay
- Second stage ignition

### 4. Coast Phase
- Ballistic flight
- Trajectory tracking
- Apogee prediction

### 5. Recovery
- Apogee detection
- Parachute deployment
- Descent tracking
- Landing detection

## Safety Features

### Hardware Safety
```go
// Ignition safety key (hardware switch)
// Prevents accidental ignition
if !safetyKeyEngaged {
    return // Cannot ignite
}

// Continuity check before launch
if !igniterContinuity {
    abort("Igniter not connected")
}
```

### Software Safety
```go
// Abort conditions
func SafetyCheck() {
    if batteryVoltage < MIN_VOLTAGE {
        Abort("Low battery")
    }
    if accelExceeds > MAX_G_FORCE {
        Abort("High G-force")
    }
    if altitude > MAX_ALTITUDE {
        DeployParachute()
    }
}
```

## Telemetry Output

### Ground Station Display
```
T+0015.234s | ALT:  1524m | VEL:  143m/s | ACC:  28.4m/s²
            | LAT: 28.5729 | LON: -80.6490
            | TEMP: 18.2°C | BATT: 11.8V
            | STATE: POWERED_FLIGHT_STAGE1
```

### Data Logging Format (CSV)
```csv
timestamp,altitude,velocity,acceleration,latitude,longitude,temperature,pressure,battery,state
15234,1524.3,143.2,28.4,28.5729,-80.6490,18.2,85.3,11.8,POWERED_FLIGHT
```

## Example Missions

### 1. Single-Stage Model Rocket
```go
// Simple altitude mission
// Launch → Coast → Apogee → Parachute → Land
```

### 2. Two-Stage Rocket
```go
// Extended altitude mission
// Launch → Stage 1 → Separation → Stage 2 → Coast → Parachute
```

### 3. Guided Rocket
```go
// Add thrust vector control or canards
// Implement trajectory guidance
```

## Performance Specifications

- **Update Rate**: 10-100 Hz (10ms-100ms cycles)
- **Telemetry Rate**: 10 Hz (100ms packets)
- **Radio Range**: 10-50 km (depending on radio module)
- **Altitude Range**: 0-30 km
- **Max Acceleration**: 50g
- **Memory Usage**: ~30KB flash, ~15KB RAM

## Testing Procedures

### Ground Testing
1. **Static Tests**: Test without ignition
2. **Continuity Tests**: Verify igniter connections
3. **Radio Tests**: Confirm telemetry reception
4. **Sensor Calibration**: Zero altitude, check GPS lock
5. **Recovery Tests**: Test parachute deployment

### Flight Testing
1. **Low-Altitude Test**: Sub-500m flight
2. **Verify all systems**: Check data logging
3. **Recovery Test**: Confirm parachute deployment
4. **Incremental Testing**: Gradually increase altitude

## Regulations and Compliance

### United States (FAA)
- Rockets >453g require FAA waiver
- Altitude restrictions apply
- No flight in controlled airspace

### Amateur Rocketry Organizations
- [NAR](http://www.nar.org/) - National Association of Rocketry
- [TRA](http://www.tripoli.org/) - Tripoli Rocketry Association

## Resources

- [OpenRocket](https://openrocket.info/) - Rocket design software
- [RocketPy](https://github.com/RocketPy-Team/RocketPy) - Trajectory simulation
- [Apogee Rockets](https://www.apogeerockets.com/) - Components and education
- [NAKKA-ROCKETRY](http://www.nakka-rocketry.net/) - Amateur rocketry resource

## Debugging

### Serial Output
```go
// Enable verbose logging
machine.Serial.Printf("T+%.3f ALT:%0.1f VEL:%.1f\n", 
    missionTime, altitude, velocity)
```

### Ground Station Software
Use a laptop or Raspberry Pi with radio receiver to monitor telemetry in real-time.

## Contributing

When adding rocketry examples:
1. Emphasize safety in all documentation
2. Include abort/emergency procedures
3. Test thoroughly on ground before flight
4. Follow rocketry best practices
5. Document all hardware interfaces

## License

MIT License - See repository root for details

**DISCLAIMER**: The authors assume no liability for any damage or injury resulting from use of this code. Rocketry is inherently dangerous. Use at your own risk.
