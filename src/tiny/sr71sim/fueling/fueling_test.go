package fueling

import (
    "fmt"
    "testing"
)

// Dummy function to simulate fueling process
func FuelingProcess(speed, altitude float64, engineType string, missionType string) string {
    if missionType == "long-range" && engineType == "turbojet" && speed > 2.0 {
        engineType = "ramjet"
    }
    if missionType == "long-range" {
        return engineType + "_refueled"
    }
    return engineType
}

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
