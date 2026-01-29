package flying

import (
	"fmt"
)

// SR71 represents the SR-71 Blackbird aircraft
type SR71 struct {
	Altitude       int     // Current altitude in feet
	Velocity       int     // Current velocity in mph
	TargetAltitude int     // Target altitude in feet
	TargetVelocity int     // Target velocity in mph
	Pitch          float64 // Pitch angle in degrees
	Roll           float64 // Roll angle in degrees
	Yaw            float64 // Yaw angle in degrees
	Mission        string  // Current mission type
}

// Mission types
const (
	MissionReconnaissance = "reconnaissance"
	MissionHighSpeed      = "high-speed"
	MissionStealth        = "stealth"
	MissionTraining       = "training"
)

// NewSR71 creates a new SR-71 aircraft instance
func NewSR71() *SR71 {
	return &SR71{
		Altitude:       0,
		Velocity:       0,
		TargetAltitude: 0,
		TargetVelocity: 0,
		Pitch:          0,
		Roll:           0,
		Yaw:            0,
		Mission:        "",
	}
}

// FlyAtHeight sets the aircraft to fly at a specified altitude
func (sr71 *SR71) FlyAtHeight(height int) {
	sr71.TargetAltitude = height
	sr71.Altitude = height
	fmt.Printf("SR-71 flying at %d feet\n", height)
}

// AdjustVelocityForMission adjusts velocity based on mission type
func (sr71 *SR71) AdjustVelocityForMission(mission string) {
	sr71.Mission = mission
	
	switch mission {
	case MissionReconnaissance:
		sr71.Velocity = 2200 // Mach 3.2
		sr71.TargetVelocity = 2200
	case MissionHighSpeed:
		sr71.Velocity = 2500 // Maximum speed
		sr71.TargetVelocity = 2500
	case MissionStealth:
		sr71.Velocity = 1800 // Lower speed for reduced detection
		sr71.TargetVelocity = 1800
	case MissionTraining:
		sr71.Velocity = 1500 // Training speed
		sr71.TargetVelocity = 1500
	default:
		sr71.Velocity = 2000 // Default cruise speed
		sr71.TargetVelocity = 2000
	}
	
	fmt.Printf("SR-71 adjusted velocity to %d mph for %s mission\n", sr71.Velocity, mission)
}

// ClimbTo gradually climbs to target altitude
func (sr71 *SR71) ClimbTo(targetAltitude int, rate int) {
	sr71.TargetAltitude = targetAltitude
	
	if rate <= 0 {
		rate = 5000 // Default climb rate: 5000 ft/min
	}
	
	if sr71.Altitude < targetAltitude {
		sr71.Altitude += rate
		if sr71.Altitude > targetAltitude {
			sr71.Altitude = targetAltitude
		}
	}
}

// DescendTo gradually descends to target altitude
func (sr71 *SR71) DescendTo(targetAltitude int, rate int) {
	sr71.TargetAltitude = targetAltitude
	
	if rate <= 0 {
		rate = 3000 // Default descent rate: 3000 ft/min
	}
	
	if sr71.Altitude > targetAltitude {
		sr71.Altitude -= rate
		if sr71.Altitude < targetAltitude {
			sr71.Altitude = targetAltitude
		}
	}
}

// Accelerate increases velocity
func (sr71 *SR71) Accelerate(amount int) {
	sr71.Velocity += amount
	if sr71.Velocity > 2500 {
		sr71.Velocity = 2500 // Maximum speed limit
	}
	fmt.Printf("SR-71 accelerated to %d mph\n", sr71.Velocity)
}

// Decelerate decreases velocity
func (sr71 *SR71) Decelerate(amount int) {
	sr71.Velocity -= amount
	if sr71.Velocity < 0 {
		sr71.Velocity = 0
	}
	fmt.Printf("SR-71 decelerated to %d mph\n", sr71.Velocity)
}

// SetPitch sets the pitch angle
func (sr71 *SR71) SetPitch(angle float64) {
	sr71.Pitch = angle
}

// SetRoll sets the roll angle
func (sr71 *SR71) SetRoll(angle float64) {
	sr71.Roll = angle
}

// SetYaw sets the yaw angle
func (sr71 *SR71) SetYaw(angle float64) {
	sr71.Yaw = angle
}

// GetMachNumber calculates current Mach number (approximate)
func (sr71 *SR71) GetMachNumber() float64 {
	// Mach 1 ≈ 767 mph at sea level
	return float64(sr71.Velocity) / 767.0
}

// GetStatus returns a formatted status string
func (sr71 *SR71) GetStatus() string {
	return fmt.Sprintf(
		"SR-71 Status:\n"+
			"  Altitude: %d ft (Target: %d ft)\n"+
			"  Velocity: %d mph (Mach %.2f)\n"+
			"  Mission: %s\n"+
			"  Attitude: Pitch=%.1f° Roll=%.1f° Yaw=%.1f°",
		sr71.Altitude, sr71.TargetAltitude,
		sr71.Velocity, sr71.GetMachNumber(),
		sr71.Mission,
		sr71.Pitch, sr71.Roll, sr71.Yaw,
	)
}
