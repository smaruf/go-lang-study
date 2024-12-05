package engine

import (
    "fmt"
)

type EngineState struct {
    Velocity          float64
    AirIntake         float64
    Altitude          float64
    CombustionChamber float64
    ExhaustPattern    string
}

func test() {
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

    for _, state := range states {
        fmt.Printf("Testing at Mach %.1f:\n", state.Velocity)
        fmt.Printf("  Air Intake: %.1f\n", state.AirIntake)
        fmt.Printf("  Altitude: %.1f meters\n", state.Altitude)
        fmt.Printf("  Combustion Chamber Heat: %.1f°C\n", state.CombustionChamber)
        fmt.Printf("  Exhaust Pattern: %s\n", state.ExhaustPattern)
    }
}
