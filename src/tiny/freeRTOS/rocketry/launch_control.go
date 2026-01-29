package main

import (
	"machine"
	"time"
)

// Rocketry launch sequence and control system for FreeRTOS
// Simulates rocket launch control, telemetry, and safety systems

const (
	// Ignition and staging pins
	IGNITION_PIN   = machine.GPIO16
	STAGE1_SEP_PIN = machine.GPIO17
	STAGE2_SEP_PIN = machine.GPIO18
	PARACHUTE_PIN  = machine.GPIO19
	
	// Sensor pins
	ALTITUDE_SENSOR = machine.GPIO20
	PRESSURE_SENSOR = machine.GPIO21
	ACCEL_SDA       = machine.GPIO22
	ACCEL_SCL       = machine.GPIO23
	
	// Status LED pins
	STATUS_LED_GREEN  = machine.GPIO24
	STATUS_LED_YELLOW = machine.GPIO25
	STATUS_LED_RED    = machine.LED
)

// LaunchState represents the current state of the rocket
type LaunchState int

const (
	StatePrelaunch LaunchState = iota
	StateIgnition
	StatePoweredFlight
	StateStage1Separation
	StateStage2Ignition
	StateCoasting
	StateApogee
	StateParachuteDeploy
	StateLanding
	StateRecovered
)

// RocketController manages the launch sequence
type RocketController struct {
	state          LaunchState
	altitude       float32
	velocity       float32
	acceleration   float32
	launchTime     time.Time
	maxAltitude    float32
	ignitionPin    machine.Pin
	stage1SepPin   machine.Pin
	stage2SepPin   machine.Pin
	parachutePin   machine.Pin
	statusGreen    machine.Pin
	statusYellow   machine.Pin
	statusRed      machine.Pin
}

// NewRocketController creates a new rocket controller
func NewRocketController() *RocketController {
	rc := &RocketController{
		state:        StatePrelaunch,
		altitude:     0,
		velocity:     0,
		acceleration: 0,
		maxAltitude:  0,
	}
	
	// Configure ignition and separation pins
	rc.ignitionPin = IGNITION_PIN
	rc.stage1SepPin = STAGE1_SEP_PIN
	rc.stage2SepPin = STAGE2_SEP_PIN
	rc.parachutePin = PARACHUTE_PIN
	
	rc.ignitionPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	rc.stage1SepPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	rc.stage2SepPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	rc.parachutePin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	
	// Configure status LEDs
	rc.statusGreen = STATUS_LED_GREEN
	rc.statusYellow = STATUS_LED_YELLOW
	rc.statusRed = STATUS_LED_RED
	
	rc.statusGreen.Configure(machine.PinConfig{Mode: machine.PinOutput})
	rc.statusYellow.Configure(machine.PinConfig{Mode: machine.PinOutput})
	rc.statusRed.Configure(machine.PinConfig{Mode: machine.PinOutput})
	
	return rc
}

// SetStatus sets the status LED based on state
func (rc *RocketController) SetStatus(state string) {
	rc.statusGreen.Low()
	rc.statusYellow.Low()
	rc.statusRed.Low()
	
	switch state {
	case "ready":
		rc.statusGreen.High()
	case "armed":
		rc.statusYellow.High()
	case "launch":
		rc.statusGreen.High()
		rc.statusYellow.High()
	case "abort":
		rc.statusRed.High()
	}
}

// Ignite triggers the rocket ignition
func (rc *RocketController) Ignite() {
	rc.SetStatus("launch")
	rc.ignitionPin.High()
	rc.launchTime = time.Now()
	rc.state = StateIgnition
	
	// Ignition pulse duration
	time.Sleep(2 * time.Second)
	rc.ignitionPin.Low()
	
	rc.state = StatePoweredFlight
}

// SeparateStage1 separates the first stage
func (rc *RocketController) SeparateStage1() {
	rc.state = StateStage1Separation
	rc.stage1SepPin.High()
	time.Sleep(500 * time.Millisecond)
	rc.stage1SepPin.Low()
	
	// Wait for separation
	time.Sleep(2 * time.Second)
}

// IgniteStage2 ignites the second stage
func (rc *RocketController) IgniteStage2() {
	rc.state = StateStage2Ignition
	rc.stage2SepPin.High() // Reusing pin for stage 2 ignition
	time.Sleep(3 * time.Second)
	rc.stage2SepPin.Low()
}

// DeployParachute deploys the recovery parachute
func (rc *RocketController) DeployParachute() {
	rc.state = StateParachuteDeploy
	rc.parachutePin.High()
	time.Sleep(1 * time.Second)
	rc.parachutePin.Low()
}

// UpdateTelemetry updates rocket telemetry data
func (rc *RocketController) UpdateTelemetry() {
	// Simulate sensor readings
	flightTime := time.Since(rc.launchTime).Seconds()
	
	switch rc.state {
	case StatePoweredFlight, StateStage2Ignition:
		// Simulate powered flight acceleration
		rc.acceleration = 30.0 // m/s²
		rc.velocity += rc.acceleration * 0.1
		rc.altitude += rc.velocity * 0.1
		
	case StateCoasting:
		// Simulate coasting with deceleration
		rc.acceleration = -9.81 // Gravity
		rc.velocity += rc.acceleration * 0.1
		if rc.velocity < 0 {
			rc.velocity = 0
		}
		rc.altitude += rc.velocity * 0.1
		
	case StateParachuteDeploy, StateLanding:
		// Simulate descent with parachute
		rc.acceleration = -9.81
		rc.velocity = -5.0 // Descent rate with parachute
		rc.altitude += rc.velocity * 0.1
		if rc.altitude < 0 {
			rc.altitude = 0
			rc.velocity = 0
			rc.state = StateRecovered
		}
	}
	
	// Track maximum altitude
	if rc.altitude > rc.maxAltitude {
		rc.maxAltitude = rc.altitude
	}
	
	// Log telemetry (in real system, transmit via radio)
	_ = flightTime
}

// CheckApogee checks if rocket has reached apogee
func (rc *RocketController) CheckApogee() bool {
	return rc.velocity <= 0 && rc.state == StateCoasting
}

// LaunchSequenceTask is the main FreeRTOS task for launch control
func LaunchSequenceTask() {
	rocket := NewRocketController()
	
	// Pre-launch checks
	rocket.SetStatus("ready")
	time.Sleep(3 * time.Second)
	
	// Arm the system
	rocket.SetStatus("armed")
	time.Sleep(2 * time.Second)
	
	// T-10 seconds countdown
	for i := 10; i > 0; i-- {
		// Blink LED for countdown
		rocket.statusYellow.High()
		time.Sleep(500 * time.Millisecond)
		rocket.statusYellow.Low()
		time.Sleep(500 * time.Millisecond)
	}
	
	// IGNITION!
	rocket.Ignite()
	
	// Main flight loop
	for {
		rocket.UpdateTelemetry()
		
		switch rocket.state {
		case StatePoweredFlight:
			// Check for stage 1 separation (at 10km altitude)
			if rocket.altitude > 10000 {
				rocket.SeparateStage1()
				time.Sleep(1 * time.Second)
				rocket.IgniteStage2()
			}
			
		case StateStage2Ignition:
			// Burn second stage for 3 seconds
			if time.Since(rocket.launchTime) > 13*time.Second {
				rocket.state = StateCoasting
			}
			
		case StateCoasting:
			// Check for apogee
			if rocket.CheckApogee() {
				rocket.state = StateApogee
				time.Sleep(500 * time.Millisecond)
				rocket.DeployParachute()
			}
			
		case StateRecovered:
			// Mission complete
			rocket.SetStatus("ready")
			return
		}
		
		// Update at 10 Hz
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	// Run launch sequence task
	go LaunchSequenceTask()
	
	// Keep main running
	select {}
}
