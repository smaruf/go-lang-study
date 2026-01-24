package avionics

import (
	"encoding/json"
	"fmt"
	"math"
)

// Coordinates represents GPS coordinates
type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// AvionicsState represents the state of aircraft avionics systems
type AvionicsState struct {
	Altitude            float64     `json:"altitude"`
	Speed               float64     `json:"speed"`
	Position            Coordinates `json:"position"`
	NavigationSystem    string      `json:"navigation_system"`
	CommunicationStatus string      `json:"communication_status"`
	AutopilotStatus     string      `json:"autopilot_status"`
	EngineChokeRecovery bool        `json:"engine_choke_recovery"`
	CabinPressure       float64     `json:"cabin_pressure"`
	GForceRecovery      bool        `json:"g_force_recovery"`
	FuelLeachingRate    float64     `json:"fuel_leaching_rate"`
	ExternalHeat        float64     `json:"external_heat"`
	Temperature         float64     `json:"temperature"`
	FuelSafety          bool        `json:"fuel_safety"`
}

// Avionics represents the avionics system controller
type Avionics struct {
	currentState AvionicsState
	waypoints    []Coordinates
	currentWP    int
}

// New creates a new Avionics instance with default values
func New() *Avionics {
	return &Avionics{
		currentState: AvionicsState{
			Altitude:            0,
			Speed:               0,
			Position:            Coordinates{Latitude: 0, Longitude: 0},
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
		waypoints: []Coordinates{},
		currentWP: 0,
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

// SetPosition sets the current GPS position
func (a *Avionics) SetPosition(lat, lon float64) {
	a.currentState.Position = Coordinates{Latitude: lat, Longitude: lon}
}

// GetPosition returns the current GPS position
func (a *Avionics) GetPosition() Coordinates {
	return a.currentState.Position
}

// SetWaypoints sets the flight waypoints
func (a *Avionics) SetWaypoints(waypoints []Coordinates) {
	a.waypoints = waypoints
	a.currentWP = 0
}

// GetCurrentWaypoint returns the current waypoint
func (a *Avionics) GetCurrentWaypoint() (Coordinates, bool) {
	if a.currentWP >= len(a.waypoints) {
		return Coordinates{}, false
	}
	return a.waypoints[a.currentWP], true
}

// NextWaypoint advances to the next waypoint
func (a *Avionics) NextWaypoint() bool {
	a.currentWP++
	return a.currentWP < len(a.waypoints)
}

// DistanceToWaypoint calculates distance to current waypoint in nautical miles
func (a *Avionics) DistanceToWaypoint() float64 {
	if a.currentWP >= len(a.waypoints) {
		return 0
	}
	return CalculateDistance(a.currentState.Position, a.waypoints[a.currentWP])
}

// CalculateDistance calculates distance between two coordinates in nautical miles
// Using the Haversine formula
func CalculateDistance(c1, c2 Coordinates) float64 {
	const earthRadiusNM = 3440.065 // Earth's radius in nautical miles
	
	lat1 := c1.Latitude * math.Pi / 180
	lat2 := c2.Latitude * math.Pi / 180
	deltaLat := (c2.Latitude - c1.Latitude) * math.Pi / 180
	deltaLon := (c2.Longitude - c1.Longitude) * math.Pi / 180
	
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	
	return earthRadiusNM * c
}

// NavigateToWaypoint updates position towards current waypoint
// Returns true if waypoint is reached
func (a *Avionics) NavigateToWaypoint(speedMph, timeMinutes float64) bool {
	waypoint, ok := a.GetCurrentWaypoint()
	if !ok {
		return true // No more waypoints
	}
	
	// Calculate distance we can travel in nautical miles
	// 1 mph ≈ 0.868976 knots
	speedKnots := speedMph * 0.868976
	distanceNM := speedKnots * (timeMinutes / 60.0)
	
	distanceToWP := CalculateDistance(a.currentState.Position, waypoint)
	
	if distanceNM >= distanceToWP {
		// We reached the waypoint
		a.currentState.Position = waypoint
		return true
	}
	
	// Move towards waypoint
	bearing := CalculateBearing(a.currentState.Position, waypoint)
	newPos := CalculateDestination(a.currentState.Position, distanceNM, bearing)
	a.currentState.Position = newPos
	
	return false
}

// CalculateBearing calculates the bearing from c1 to c2 in degrees
func CalculateBearing(c1, c2 Coordinates) float64 {
	lat1 := c1.Latitude * math.Pi / 180
	lat2 := c2.Latitude * math.Pi / 180
	deltaLon := (c2.Longitude - c1.Longitude) * math.Pi / 180
	
	y := math.Sin(deltaLon) * math.Cos(lat2)
	x := math.Cos(lat1)*math.Sin(lat2) - math.Sin(lat1)*math.Cos(lat2)*math.Cos(deltaLon)
	bearing := math.Atan2(y, x)
	
	return math.Mod(bearing*180/math.Pi+360, 360)
}

// CalculateDestination calculates a new position given start position, distance (NM), and bearing (degrees)
func CalculateDestination(start Coordinates, distanceNM, bearing float64) Coordinates {
	const earthRadiusNM = 3440.065
	
	lat1 := start.Latitude * math.Pi / 180
	lon1 := start.Longitude * math.Pi / 180
	brng := bearing * math.Pi / 180
	
	lat2 := math.Asin(math.Sin(lat1)*math.Cos(distanceNM/earthRadiusNM) +
		math.Cos(lat1)*math.Sin(distanceNM/earthRadiusNM)*math.Cos(brng))
	
	lon2 := lon1 + math.Atan2(math.Sin(brng)*math.Sin(distanceNM/earthRadiusNM)*math.Cos(lat1),
		math.Cos(distanceNM/earthRadiusNM)-math.Sin(lat1)*math.Sin(lat2))
	
	return Coordinates{
		Latitude:  lat2 * 180 / math.Pi,
		Longitude: lon2 * 180 / math.Pi,
	}
}
