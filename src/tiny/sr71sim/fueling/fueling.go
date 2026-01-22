package fueling

import (
	"fmt"
	"time"
)

// FuelTank represents the aircraft's fuel tank
type FuelTank struct {
	Capacity       float64 // Maximum fuel capacity in gallons
	CurrentLevel   float64 // Current fuel level in gallons
	LeakRate       float64 // Fuel leak rate (normal for SR-71 on ground)
	Refueling      bool    // Is currently being refueled
	RefuelRate     float64 // Refueling rate in gallons/minute
}

// FuelSystem represents the complete fuel management system
type FuelSystem struct {
	MainTank       *FuelTank
	ConsumptionRate float64 // Current fuel consumption in gallons/hour
	EngineType     string  // Current engine type (turbojet/ramjet)
	MissionType    string  // Current mission type
}

// NewFuelSystem creates a new fuel system
func NewFuelSystem() *FuelSystem {
	return &FuelSystem{
		MainTank: &FuelTank{
			Capacity:     12000, // SR-71 carries ~12,000 gallons
			CurrentLevel: 10000, // Start with 10,000 gallons
			LeakRate:     0.5,   // Leaks ~0.5 gal/min on ground (by design)
			Refueling:    false,
			RefuelRate:   300,   // 300 gal/min refuel rate
		},
		ConsumptionRate: 3000,      // Default consumption
		EngineType:     "turbojet",
		MissionType:    "standard",
	}
}

// GetFuelLevel returns current fuel level as percentage
func (fs *FuelSystem) GetFuelLevel() float64 {
	return (fs.MainTank.CurrentLevel / fs.MainTank.Capacity) * 100
}

// GetFuelAmount returns current fuel amount in gallons
func (fs *FuelSystem) GetFuelAmount() float64 {
	return fs.MainTank.CurrentLevel
}

// ConsumeFuel consumes fuel based on current consumption rate
func (fs *FuelSystem) ConsumeFuel(duration time.Duration) {
	hours := duration.Hours()
	consumed := fs.ConsumptionRate * hours
	
	fs.MainTank.CurrentLevel -= consumed
	if fs.MainTank.CurrentLevel < 0 {
		fs.MainTank.CurrentLevel = 0
	}
}

// LeakFuel simulates fuel leaking (normal for SR-71 on ground)
func (fs *FuelSystem) LeakFuel(duration time.Duration) {
	minutes := duration.Minutes()
	leaked := fs.MainTank.LeakRate * minutes
	
	fs.MainTank.CurrentLevel -= leaked
	if fs.MainTank.CurrentLevel < 0 {
		fs.MainTank.CurrentLevel = 0
	}
}

// StartRefueling starts the refueling process
func (fs *FuelSystem) StartRefueling() {
	fs.MainTank.Refueling = true
	fmt.Println("Starting aerial refueling...")
}

// StopRefueling stops the refueling process
func (fs *FuelSystem) StopRefueling() {
	fs.MainTank.Refueling = false
	fmt.Println("Refueling complete")
}

// Refuel adds fuel to the tank
func (fs *FuelSystem) Refuel(duration time.Duration) {
	if !fs.MainTank.Refueling {
		return
	}
	
	minutes := duration.Minutes()
	added := fs.MainTank.RefuelRate * minutes
	
	fs.MainTank.CurrentLevel += added
	if fs.MainTank.CurrentLevel > fs.MainTank.Capacity {
		fs.MainTank.CurrentLevel = fs.MainTank.Capacity
	}
}

// UpdateEngineType updates engine type and adjusts consumption rate
func (fs *FuelSystem) UpdateEngineType(speed float64) {
	oldType := fs.EngineType
	
	if speed > 2.0 {
		fs.EngineType = "ramjet"
		fs.ConsumptionRate = 5600 // High consumption at Mach 2+
	} else {
		fs.EngineType = "turbojet"
		fs.ConsumptionRate = 3000 // Lower consumption at subsonic speeds
	}
	
	if oldType != fs.EngineType {
		fmt.Printf("Engine mode switched from %s to %s\n", oldType, fs.EngineType)
	}
}

// SetMissionType sets the mission type and may affect fuel planning
func (fs *FuelSystem) SetMissionType(mission string) {
	fs.MissionType = mission
	
	// Adjust consumption based on mission type
	switch mission {
	case "long-range":
		fs.ConsumptionRate *= 1.2 // 20% more for extended range
	case "high-speed":
		fs.ConsumptionRate *= 1.5 // 50% more for max speed
	case "stealth":
		fs.ConsumptionRate *= 0.8 // 20% less for efficient cruise
	}
}

// EstimatedFlightTime calculates remaining flight time at current consumption rate
func (fs *FuelSystem) EstimatedFlightTime() time.Duration {
	if fs.ConsumptionRate == 0 {
		return 0
	}
	
	hours := fs.MainTank.CurrentLevel / fs.ConsumptionRate
	return time.Duration(hours * float64(time.Hour))
}

// NeedsRefueling checks if fuel level is below safe threshold
func (fs *FuelSystem) NeedsRefueling() bool {
	return fs.GetFuelLevel() < 30.0 // Less than 30%
}

// IsCritical checks if fuel level is critically low
func (fs *FuelSystem) IsCritical() bool {
	return fs.GetFuelLevel() < 10.0 // Less than 10%
}

// GetStatus returns a formatted status string
func (fs *FuelSystem) GetStatus() string {
	return fmt.Sprintf(
		"Fuel System Status:\n"+
			"  Fuel Level: %.0f gallons (%.1f%%)\n"+
			"  Consumption Rate: %.0f gal/hr\n"+
			"  Engine Type: %s\n"+
			"  Mission Type: %s\n"+
			"  Refueling: %v\n"+
			"  Estimated Flight Time: %v",
		fs.MainTank.CurrentLevel,
		fs.GetFuelLevel(),
		fs.ConsumptionRate,
		fs.EngineType,
		fs.MissionType,
		fs.MainTank.Refueling,
		fs.EstimatedFlightTime().Round(time.Minute),
	)
}

// FuelingProcess is the main fueling logic function (kept for compatibility with tests)
func FuelingProcess(speed, altitude float64, engineType string, missionType string) string {
	// If speed is high enough, switch from turbojet to ramjet
	if engineType == "turbojet" && speed > 2.0 {
		engineType = "ramjet"
	}
	
	// If mission is long-range, add refueling suffix
	if missionType == "long-range" {
		return engineType + "_refueled"
	}
	
	return engineType
}
