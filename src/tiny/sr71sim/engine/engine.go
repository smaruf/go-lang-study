package engine

import (
	"encoding/json"
	"fmt"
)

// EngineState represents the state of a jet engine
type EngineState struct {
	Velocity          float64 `json:"velocity"`          // Mach number
	AirIntake         float64 `json:"air_intake"`        // Air intake percentage
	Altitude          float64 `json:"altitude"`          // Altitude in feet
	CombustionChamber float64 `json:"combustion_chamber"` // Temperature in Fahrenheit
	ExhaustPattern    string  `json:"exhaust_pattern"`   // Exhaust pattern description
	Mode              string  `json:"mode"`              // Engine mode: Turbojet, Ramjet, Scramjet
	Thrust            float64 `json:"thrust"`            // Thrust in pounds
	FuelFlow          float64 `json:"fuel_flow"`         // Fuel flow in gallons/hour
}

// Engine represents the SR-71's Pratt & Whitney J58 engine
type Engine struct {
	currentState EngineState
}

// New creates a new Engine instance
func New() *Engine {
	return &Engine{
		currentState: EngineState{
			Velocity:          0,
			AirIntake:         0,
			Altitude:          0,
			CombustionChamber: 200,
			ExhaustPattern:    "Idle",
			Mode:              "Turbojet",
			Thrust:            0,
			FuelFlow:          0,
		},
	}
}

// GetState returns the current engine state
func (e *Engine) GetState() EngineState {
	return e.currentState
}

// SetSpeed sets the engine speed (Mach number) and updates related parameters
func (e *Engine) SetSpeed(mach float64) {
	e.currentState.Velocity = mach
	e.updateEngineMode()
	e.updateParameters()
}

// SetAltitude sets the altitude and updates related parameters
func (e *Engine) SetAltitude(altitude float64) {
	e.currentState.Altitude = altitude
	e.updateParameters()
}

// Mode returns the current engine mode
func (e *Engine) Mode() string {
	return e.currentState.Mode
}

// updateEngineMode determines the engine mode based on Mach number
func (e *Engine) updateEngineMode() {
	mach := e.currentState.Velocity
	
	if mach < 2.0 {
		e.currentState.Mode = "Turbojet"
		e.currentState.ExhaustPattern = "Normal"
	} else if mach < 5.0 {
		e.currentState.Mode = "Ramjet"
		e.currentState.ExhaustPattern = "Supersonic"
	} else {
		e.currentState.Mode = "Scramjet"
		e.currentState.ExhaustPattern = "Hypersonic"
	}
}

// updateParameters updates all engine parameters based on current state
func (e *Engine) updateParameters() {
	mach := e.currentState.Velocity
	
	// Air intake calculation
	switch e.currentState.Mode {
	case "Turbojet":
		e.currentState.AirIntake = mach * 2.0
	case "Ramjet":
		e.currentState.AirIntake = mach * 2.2
	case "Scramjet":
		e.currentState.AirIntake = mach * 2.4
	}
	
	// Combustion chamber temperature
	if mach < 3.0 {
		e.currentState.CombustionChamber = 500 + mach*100
	} else if mach < 6.5 {
		e.currentState.CombustionChamber = 800 + (mach-3.0)*150
	} else {
		e.currentState.CombustionChamber = 1200 + (mach-6.5)*100
	}
	
	// Thrust calculation (simplified)
	// SR-71 J58 produces approximately 32,500 lbf with afterburner
	e.currentState.Thrust = 15000 + mach*5000
	if e.currentState.Thrust > 34000 {
		e.currentState.Thrust = 34000
	}
	
	// Fuel flow calculation
	// At cruise (Mach 3), consumes ~5,600 gal/hour per engine
	if mach < 1.0 {
		e.currentState.FuelFlow = 1000 + mach*500
	} else if mach < 3.0 {
		e.currentState.FuelFlow = 1500 + (mach-1.0)*2000
	} else {
		e.currentState.FuelFlow = 5500 + (mach-3.0)*500
	}
}

// Update performs a full state update
func (e *Engine) Update() {
	e.updateEngineMode()
	e.updateParameters()
}

// Test returns test data for the engine system
func Test() ([]EngineState, error) {
	fmt.Println("Start SR-71 Single Engine Test....")

	var states []EngineState
	for mach := 0.5; mach <= 15.5; mach += 0.2 {
		engine := New()
		engine.SetSpeed(mach)
		
		// Calculate altitude based on typical flight profile
		var altitude float64
		switch {
		case mach < 3.0:
			altitude = 5000 + mach*1000
		case mach < 6.5:
			altitude = 15000 + (mach-3.0)*2000
		case mach <= 15.5:
			altitude = 25000 + (mach-6.5)*1000
		}
		engine.SetAltitude(altitude)
		
		states = append(states, engine.GetState())
	}

	return states, nil
}

// TestJSON returns test data as JSON string
func TestJSON() (string, error) {
	states, err := Test()
	if err != nil {
		return "", err
	}

	jsonData, err := json.Marshal(states)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}
