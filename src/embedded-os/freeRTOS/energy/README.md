# FreeRTOS Renewable Energy Monitoring Systems

This directory contains monitoring and control systems for renewable energy generators running on FreeRTOS using TinyGo for Raspberry Pi Pico and Arduino platforms.

## Overview

These examples demonstrate real-time monitoring and control of various renewable energy systems including wind turbines, solar panels, hydroelectric generators, and thermoelectric generators. Each system includes power optimization, safety monitoring, and data logging.

## Files

### Wind Generator (`wind_generator.go`)
Complete wind turbine monitoring and control system:
- **Wind Speed Measurement**: Anemometer pulse counting
- **Wind Direction**: Potentiometer-based wind vane
- **Power Monitoring**: Voltage, current, and power calculation
- **Yaw Control**: Automatic turbine orientation
- **Brake Control**: Emergency shutdown and speed limiting
- **Safety Shutdown**: Over-speed and over-voltage protection

**Hardware Requirements:**
- Anemometer (pulse output)
- Wind direction sensor (potentiometer)
- Voltage sensor (0-100V with voltage divider)
- Current sensor (Hall effect, 0-20A)
- Yaw motor with direction control
- Brake actuator (solenoid or electromagnetic brake)

**Key Features:**
- Automatic yaw adjustment to face wind
- Emergency brake on high winds (>25 m/s)
- Power output optimization
- Cut-in wind speed: 3 m/s
- Shutdown wind speed: 25 m/s

### Solar Monitor (`solar_monitor.go`)
Solar panel system with MPPT (Maximum Power Point Tracking):
- **Panel Monitoring**: Voltage, current, power measurement
- **Battery Management**: Charge state monitoring
- **MPPT Algorithm**: Perturb and Observe method
- **Environmental Sensing**: Irradiance and temperature
- **Load Management**: Automatic load switching
- **Charge States**: Low battery, bulk charge, float charge

**Hardware Requirements:**
- Solar panel voltage sensor (0-50V)
- Panel current sensor (0-10A)
- Battery voltage monitor (12V/24V system)
- Light sensor for irradiance
- Temperature sensor (TMP36)
- PWM-controlled DC-DC converter for MPPT
- Load relay

**Charge States:**
- **Low Battery**: <11.5V - Loads disconnected
- **Bulk Charge**: 11.5V-14.4V - Maximum charging
- **Float Charge**: 14.4V-14.7V - Maintenance charge
- **Over-Voltage**: >14.7V - Charging stopped

### Hydro Monitor (`hydro_monitor.go`)
Hydroelectric turbine monitoring system:
- **Flow Measurement**: Water flow sensor (L/min)
- **Pressure Monitoring**: Hydraulic pressure (kPa)
- **Power Generation**: Voltage and current monitoring
- **Gate Control**: Water gate servo positioning
- **Bypass Valve**: Emergency water bypass
- **Optimization**: Automatic power output optimization

**Hardware Requirements:**
- Flow sensor (YF-S201 or similar)
- Pressure sensor (0-1000 kPa)
- Voltage and current sensors
- Gate control servo
- Bypass valve actuator
- RPM sensor (optional)

**Safety Features:**
- Over-speed shutdown (>3000 RPM)
- Over-voltage protection (>60V)
- Over-current protection (>15A)
- Low-flow cavitation prevention
- Emergency bypass activation

### Thermo Generator (`thermo_generator.go`)
Thermoelectric generator (TEG) monitoring:
- **Temperature Differential**: Hot and cold side monitoring
- **Power Optimization**: Maximize efficiency
- **Cooling Control**: Fan and pump management
- **Heat Source Control**: Optional heater management
- **Efficiency Calculation**: Real-time efficiency monitoring

**Hardware Requirements:**
- K-type thermocouples (hot and cold sides)
- TEG voltage and current sensors
- Cooling fan with PWM control
- Water pump (optional)
- Heater control (for controlled heat source)
- Thermocouple amplifier

**Operating Principles:**
- Seebeck effect: Temperature differential → Electricity
- Carnot efficiency limit
- Optimal cold side: 20-30°C
- Target temperature differential: 150°C
- Typical efficiency: 5-8%

## Common Features

All energy monitoring systems include:
- **Real-time Monitoring**: Continuous sensor reading
- **Power Calculation**: Voltage × Current
- **Energy Accumulation**: Total kWh generated
- **Safety Checks**: Over-voltage, over-current protection
- **Status Indicators**: LED indicators for system state
- **Data Logging**: Periodic data logging capability

## Building and Flashing

### For Raspberry Pi Pico:
```bash
# Wind turbine monitor
tinygo flash -target=pico wind_generator.go

# Solar MPPT controller
tinygo flash -target=pico solar_monitor.go

# Hydro turbine monitor
tinygo flash -target=pico hydro_monitor.go

# Thermoelectric monitor
tinygo flash -target=pico thermo_generator.go
```

### For Arduino:
```bash
tinygo flash -target=arduino wind_generator.go
tinygo flash -target=arduino-nano solar_monitor.go
```

## Pin Configurations

### Wind Generator
```
Sensors:
- GPIO16: Anemometer pulse input
- ADC0: Wind direction (0-360°)
- ADC1: Voltage (0-100V)
- ADC2: Current (0-20A)

Controls:
- GPIO17: Brake control
- GPIO18: Yaw motor PWM
- GPIO19: Yaw motor direction

Status:
- GPIO20: Power generation LED
- GPIO21: Warning LED
- LED: Shutdown LED
```

### Solar Monitor
```
Sensors:
- ADC0: Panel voltage
- ADC1: Panel current
- ADC2: Battery voltage
- ADC3: Light sensor
- ADC4: Temperature sensor

Controls:
- GPIO16: MPPT PWM control
- GPIO17: Load relay

Status:
- GPIO18: Charging LED
- GPIO19: Full charge LED
- LED: Low battery LED
```

## Performance Specifications

### Wind Generator
- **Wind Speed Range**: 0-40 m/s
- **Power Range**: 0-500W (typical small turbine)
- **Voltage Range**: 12V/24V/48V systems
- **Update Rate**: 10 Hz
- **Yaw Adjustment**: ±5° steps

### Solar Monitor
- **Panel Voltage**: Up to 50V
- **Panel Current**: Up to 10A
- **MPPT Efficiency**: 95-98%
- **Battery Types**: Lead-acid, Li-ion (with voltage adjustments)
- **Update Rate**: 10 Hz

### Hydro Monitor
- **Flow Range**: 0-200 L/min
- **Pressure Range**: 0-1000 kPa
- **Power Range**: 0-1000W
- **RPM Range**: 0-3000 RPM
- **Update Rate**: 10 Hz

### Thermo Generator
- **Temperature Range**: 0-500°C (hot side)
- **Temperature Differential**: 50-350°C
- **Power Range**: 0-100W (typical TEG)
- **Efficiency**: 5-8% typical
- **Update Rate**: 10 Hz

## Safety Features

### All Systems
```go
// Voltage protection
if voltage > MAX_VOLTAGE {
    DisconnectLoad()
    TriggerAlarm()
}

// Current protection
if current > MAX_CURRENT {
    ReducePower()
    TriggerAlarm()
}

// Temperature protection (where applicable)
if temperature > MAX_TEMP {
    EmergencyShutdown()
}
```

### Wind Turbine Specific
- Furling or braking in high winds (>25 m/s)
- Tower vibration monitoring (if sensor available)
- Lightning protection recommendations

### Solar Specific
- Reverse current protection
- Battery over-charge protection
- Deep discharge protection

### Hydro Specific
- Cavitation prevention
- Trash rack blockage detection
- Emergency shutdown valve

## Data Logging

All systems support data logging:

```go
type LogEntry struct {
    Timestamp   int64
    Voltage     float32
    Current     float32
    Power       float32
    TotalEnergy float32
    // System-specific fields
}

// Log to SD card (requires SD card module)
// Log to flash memory
// Transmit via radio/WiFi
```

## Monitoring Dashboard

Create a ground station or web dashboard to monitor:
- Real-time power generation
- Energy production graphs
- System efficiency
- Environmental conditions
- Alert notifications

## Example Multi-System Setup

Combine multiple renewable sources:

```go
func main() {
    // Run multiple energy monitors
    go WindTurbineTask()
    go SolarMonitorTask()
    go HydroMonitorTask()
    go ThermoMonitorTask()
    
    // Aggregate power data
    go PowerAggregatorTask()
    
    // Battery management
    go BatteryManagerTask()
    
    select {}
}
```

## Efficiency Optimization

### Wind Turbine
- Blade pitch control (advanced)
- MPPT for optimal load matching
- Yaw alignment with wind direction

### Solar Panel
- Perturb and Observe MPPT
- Temperature coefficient compensation
- Panel cleaning schedule based on efficiency drop

### Hydro Turbine
- Gate position optimization
- Minimize cavitation
- Maintain optimal RPM

### Thermoelectric
- Maximize temperature differential
- Optimize cold-side cooling
- Minimize thermal resistance

## Integration Examples

### Off-Grid Home System
```
Solar (primary) + Wind (supplementary) + Hydro (baseload)
→ Battery Bank → Inverter → Home Loads
```

### Remote Monitoring Station
```
Solar + Thermoelectric (backup) → Battery → Sensors + Radio
```

### Educational Laboratory
```
All four systems → Data acquisition → Real-time display
```

## Resources

### Wind Energy
- [AWEA](https://www.awea.org/) - American Wind Energy Association
- [Wind Turbine Design](https://www.wind-power-program.com/)

### Solar Energy
- [NREL](https://www.nrel.gov/) - National Renewable Energy Laboratory
- [PVEducation](https://www.pveducation.org/) - Solar cell resources

### Hydro Energy
- [Micro-Hydro Power](https://www.microhydropower.net/)
- [Small Hydro](https://www.smallhydro.com/)

### Thermoelectric
- [TEG Resources](https://www.tegpower.com/)
- [Seebeck Effect](https://en.wikipedia.org/wiki/Seebeck_effect)

## Testing and Calibration

### Wind Turbine
1. Anemometer calibration curve
2. Yaw mechanism alignment
3. Brake system verification

### Solar Panel
1. Open-circuit voltage test
2. Short-circuit current test
3. MPPT algorithm tuning

### Hydro Turbine
1. Flow sensor calibration
2. Pressure sensor zeroing
3. Gate position calibration

### Thermoelectric
1. Thermocouple linearization
2. Cooling system optimization
3. Heat source characterization

## Contributing

When adding energy monitoring examples:
1. Include detailed sensor specifications
2. Provide calibration procedures
3. Document safety considerations
4. Add efficiency calculations
5. Test with actual hardware

## License

MIT License - See repository root for details
