package engine

import (
    "encoding/json"
    "fmt"
)

type EngineState struct {
    Velocity          float64 `json:"velocity"`
    AirIntake         float64 `json:"air_intake"`
    Altitude          float64 `json:"altitude"`
    CombustionChamber float64 `json:"combustion_chamber"`
    ExhaustPattern    string  `json:"exhaust_pattern"`
}

func test() (string, error) {
    fmt.Println("Start SR-71 Single Engine Test....")

    var states []EngineState
    for mach := 0.5; mach <= 15.5; mach += 0.2 {
        var airIntake, altitude, combustionChamber float64
        var exhaustPattern string

        switch {
        case mach < 3.0:
            airIntake = mach * 2
            altitude = 5000 + mach*1000
            combustionChamber = 500 + mach*10
            exhaustPattern = "Normal"
        case mach < 6.5:
            airIntake = mach * 2.2
            altitude = 15000 + (mach-3.0)*2000
            combustionChamber = 800 + (mach-3.0)*50
            exhaustPattern = "Supersonic"
        case mach <= 15.5:
            airIntake = mach * 2.4
            altitude = 25000 + (mach-6.5)*1000
            combustionChamber = 1200 + (mach-6.5)*100
            exhaustPattern = "Hypersonic"
        }

        states = append(states, EngineState{
            Velocity:          mach,
            AirIntake:         airIntake,
            Altitude:          altitude,
            CombustionChamber: combustionChamber,
            ExhaustPattern:    exhaustPattern,
        })
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
        fmt.Println("Test Data:", result)
    }
}
