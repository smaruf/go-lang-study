package fueling

import (
    "fmt"
    "testing"
    "time"
)

func TestFuelingProcess(t *testing.T) {
    tests := []struct {
        speed      float64
        altitude   float64
        engineType string
        missionType string
        expected   string
    }{
        {1.5, 30000, "turbojet", "standard", "turbojet"},
        {2.5, 35000, "turbojet", "standard", "ramjet"},
        {3.0, 40000, "ramjet", "standard", "ramjet"},
        {0.9, 20000, "turbojet", "standard", "turbojet"},
        {2.2, 36000, "turbojet", "standard", "ramjet"},
        {2.5, 35000, "turbojet", "long-range", "ramjet_refueled"},
        {1.5, 30000, "turbojet", "long-range", "turbojet_refueled"},
        {3.0, 40000, "ramjet", "long-range", "ramjet_refueled"},
    }

    for _, tt := range tests {
        t.Run(fmt.Sprintf("speed=%v,altitude=%v,engineType=%v,missionType=%v", tt.speed, tt.altitude, tt.engineType, tt.missionType), func(t *testing.T) {
            result := FuelingProcess(tt.speed, tt.altitude, tt.engineType, tt.missionType)
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}

func TestFuelSystemNew(t *testing.T) {
    fs := NewFuelSystem()
    if fs == nil {
        t.Error("NewFuelSystem() should return a non-nil FuelSystem instance")
    }
    
    if fs.GetFuelLevel() <= 0 {
        t.Error("New fuel system should have fuel")
    }
}

func TestFuelConsumption(t *testing.T) {
    fs := NewFuelSystem()
    initialFuel := fs.GetFuelAmount()
    
    fs.ConsumeFuel(1 * time.Hour)
    
    if fs.GetFuelAmount() >= initialFuel {
        t.Error("Fuel should decrease after consumption")
    }
}

func TestRefueling(t *testing.T) {
    fs := NewFuelSystem()
    fs.SetFuelLevel(1000) // Set low fuel using proper method
    initialFuel := fs.GetFuelAmount()
    
    fs.StartRefueling()
    fs.Refuel(1 * time.Minute)
    
    if fs.GetFuelAmount() <= initialFuel {
        t.Error("Fuel should increase during refueling")
    }
}
