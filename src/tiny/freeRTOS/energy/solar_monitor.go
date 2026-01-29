package main

import (
	"machine"
	"time"
)

// Solar panel monitoring and Maximum Power Point Tracking (MPPT) system
// Monitors solar irradiance, panel temperature, and optimizes power output

const (
	// Solar panel measurement pins
	PANEL_VOLTAGE_ADC = machine.ADC0
	PANEL_CURRENT_ADC = machine.ADC1
	BATTERY_VOLT_ADC  = machine.ADC2
	
	// Environmental sensors
	LIGHT_SENSOR_ADC = machine.ADC3
	TEMP_SENSOR_ADC  = machine.ADC4
	
	// MPPT control pins
	PWM_CONTROL_PIN = machine.GPIO16
	
	// Relay control for load switching
	LOAD_RELAY_PIN = machine.GPIO17
	
	// Status indicators
	CHARGING_LED  = machine.GPIO18
	FULL_LED      = machine.GPIO19
	LOW_BATT_LED  = machine.LED
)

// SolarController manages solar panel and battery charging
type SolarController struct {
	panelVoltage    float32
	panelCurrent    float32
	panelPower      float32
	batteryVoltage  float32
	batteryCapacity float32
	temperature     float32
	irradiance      float32
	totalEnergy     float32
	mpptDutyCycle   uint32
	chargeState     string
	panelVoltADC    machine.ADC
	panelCurrADC    machine.ADC
	batteryVoltADC  machine.ADC
	lightSensorADC  machine.ADC
	tempSensorADC   machine.ADC
	pwmControl      machine.Pin
	loadRelay       machine.Pin
	chargingLED     machine.Pin
	fullLED         machine.Pin
	lowBattLED      machine.Pin
}

// NewSolarController creates a new solar controller
func NewSolarController() *SolarController {
	sc := &SolarController{
		mpptDutyCycle:   50,
		chargeState:     "INIT",
		totalEnergy:     0,
		batteryCapacity: 0, // Initial state of charge unknown
	}
	
	// Configure ADC pins
	machine.InitADC()
	sc.panelVoltADC = PANEL_VOLTAGE_ADC
	sc.panelCurrADC = PANEL_CURRENT_ADC
	sc.batteryVoltADC = BATTERY_VOLT_ADC
	sc.lightSensorADC = LIGHT_SENSOR_ADC
	sc.tempSensorADC = TEMP_SENSOR_ADC
	
	sc.panelVoltADC.Configure(machine.ADCConfig{})
	sc.panelCurrADC.Configure(machine.ADCConfig{})
	sc.batteryVoltADC.Configure(machine.ADCConfig{})
	sc.lightSensorADC.Configure(machine.ADCConfig{})
	sc.tempSensorADC.Configure(machine.ADCConfig{})
	
	// Configure control pins
	sc.pwmControl = PWM_CONTROL_PIN
	sc.loadRelay = LOAD_RELAY_PIN
	
	sc.pwmControl.Configure(machine.PinConfig{Mode: machine.PinPWM})
	sc.loadRelay.Configure(machine.PinConfig{Mode: machine.PinOutput})
	
	// Configure status LEDs
	sc.chargingLED = CHARGING_LED
	sc.fullLED = FULL_LED
	sc.lowBattLED = LOW_BATT_LED
	
	sc.chargingLED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	sc.fullLED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	sc.lowBattLED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	
	return sc
}

// MeasurePanelOutput measures solar panel voltage and current
func (sc *SolarController) MeasurePanelOutput() {
	// Read panel voltage (0-50V range)
	voltageRaw := sc.panelVoltADC.Get()
	sc.panelVoltage = float32(voltageRaw) * 50.0 / 65535.0
	
	// Read panel current (0-10A range)
	currentRaw := sc.panelCurrADC.Get()
	sc.panelCurrent = float32(currentRaw) * 10.0 / 65535.0
	
	// Calculate power
	sc.panelPower = sc.panelVoltage * sc.panelCurrent
}

// MeasureBatteryVoltage measures battery voltage
func (sc *SolarController) MeasureBatteryVoltage() {
	// Read battery voltage (0-15V range for 12V system)
	voltageRaw := sc.batteryVoltADC.Get()
	sc.batteryVoltage = float32(voltageRaw) * 15.0 / 65535.0
}

// MeasureEnvironment measures temperature and light intensity
func (sc *SolarController) MeasureEnvironment() {
	// Read light sensor (0-1000 W/m² irradiance)
	lightRaw := sc.lightSensorADC.Get()
	sc.irradiance = float32(lightRaw) * 1000.0 / 65535.0
	
	// Read temperature sensor (TMP36: -40 to 125°C)
	tempRaw := sc.tempSensorADC.Get()
	tempVoltage := float32(tempRaw) * 3.3 / 65535.0
	sc.temperature = (tempVoltage - 0.5) * 100.0
}

// MPPT performs Maximum Power Point Tracking
func (sc *SolarController) MPPT() {
	// Perturb and Observe algorithm
	static := struct {
		prevPower float32
		prevDuty  uint32
		direction int
	}{0, 50, 1}
	
	powerDiff := sc.panelPower - static.prevPower
	
	if powerDiff > 0 {
		// Power increased, continue in same direction
		if static.direction > 0 {
			sc.mpptDutyCycle++
		} else {
			sc.mpptDutyCycle--
		}
	} else {
		// Power decreased, reverse direction
		static.direction = -static.direction
		if static.direction > 0 {
			sc.mpptDutyCycle++
		} else {
			sc.mpptDutyCycle--
		}
	}
	
	// Limit duty cycle range
	if sc.mpptDutyCycle > 95 {
		sc.mpptDutyCycle = 95
	} else if sc.mpptDutyCycle < 10 {
		sc.mpptDutyCycle = 10
	}
	
	// Apply PWM duty cycle
	// In real implementation, use proper PWM peripheral
	
	static.prevPower = sc.panelPower
	static.prevDuty = sc.mpptDutyCycle
}

// ManageCharging manages battery charging state
func (sc *SolarController) ManageCharging() {
	// Determine charge state based on battery voltage
	if sc.batteryVoltage < 11.5 {
		// Low battery - disconnect loads
		sc.chargeState = "LOW_BATTERY"
		sc.loadRelay.Low()
		sc.lowBattLED.High()
		sc.chargingLED.Low()
		sc.fullLED.Low()
		
	} else if sc.batteryVoltage < 14.4 {
		// Bulk/absorption charging
		sc.chargeState = "CHARGING"
		sc.loadRelay.High()
		sc.chargingLED.High()
		sc.lowBattLED.Low()
		sc.fullLED.Low()
		
		// Accumulate energy
		sc.totalEnergy += sc.panelPower * 0.1 / 3600.0 // Wh
		
	} else if sc.batteryVoltage < 14.7 {
		// Float charging
		sc.chargeState = "FLOAT"
		sc.loadRelay.High()
		sc.chargingLED.Low()
		sc.fullLED.High()
		sc.lowBattLED.Low()
		
	} else {
		// Over-voltage - stop charging
		sc.chargeState = "OVER_VOLTAGE"
		sc.loadRelay.Low()
		sc.chargingLED.Low()
		sc.fullLED.Low()
		sc.lowBattLED.High()
	}
}

// CalculateEfficiency calculates system efficiency
func (sc *SolarController) CalculateEfficiency() float32 {
	// Theoretical maximum power based on irradiance
	// Standard 100W panel at 1000 W/m²
	theoreticalPower := sc.irradiance * 0.1 // 100W panel
	
	if theoreticalPower > 0 {
		efficiency := (sc.panelPower / theoreticalPower) * 100.0
		return efficiency
	}
	return 0
}

// LogData logs system data
func (sc *SolarController) LogData() {
	// In real system, log to SD card or transmit via network
	efficiency := sc.CalculateEfficiency()
	_ = efficiency
	
	// Data to log: timestamp, panel V/I/P, battery V, irradiance,
	// temperature, charge state, total energy, efficiency
}

// SolarMonitorTask is the main FreeRTOS task for solar monitoring
func SolarMonitorTask() {
	controller := NewSolarController()
	
	ticker := time.NewTicker(100 * time.Millisecond) // 10 Hz update rate
	defer ticker.Stop()
	
	for range ticker.C {
		// Measure panel output
		controller.MeasurePanelOutput()
		
		// Measure battery voltage
		controller.MeasureBatteryVoltage()
		
		// Measure environmental conditions
		controller.MeasureEnvironment()
		
		// Perform MPPT
		if controller.irradiance > 100 {
			controller.MPPT()
		}
		
		// Manage battery charging
		controller.ManageCharging()
		
		// Log data
		controller.LogData()
	}
}

func main() {
	// Run solar monitoring task
	go SolarMonitorTask()
	
	// Keep main running
	select {}
}
