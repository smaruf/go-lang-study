package frames

import (
	"github.com/smaruf/go-lang-study/src/remote-aircraft/pkg/materials"
	"github.com/smaruf/go-lang-study/src/remote-aircraft/pkg/parts"
)

// QuadFrame represents a complete quadcopter frame
type QuadFrame struct {
	ArmLength     float64            // Diagonal arm length in mm
	Arms          []*parts.Arm       // Four arms
	MotorMounts   []*parts.MotorMount // Four motor mounts
	CenterPlate   *CenterPlate       // Center plate
	Material      materials.Material // Material type
}

// CenterPlate represents the center plate of a quadcopter
type CenterPlate struct {
	Diameter  float64            // Plate diameter in mm
	Thickness float64            // Plate thickness in mm
	Material  materials.Material // Material type
}

// NewQuadFrame creates a complete quadcopter frame
func NewQuadFrame(armLength float64, material materials.Material) *QuadFrame {
	frame := &QuadFrame{
		ArmLength: armLength,
		Material:  material,
	}
	
	// Create four arms
	frame.Arms = make([]*parts.Arm, 4)
	for i := 0; i < 4; i++ {
		frame.Arms[i] = parts.NewArm(armLength, material)
	}
	
	// Create four motor mounts (for 2207 motors - common size)
	frame.MotorMounts = make([]*parts.MotorMount, 4)
	for i := 0; i < 4; i++ {
		frame.MotorMounts[i] = parts.NewMotorMount(28.0, material)
	}
	
	// Create center plate
	centerDiameter := armLength * 0.5 // Proportional to arm length
	if centerDiameter < 60 {
		centerDiameter = 60 // Minimum size
	}
	frame.CenterPlate = &CenterPlate{
		Diameter:  centerDiameter,
		Thickness: 3.0,
		Material:  material,
	}
	
	return frame
}

// TotalWeight calculates total frame weight in grams
func (q *QuadFrame) TotalWeight() float64 {
	weight := 0.0
	
	// Arms
	for _, arm := range q.Arms {
		weight += arm.Weight()
	}
	
	// Motor mounts
	for _, mount := range q.MotorMounts {
		weight += mount.Weight()
	}
	
	// Center plate
	weight += q.CenterPlate.Weight()
	
	return weight
}

// Weight calculates center plate weight
func (c *CenterPlate) Weight() float64 {
	// Simple circular plate volume: π * r² * h
	radius := c.Diameter / 2
	volume := 3.14159 * radius * radius * c.Thickness // mm³
	volumeCm3 := volume / 1000.0
	return volumeCm3 * c.Material.Density
}

// Wheelbase returns the motor-to-motor diagonal distance
func (q *QuadFrame) Wheelbase() float64 {
	return q.ArmLength * 2 // Diagonal measurement
}

// PropClearance estimates propeller ground clearance
func (q *QuadFrame) PropClearance(propDiameter float64) float64 {
	// Estimate assuming arms at 45 degrees
	armHeight := q.Arms[0].Height
	clearance := armHeight + 10 // 10mm above arm
	return clearance
}

// Specifications returns a map of specifications
func (q *QuadFrame) Specifications() map[string]interface{} {
	return map[string]interface{}{
		"type":                "Quadcopter Frame",
		"arm_length":          q.ArmLength,
		"wheelbase":           q.Wheelbase(),
		"center_diameter":     q.CenterPlate.Diameter,
		"center_thickness":    q.CenterPlate.Thickness,
		"material":            q.Material.Name,
		"total_weight_g":      q.TotalWeight(),
		"arm_weight_g":        q.Arms[0].Weight(),
		"motor_mount_weight":  q.MotorMounts[0].Weight(),
		"center_plate_weight": q.CenterPlate.Weight(),
		"num_arms":            len(q.Arms),
		"num_motors":          len(q.MotorMounts),
	}
}
