package avionics

import (
	"encoding/json"
	"fmt"
)

// AvionicsState represents the state of aircraft avionics systems
type AvionicsState struct {
	Altitude            float64 `json:"altitude"`
	Speed               float64 `json:"speed"`
	NavigationSystem    string  `json:"navigation_system"`
	CommunicationStatus string  `json:"communication_status"`
	AutopilotStatus     string  `json:"autopilot_status"`
	EngineChokeRecovery bool    `json:"engine_choke_recovery"`
	CabinPressure       float64 `json:"cabin_pressure"`
	GForceRecovery      bool    `json:"g_force_recovery"`
	FuelLeachingRate    float64 `json:"fuel_leaching_rate"`
	ExternalHeat        float64 `json:"external_heat"`
	Temperature         float64 `json:"temperature"`
	FuelSafety          bool    `json:"fuel_safety"`
}

// Avionics represents the avionics system controller
type Avionics struct {
	currentState AvionicsState
}

// New creates a new Avionics instance with default values
func New() *Avionics {
	return &Avionics{
		currentState: AvionicsState{
			Altitude:            0,
			Speed:               0,
			NavigationSystem:    "GPS",
			CommunicationStatus: "Active",
			AutopilotStatus:     "Disengaged",
			EngineChokeRecovery: true,
			CabinPressure:       14.7, // Sea level pressure
			GForceRecovery:      true,
			FuelLeachingRate:    0.0,
			ExternalHeat:        60.0,
			Temperature:         70.0,
			FuelSafety:          true,
		},
	}
}

// GetState returns the current avionics state
func (a *Avionics) GetState() AvionicsState {
	return a.currentState
}

// SetAltitude updates the current altitude
func (a *Avionics) SetAltitude(altitude float64) {
	a.currentState.Altitude = altitude
	a.updateCabinPressure()
}

// SetSpeed updates the current speed
func (a *Avionics) SetSpeed(speed float64) {
	a.currentState.Speed = speed
	a.updateExternalHeat()
}

// EnableAutopilot enables the autopilot system
func (a *Avionics) EnableAutopilot() {
	a.currentState.AutopilotStatus = "Engaged"
}

// DisableAutopilot disables the autopilot system
func (a *Avionics) DisableAutopilot() {
	a.currentState.AutopilotStatus = "Disengaged"
}

// SwitchNavigationSystem switches between GPS and INS
func (a *Avionics) SwitchNavigationSystem(system string) error {
	if system != "GPS" && system != "INS" {
		return fmt.Errorf("invalid navigation system: %s", system)
	}
	a.currentState.NavigationSystem = system
	return nil
}

// updateCabinPressure calculates cabin pressure based on altitude
func (a *Avionics) updateCabinPressure() {
	// Simplified pressure calculation: decreases with altitude
	// Sea level: 14.7 psi, at 85,000 ft: ~0.5 psi
	// Cabin is pressurized, so it doesn't drop as fast
	if a.currentState.Altitude <= 10000 {
		a.currentState.CabinPressure = 14.7
	} else {
		// Linear approximation for pressurized cabin
		a.currentState.CabinPressure = 14.7 - (a.currentState.Altitude-10000)/85000*6.0
		if a.currentState.CabinPressure < 8.0 {
			a.currentState.CabinPressure = 8.0 // Minimum cabin pressure
		}
	}
}

// updateExternalHeat calculates external heat based on speed
func (a *Avionics) updateExternalHeat() {
	// At Mach 3+, surface temps reach 600°F (316°C)
	// Speed in mph, Mach 1 ≈ 767 mph
	mach := a.currentState.Speed / 767.0
	if mach < 1.0 {
		a.currentState.ExternalHeat = 60.0 + mach*100
	} else if mach < 3.0 {
		a.currentState.ExternalHeat = 160.0 + (mach-1.0)*200
	} else {
		// At Mach 3+, heat increases dramatically
		a.currentState.ExternalHeat = 560.0 + (mach-3.0)*100
		if a.currentState.ExternalHeat > 600.0 {
			a.currentState.ExternalHeat = 600.0 // Max temp
		}
	}
}

// Update performs a full state update
func (a *Avionics) Update() {
	a.updateCabinPressure()
	a.updateExternalHeat()
	
	// Update fuel leaching rate (expected at high speeds due to thermal expansion)
	mach := a.currentState.Speed / 767.0
	if mach > 2.0 {
		a.currentState.FuelLeachingRate = (mach - 2.0) * 0.5
	} else {
		a.currentState.FuelLeachingRate = 0.0
	}
}

// Test returns test data for the avionics system
func Test() ([]AvionicsState, error) {
	fmt.Println("Start SR-71 Avionics Test...")

	states := []AvionicsState{
		{Altitude: 10000, Speed: 300, NavigationSystem: "GPS", CommunicationStatus: "Active", AutopilotStatus: "Engaged", EngineChokeRecovery: true, CabinPressure: 10.5, GForceRecovery: true, FuelLeachingRate: 1.0, ExternalHeat: 50.0, Temperature: 24.0, FuelSafety: true},
		{Altitude: 15000, Speed: 500, NavigationSystem: "INS", CommunicationStatus: "Active", AutopilotStatus: "Disengaged", EngineChokeRecovery: false, CabinPressure: 9.8, GForceRecovery: true, FuelLeachingRate: 1.2, ExternalHeat: 55.0, Temperature: 25.0, FuelSafety: true},
		{Altitude: 20000, Speed: 700, NavigationSystem: "GPS", CommunicationStatus: "Inactive", AutopilotStatus: "Engaged", EngineChokeRecovery: true, CabinPressure: 8.5, GForceRecovery: false, FuelLeachingRate: 1.4, ExternalHeat: 60.0, Temperature: 26.0, FuelSafety: false},
		{Altitude: 25000, Speed: 900, NavigationSystem: "INS", CommunicationStatus: "Active", AutopilotStatus: "Engaged", EngineChokeRecovery: false, CabinPressure: 7.9, GForceRecovery: true, FuelLeachingRate: 1.6, ExternalHeat: 65.0, Temperature: 27.0, FuelSafety: true},
		{Altitude: 30000, Speed: 1100, NavigationSystem: "GPS", CommunicationStatus: "Active", AutopilotStatus: "Disengaged", EngineChokeRecovery: true, CabinPressure: 7.2, GForceRecovery: false, FuelLeachingRate: 1.8, ExternalHeat: 70.0, Temperature: 28.0, FuelSafety: false},
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
