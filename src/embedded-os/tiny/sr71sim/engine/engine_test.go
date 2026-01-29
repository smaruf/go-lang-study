package engine

import (
    "testing"
)

func TestEngineNew(t *testing.T) {
    eng := New()
    if eng == nil {
        t.Error("New() should return a non-nil Engine instance")
    }
    
    if eng.Mode() != "Turbojet" {
        t.Errorf("Expected default mode to be Turbojet, got %s", eng.Mode())
    }
}

func TestEngineModeSwitching(t *testing.T) {
    eng := New()
    
    // Test Turbojet mode (Mach < 2.0)
    eng.SetSpeed(1.5)
    if eng.Mode() != "Turbojet" {
        t.Errorf("Expected Turbojet mode at Mach 1.5, got %s", eng.Mode())
    }
    
    // Test Ramjet mode (Mach 2.0-5.0)
    eng.SetSpeed(3.0)
    if eng.Mode() != "Ramjet" {
        t.Errorf("Expected Ramjet mode at Mach 3.0, got %s", eng.Mode())
    }
    
    // Test Scramjet mode (Mach > 5.0)
    eng.SetSpeed(6.0)
    if eng.Mode() != "Scramjet" {
        t.Errorf("Expected Scramjet mode at Mach 6.0, got %s", eng.Mode())
    }
}

func TestEngineSetSpeed(t *testing.T) {
    eng := New()
    eng.SetSpeed(3.21)
    
    state := eng.GetState()
    if state.Velocity != 3.21 {
        t.Errorf("Expected velocity 3.21, got %f", state.Velocity)
    }
}

func TestEngineTest(t *testing.T) {
    states, err := Test()
    if err != nil {
        t.Fatalf("Test() returned error: %v", err)
    }
    
    if len(states) == 0 {
        t.Error("Test() should return at least one state")
    }
}
