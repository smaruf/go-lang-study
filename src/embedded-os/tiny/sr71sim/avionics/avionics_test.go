package avionics

import (
    "testing"
)

func TestAvionicsNew(t *testing.T) {
    avio := New()
    if avio == nil {
        t.Error("New() should return a non-nil Avionics instance")
    }
    
    state := avio.GetState()
    if state.NavigationSystem != "GPS" {
        t.Errorf("Expected default navigation system to be GPS, got %s", state.NavigationSystem)
    }
}

func TestAvionicsSetAltitude(t *testing.T) {
    avio := New()
    avio.SetAltitude(50000)
    
    state := avio.GetState()
    if state.Altitude != 50000 {
        t.Errorf("Expected altitude 50000, got %f", state.Altitude)
    }
}

func TestAvionicsSetSpeed(t *testing.T) {
    avio := New()
    avio.SetSpeed(2200)
    
    state := avio.GetState()
    if state.Speed != 2200 {
        t.Errorf("Expected speed 2200, got %f", state.Speed)
    }
}

func TestAvionicsAutopilot(t *testing.T) {
    avio := New()
    
    avio.EnableAutopilot()
    if avio.GetState().AutopilotStatus != "Engaged" {
        t.Error("Autopilot should be Engaged")
    }
    
    avio.DisableAutopilot()
    if avio.GetState().AutopilotStatus != "Disengaged" {
        t.Error("Autopilot should be Disengaged")
    }
}

func TestAvionicsTest(t *testing.T) {
    states, err := Test()
    if err != nil {
        t.Fatalf("Test() returned error: %v", err)
    }
    
    if len(states) != 5 {
        t.Errorf("Expected 5 test states, got %d", len(states))
    }
}
