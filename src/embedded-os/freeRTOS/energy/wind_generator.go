package main

import (
	"machine"
	"time"
)

// Wind turbine generator monitoring and control system
// Monitors wind speed, power output, and controls turbine operation

const (
	// Anemometer (wind speed sensor) pins
	ANEMOMETER_PIN = machine.GPIO16
	
	// Wind direction sensor (potentiometer on ADC)
	WIND_DIR_ADC = machine.ADC0
	
	// Power measurement pins
	VOLTAGE_ADC = machine.ADC1
	CURRENT_ADC = machine.ADC2
	
	// Turbine control pins
	BRAKE_PIN      = machine.GPIO17
	YAW_MOTOR_PWM  = machine.GPIO18
	YAW_MOTOR_DIR  = machine.GPIO19
	
	// Status LEDs
	POWER_LED    = machine.GPIO20
	WARNING_LED  = machine.GPIO21
	SHUTDOWN_LED = machine.LED
)

// WindTurbineController manages wind turbine operation
type WindTurbineController struct {
	windSpeed      float32
	windDirection  float32
	voltage        float32
	current        float32
	power          float32
	totalEnergy    float32
	rpm            float32
	brakePower     bool
	yawAngle       float32
	anemometer     machine.Pin
	voltageADC     machine.ADC
	currentADC     machine.ADC
	windDirADC     machine.ADC
	brakePin       machine.Pin
	yawMotorPWM    machine.Pin
	yawMotorDir    machine.Pin
	powerLED       machine.Pin
	warningLED     machine.Pin
	shutdownLED    machine.Pin
}

// NewWindTurbineController creates a new wind turbine controller
func NewWindTurbineController() *WindTurbineController {
	wtc := &WindTurbineController{
		totalEnergy: 0,
		brakePower:  false,
		yawAngle:    0,
	}
	
	// Configure anemometer pin
	wtc.anemometer = ANEMOMETER_PIN
	wtc.anemometer.Configure(machine.PinConfig{Mode: machine.PinInput})
	
	// Configure ADC pins
	machine.InitADC()
	wtc.voltageADC = VOLTAGE_ADC
	wtc.currentADC = CURRENT_ADC
	wtc.windDirADC = WIND_DIR_ADC
	wtc.voltageADC.Configure(machine.ADCConfig{})
	wtc.currentADC.Configure(machine.ADCConfig{})
	wtc.windDirADC.Configure(machine.ADCConfig{})
	
	// Configure control pins
	wtc.brakePin = BRAKE_PIN
	wtc.yawMotorPWM = YAW_MOTOR_PWM
	wtc.yawMotorDir = YAW_MOTOR_DIR
	
	wtc.brakePin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	wtc.yawMotorPWM.Configure(machine.PinConfig{Mode: machine.PinOutput})
	wtc.yawMotorDir.Configure(machine.PinConfig{Mode: machine.PinOutput})
	
	// Configure status LEDs
	wtc.powerLED = POWER_LED
	wtc.warningLED = WARNING_LED
	wtc.shutdownLED = SHUTDOWN_LED
	
	wtc.powerLED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	wtc.warningLED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	wtc.shutdownLED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	
	return wtc
}

// MeasureWindSpeed measures wind speed from anemometer
func (wtc *WindTurbineController) MeasureWindSpeed() {
	// Count pulses from anemometer over 1 second
	// Each pulse = 1 rotation, calibration factor converts to m/s
	pulseCount := 0
	start := time.Now()
	
	for time.Since(start) < time.Second {
		if wtc.anemometer.Get() {
			pulseCount++
			time.Sleep(10 * time.Millisecond) // Debounce
		}
	}
	
	// Calibration: 1 pulse/sec = 2.4 km/h = 0.667 m/s
	wtc.windSpeed = float32(pulseCount) * 0.667
	wtc.rpm = float32(pulseCount) * 60.0
}

// MeasureWindDirection measures wind direction
func (wtc *WindTurbineController) MeasureWindDirection() {
	// Read ADC value from wind vane potentiometer
	adcValue := wtc.windDirADC.Get()
	
	// Convert to degrees (0-360)
	wtc.windDirection = float32(adcValue) * 360.0 / 65535.0
}

// MeasurePower measures voltage, current, and calculates power
func (wtc *WindTurbineController) MeasurePower() {
	// Read voltage (0-100V range with voltage divider)
	voltageRaw := wtc.voltageADC.Get()
	wtc.voltage = float32(voltageRaw) * 100.0 / 65535.0
	
	// Read current (0-20A range with current sensor)
	currentRaw := wtc.currentADC.Get()
	wtc.current = float32(currentRaw) * 20.0 / 65535.0
	
	// Calculate power (W)
	wtc.power = wtc.voltage * wtc.current
	
	// Accumulate energy (Wh) - sample every 100ms
	wtc.totalEnergy += wtc.power * 0.1 / 3600.0
}

// ControlBrake activates or releases the turbine brake
func (wtc *WindTurbineController) ControlBrake(engage bool) {
	wtc.brakePower = engage
	if engage {
		wtc.brakePin.High()
		wtc.shutdownLED.High()
	} else {
		wtc.brakePin.Low()
		wtc.shutdownLED.Low()
	}
}

// AdjustYaw adjusts turbine yaw to face wind
func (wtc *WindTurbineController) AdjustYaw() {
	// Calculate yaw error (difference between turbine angle and wind direction)
	yawError := wtc.windDirection - wtc.yawAngle
	
	// Normalize to -180 to 180
	if yawError > 180 {
		yawError -= 360
	} else if yawError < -180 {
		yawError += 360
	}
	
	// Apply yaw correction if error > 5 degrees
	if yawError > 5 {
		// Turn clockwise
		wtc.yawMotorDir.High()
		wtc.yawMotorPWM.High()
		time.Sleep(100 * time.Millisecond)
		wtc.yawMotorPWM.Low()
		wtc.yawAngle += 5
	} else if yawError < -5 {
		// Turn counter-clockwise
		wtc.yawMotorDir.Low()
		wtc.yawMotorPWM.High()
		time.Sleep(100 * time.Millisecond)
		wtc.yawMotorPWM.Low()
		wtc.yawAngle -= 5
	}
	
	// Normalize yaw angle
	if wtc.yawAngle >= 360 {
		wtc.yawAngle -= 360
	} else if wtc.yawAngle < 0 {
		wtc.yawAngle += 360
	}
}

// SafetyCheck performs safety checks and controls
func (wtc *WindTurbineController) SafetyCheck() {
	// Shutdown conditions
	if wtc.windSpeed > 25.0 { // Over-speed shutdown (hurricane force winds)
		wtc.ControlBrake(true)
		wtc.warningLED.High()
		return
	}
	
	if wtc.voltage > 60.0 { // Over-voltage protection
		wtc.ControlBrake(true)
		wtc.warningLED.High()
		return
	}
	
	if wtc.current > 15.0 { // Over-current protection
		wtc.ControlBrake(true)
		wtc.warningLED.High()
		return
	}
	
	// Normal operation
	if wtc.windSpeed < 3.0 { // Cut-in wind speed
		wtc.ControlBrake(true) // Brake to prevent motoring
		wtc.powerLED.Low()
	} else {
		wtc.ControlBrake(false) // Release brake for generation
		wtc.powerLED.High()
		wtc.warningLED.Low()
	}
}

// WindTurbineTask is the main FreeRTOS task for wind turbine control
func WindTurbineTask() {
	turbine := NewWindTurbineController()
	
	ticker := time.NewTicker(100 * time.Millisecond) // 10 Hz update rate
	defer ticker.Stop()
	
	measureCount := 0
	
	for range ticker.C {
		// Measure wind speed every second
		if measureCount%10 == 0 {
			turbine.MeasureWindSpeed()
		}
		
		// Measure wind direction
		turbine.MeasureWindDirection()
		
		// Measure power output
		turbine.MeasurePower()
		
		// Adjust yaw to face wind
		turbine.AdjustYaw()
		
		// Perform safety checks
		turbine.SafetyCheck()
		
		// Log data (in real system, transmit to monitoring station)
		// Data: wind speed, direction, voltage, current, power, total energy
		
		measureCount++
	}
}

func main() {
	// Run wind turbine control task
	go WindTurbineTask()
	
	// Keep main running
	select {}
}
