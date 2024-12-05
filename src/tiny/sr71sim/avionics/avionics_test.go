package avionics

import (
    "encoding/json"
    "fmt"
)

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

func test() (string, error) {
    fmt.Println("Start SR-71 Avionics Test...")

    states := []AvionicsState{
        {Altitude: 10000, Speed: 300, NavigationSystem: "GPS", CommunicationStatus: "Active", AutopilotStatus: "Engaged", EngineChokeRecovery: true, CabinPressure: 10.5, GForceRecovery: true, FuelLeachingRate: 1.0, ExternalHeat: 50.0, Temperature: 24.0, FuelSafety: true},
        {Altitude: 15000, Speed: 500, NavigationSystem: "INS", CommunicationStatus: "Active", AutopilotStatus: "Disengaged", EngineChokeRecovery: false, CabinPressure: 9.8, GForceRecovery: true, FuelLeachingRate: 1.2, ExternalHeat: 55.0, Temperature: 25.0, FuelSafety: true},
        {Altitude: 20000, Speed: 700, NavigationSystem: "GPS", CommunicationStatus: "Inactive", AutopilotStatus: "Engaged", EngineChokeRecovery: true, CabinPressure: 8.5, GForceRecovery: false, FuelLeachingRate: 1.4, ExternalHeat: 60.0, Temperature: 26.0, FuelSafety: false},
        {Altitude: 25000, Speed: 900, NavigationSystem: "INS", CommunicationStatus: "Active", AutopilotStatus: "Engaged", EngineChokeRecovery: false, CabinPressure: 7.9, GForceRecovery: true, FuelLeachingRate: 1.6, ExternalHeat: 65.0, Temperature: 27.0, FuelSafety: true},
        {Altitude: 30000, Speed: 1100, NavigationSystem: "GPS", CommunicationStatus: "Active", AutopilotStatus: "Disengaged", EngineChokeRecovery: true, CabinPressure: 7.2, GForceRecovery: false, FuelLeachingRate: 1.8, ExternalHeat: 70.0, Temperature: 28.0, FuelSafety: false},
    }

    jsonData, err := json.Marshal(states)
    if err != nil {
        return "", err
    }

    return string(jsonData), nil
}

func main() {
    result, err := test()
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Simulation Data:", result)
    }
}
