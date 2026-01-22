package parts

import (
	"math"

	"github.com/smaruf/go-lang-study/src/remote-aircraft/pkg/materials"
)

// MotorMount represents a motor mounting plate
type MotorMount struct {
	MotorDiameter float64            // Motor diameter in mm
	HoleSpacing   float64            // Mounting hole spacing in mm
	Thickness     float64            // Plate thickness in mm
	Material      materials.Material // Material type
	Width         float64            // Plate width in mm (auto-calculated if 0)
	Height        float64            // Plate height in mm (auto-calculated if 0)
}

// NewMotorMount creates a motor mount with default dimensions
func NewMotorMount(motorDiameter float64, material materials.Material) *MotorMount {
	return &MotorMount{
		MotorDiameter: motorDiameter,
		HoleSpacing:   motorDiameter * 0.68, // Typical ratio
		Thickness:     3.0,
		Material:      material,
		Width:         0,  // Auto-calculate
		Height:        0,  // Auto-calculate
	}
}

// calculateDimensions calculates width and height if not set
func (m *MotorMount) calculateDimensions() {
	if m.Width == 0 {
		m.Width = m.HoleSpacing + 10 // Add margin
	}
	if m.Height == 0 {
		m.Height = m.HoleSpacing + 10 // Add margin
	}
}

// Volume calculates the volume in mm³
func (m *MotorMount) Volume() float64 {
	m.calculateDimensions()
	
	// Simple rectangular plate
	plateVolume := m.Width * m.Height * m.Thickness
	
	// Subtract motor hole
	motorHoleVolume := math.Pi * math.Pow(m.MotorDiameter/2, 2) * m.Thickness
	
	// Subtract mounting holes (4 holes, 3mm diameter each)
	mountingHolesVolume := 4 * math.Pi * math.Pow(1.5, 2) * m.Thickness
	
	return plateVolume - motorHoleVolume - mountingHolesVolume
}

// Weight calculates the weight in grams
func (m *MotorMount) Weight() float64 {
	volumeCm3 := m.Volume() / 1000.0 // Convert mm³ to cm³
	return volumeCm3 * m.Material.Density
}

// CalculateStress estimates stress under thrust load (simplified beam bending)
func (m *MotorMount) CalculateStress(thrustGrams float64) float64 {
	m.calculateDimensions()
	
	// Convert thrust to Newtons
	thrustN := thrustGrams * 9.81 / 1000.0
	
	// Simple cantilever beam stress calculation
	// σ = M*c/I where M = F*L/4, c = thickness/2, I = w*t³/12
	length := m.Width / 2 // Effective cantilever length
	moment := thrustN * length / 4
	c := m.Thickness / 2
	i := m.Width * math.Pow(m.Thickness, 3) / 12
	
	stress := moment * c / i // Pa
	return stress / 1e6 // Convert to MPa
}

// Specifications returns a map of specifications
func (m *MotorMount) Specifications() map[string]interface{} {
	m.calculateDimensions()
	
	return map[string]interface{}{
		"type":           "Motor Mount",
		"motor_diameter": m.MotorDiameter,
		"hole_spacing":   m.HoleSpacing,
		"thickness":      m.Thickness,
		"width":          m.Width,
		"height":         m.Height,
		"material":       m.Material.Name,
		"weight_g":       m.Weight(),
		"volume_cm3":     m.Volume() / 1000.0,
	}
}
