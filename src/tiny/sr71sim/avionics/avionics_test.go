package avionics

import (
    "fmt"
)

type AvionicsState struct {
    Altitude            float64
    Speed               float64
    NavigationSystem    string
    CommunicationStatus string
    AutopilotStatus     string
    EngineChokeRecovery bool
    CabinPressure       float64
    GForceRecovery      bool
    FuelLeachingRate    float64
    ExternalHeat        float64
    Temperature         float64
    FuelSafety          bool
}

func test() {
    fmt.Println("Start SR-71 Avionics Test...")

    states := []AvionicsState{
        {Altitude: 10000, Speed: 300, NavigationSystem: "GPS", CommunicationStatus: "Active", AutopilotStatus: "Engaged", EngineChokeRecovery: true, CabinPressure: 10.5, GForceRecovery: true, FuelLeachingRate: 0.5, ExternalHeat: 45, Temperature: 15, FuelSafety: true},
        {Altitude: 15000, Speed: 500, NavigationSystem: "INS", CommunicationStatus: "Active", AutopilotStatus: "Disengaged", EngineChokeRecovery: false, CabinPressure: 9.8, GForceRecovery: true, FuelLeachingRate: 0.4, ExternalHeat: 50, Temperature: 20, FuelSafety: true},
        {Altitude: 20000, Speed: 700, NavigationSystem: "GPS", CommunicationStatus: "Inactive", AutopilotStatus: "Engaged", EngineChokeRecovery: true, CabinPressure: 8.5, GForceRecovery: false, FuelLeachingRate: 0.3, ExternalHeat: 55, Temperature: 25, FuelSafety: false},
        {Altitude: 25000, Speed: 900, NavigationSystem: "INS", CommunicationStatus: "Active", AutopilotStatus: "Engaged", EngineChokeRecovery: false, CabinPressure: 7.9, GForceRecovery: true, FuelLeachingRate: 0.2, ExternalHeat: 60, Temperature: 30, FuelSafety: true},
        {Altitude: 30000, Speed: 1100, NavigationSystem: "GPS", CommunicationStatus: "Active", AutopilotStatus: "Disengaged", EngineChokeRecovery: true, CabinPressure: 7.2, GForceRecovery: false, FuelLeachingRate: 0.1, ExternalHeat: 65, Temperature: 35, FuelSafety: false},
    }

    for _, state := range states {
        fmt.Printf("Testing at Altitude %.1f meters and Speed %.1f km/h:\n", state.Altitude, state.Speed)
        fmt.Printf("  Navigation System: %s\n", state.NavigationSystem)
        fmt.Printf("  Communication Status: %s\n", state.CommunicationStatus)
        fmt.Printf("  Autopilot Status: %s\n", state.AutopilotStatus)
        fmt.Printf("  Engine Choke Recovery: %t\n", state.EngineChokeRecovery)
        fmt.Printf("  Cabin Pressure: %.1f psi\n", state.CabinPressure)
        fmt.Printf("  G-Force Recovery: %t\n", state.GForceRecovery)
        fmt.Printf("  Fuel Leaching Rate: %.1f\n", state.FuelLeachingRate)
        fmt.Printf("  External Heat: %.1f°C\n", state.ExternalHeat)
        fmt.Printf("  Temperature: %.1f°C\n", state.Temperature)
        fmt.Printf("  Fuel Safety: %t\n", state.FuelSafety)
    }
}
