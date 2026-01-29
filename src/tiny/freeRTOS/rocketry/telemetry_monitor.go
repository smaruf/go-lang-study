package main

import (
	"machine"
	"time"
)

// Rocket telemetry monitoring and data logging system
// Tracks altitude, velocity, acceleration, temperature, and GPS

const (
	// Telemetry transmission pins
	RADIO_TX = machine.GPIO26
	RADIO_RX = machine.GPIO27
	
	// Data logging LED
	LOG_LED = machine.GPIO28
)

// TelemetryData represents a complete telemetry packet
type TelemetryData struct {
	Timestamp    int64
	Altitude     float32
	Velocity     float32
	Acceleration float32
	Latitude     float32
	Longitude    float32
	Temperature  float32
	Pressure     float32
	BatteryVolt  float32
	State        string
}

// TelemetrySystem manages telemetry collection and transmission
type TelemetrySystem struct {
	data         TelemetryData
	radioTx      machine.Pin
	logLED       machine.Pin
	packetCount  int
	startTime    time.Time
}

// NewTelemetrySystem creates a new telemetry system
func NewTelemetrySystem() *TelemetrySystem {
	ts := &TelemetrySystem{
		startTime:   time.Now(),
		packetCount: 0,
	}
	
	// Configure pins
	ts.radioTx = RADIO_TX
	ts.logLED = LOG_LED
	
	ts.radioTx.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ts.logLED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	
	return ts
}

// CollectData collects sensor data
func (ts *TelemetrySystem) CollectData() {
	ts.data.Timestamp = time.Since(ts.startTime).Milliseconds()
	
	// Simulate sensor readings
	ts.data.Altitude = 15000.0 + float32(ts.packetCount)*10.0
	ts.data.Velocity = 200.0 + float32(ts.packetCount)*2.0
	ts.data.Acceleration = 25.0
	ts.data.Latitude = 28.5729 + float32(ts.packetCount)*0.0001
	ts.data.Longitude = -80.6490 + float32(ts.packetCount)*0.0001
	ts.data.Temperature = 20.0 - float32(ts.packetCount)*0.5
	ts.data.Pressure = 101.3 - float32(ts.packetCount)*0.1
	ts.data.BatteryVolt = 12.0 - float32(ts.packetCount)*0.01
	ts.data.State = "POWERED_FLIGHT"
	
	ts.packetCount++
}

// TransmitData transmits telemetry data via radio
func (ts *TelemetrySystem) TransmitData() {
	// Indicate transmission with LED
	ts.logLED.High()
	
	// Simulate radio transmission
	// In real implementation, encode data and send via UART/SPI to radio module
	ts.radioTx.High()
	time.Sleep(50 * time.Millisecond)
	ts.radioTx.Low()
	
	ts.logLED.Low()
}

// LogData logs telemetry data to onboard storage
func (ts *TelemetrySystem) LogData() {
	// In real implementation, write to SD card or flash memory
	// For now, just blink LED to indicate logging
	ts.logLED.High()
	time.Sleep(10 * time.Millisecond)
	ts.logLED.Low()
}

// CheckAlerts checks for critical conditions
func (ts *TelemetrySystem) CheckAlerts() []string {
	alerts := []string{}
	
	if ts.data.BatteryVolt < 10.0 {
		alerts = append(alerts, "LOW_BATTERY")
	}
	
	if ts.data.Temperature > 50.0 || ts.data.Temperature < -20.0 {
		alerts = append(alerts, "TEMP_OUT_OF_RANGE")
	}
	
	if ts.data.Altitude > 30000.0 {
		alerts = append(alerts, "MAX_ALTITUDE_WARNING")
	}
	
	if ts.data.Acceleration > 50.0 {
		alerts = append(alerts, "HIGH_G_FORCE")
	}
	
	return alerts
}

// FlightComputer represents the main flight computer
type FlightComputer struct {
	telemetry   *TelemetrySystem
	flightState string
	missionTime time.Duration
}

// NewFlightComputer creates a new flight computer
func NewFlightComputer() *FlightComputer {
	return &FlightComputer{
		telemetry:   NewTelemetrySystem(),
		flightState: "PRELAUNCH",
		missionTime: 0,
	}
}

// UpdateFlightState updates the current flight state
func (fc *FlightComputer) UpdateFlightState() {
	// State machine based on telemetry
	altitude := fc.telemetry.data.Altitude
	velocity := fc.telemetry.data.Velocity
	
	if altitude < 100 && velocity < 10 {
		fc.flightState = "PRELAUNCH"
	} else if altitude < 10000 && velocity > 50 {
		fc.flightState = "POWERED_FLIGHT_STAGE1"
	} else if altitude >= 10000 && altitude < 20000 {
		fc.flightState = "STAGE_SEPARATION"
	} else if altitude >= 20000 && velocity > 100 {
		fc.flightState = "POWERED_FLIGHT_STAGE2"
	} else if velocity > 0 && velocity < 50 {
		fc.flightState = "COASTING"
	} else if velocity <= 0 {
		fc.flightState = "DESCENT"
	}
	
	fc.telemetry.data.State = fc.flightState
}

// PerformGuidance performs guidance computations
func (fc *FlightComputer) PerformGuidance() {
	// In real implementation, calculate trajectory corrections
	// and send commands to thrust vector control or RCS
	switch fc.flightState {
	case "POWERED_FLIGHT_STAGE1", "POWERED_FLIGHT_STAGE2":
		// Active guidance during powered flight
		// Calculate pitch and yaw corrections
		
	case "COASTING":
		// Passive guidance during coast
		// Monitor trajectory
		
	case "DESCENT":
		// Descent guidance
		// Monitor landing site
	}
}

// TelemetryMonitorTask is the FreeRTOS task for telemetry
func TelemetryMonitorTask() {
	computer := NewFlightComputer()
	
	statusLED := machine.LED
	statusLED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	
	ticker := time.NewTicker(100 * time.Millisecond) // 10 Hz telemetry rate
	defer ticker.Stop()
	
	for range ticker.C {
		// Collect sensor data
		computer.telemetry.CollectData()
		
		// Update flight state
		computer.UpdateFlightState()
		
		// Perform guidance computations
		computer.PerformGuidance()
		
		// Check for alerts
		alerts := computer.telemetry.CheckAlerts()
		if len(alerts) > 0 {
			// Flash status LED rapidly for alerts
			for i := 0; i < 5; i++ {
				statusLED.High()
				time.Sleep(50 * time.Millisecond)
				statusLED.Low()
				time.Sleep(50 * time.Millisecond)
			}
		}
		
		// Transmit telemetry data
		computer.telemetry.TransmitData()
		
		// Log data to storage
		computer.telemetry.LogData()
		
		// Heartbeat blink
		statusLED.High()
		time.Sleep(10 * time.Millisecond)
		statusLED.Low()
		
		// Update mission time
		computer.missionTime += 100 * time.Millisecond
	}
}

func main() {
	// Run telemetry monitoring task
	go TelemetryMonitorTask()
	
	// Keep main running
	select {}
}
