package flying

import (
  "fmt"
  "testing"
)

type SR71 struct {
  Altitude int
  Velocity int
}

func (sr71 *SR71) FlyAtHeight(height int) {
  sr71.Altitude = height
  fmt.Printf("SR-71 flying at %d feet\n", height)
}

func (sr71 *SR71) AdjustVelocityForMission(mission string) {
  switch mission {
  case "reconnaissance":
    sr71.Velocity = 2200 // speed in mph
  case "high-speed":
    sr71.Velocity = 2500
  case "stealth":
    sr71.Velocity = 1800
  default:
    sr71.Velocity = 2000
  }
  fmt.Printf("SR-71 adjusted velocity to %d mph for %s mission\n", sr71.Velocity, mission)
}

func TestSR71FlyingAtDifferentHeights(t *testing.T) {
  sr71 := &SR71{}

  heights := []int{10000, 20000, 30000, 40000, 50000, 60000, 70000, 80000, 85000}
  for _, height := range heights {
    sr71.FlyAtHeight(height)
    if sr71.Altitude != height {
      t.Errorf("Expected altitude %d, but got %d", height, sr71.Altitude)
    }
  }
}

func TestSR71VelocityShiftingForDifferentMissions(t *testing.T) {
  sr71 := &SR71{}

  missions := []string{"reconnaissance", "high-speed", "stealth", "default"}
  expectedVelocities := map[string]int{
    "reconnaissance": 2200,
    "high-speed":     2500,
    "stealth":        1800,
    "default":        2000,
  }

  for _, mission := range missions {
    sr71.AdjustVelocityForMission(mission)
    if sr71.Velocity != expectedVelocities[mission] {
      t.Errorf("Expected velocity %d for %s mission, but got %d", expectedVelocities[mission], mission, sr71.Velocity)
    }
  }
}
