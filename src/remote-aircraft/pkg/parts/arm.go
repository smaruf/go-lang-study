package parts

import (
	"math"

	"github.com/smaruf/go-lang-study/src/remote-aircraft/pkg/materials"
)

// Arm represents a quadcopter arm
type Arm struct {
	Length       float64            // Arm length in mm
	Width        float64            // Arm width in mm
	Height       float64            // Arm height in mm
	WallThickness float64           // Wall thickness for hollow arms in mm
	Material     materials.Material // Material type
	Hollow       bool               // Whether the arm is hollow
}

// NewArm creates a new arm with default dimensions
func NewArm(length float64, material materials.Material) *Arm {
	return &Arm{
		Length:       length,
		Width:        12.0, // Default width
		Height:       8.0,  // Default height
		WallThickness: 2.0, // Default wall thickness
		Material:     material,
		Hollow:       true, // Default to hollow for weight savings
	}
}

// Volume calculates the volume in mm³
func (a *Arm) Volume() float64 {
	if a.Hollow {
		// Outer volume - inner volume
		outerVolume := a.Length * a.Width * a.Height
		innerWidth := a.Width - 2*a.WallThickness
		innerHeight := a.Height - 2*a.WallThickness
		innerVolume := a.Length * innerWidth * innerHeight
		
		// Ensure inner dimensions are positive
		if innerWidth <= 0 || innerHeight <= 0 {
			return outerVolume
		}
		
		return outerVolume - innerVolume
	}
	
	// Solid arm
	return a.Length * a.Width * a.Height
}

// Weight calculates the weight in grams
func (a *Arm) Weight() float64 {
	volumeCm3 := a.Volume() / 1000.0 // Convert mm³ to cm³
	return volumeCm3 * a.Material.Density
}

// CalculateStress estimates stress under crash load
func (a *Arm) CalculateStress(loadGrams float64) float64 {
	// Convert load to Newtons
	loadN := loadGrams * 9.81 / 1000.0
	
	// Simple cantilever beam stress: σ = M*c/I
	moment := loadN * a.Length
	c := a.Height / 2
	
	// Second moment of area for rectangular cross-section
	var i float64
	if a.Hollow {
		innerWidth := a.Width - 2*a.WallThickness
		innerHeight := a.Height - 2*a.WallThickness
		if innerWidth > 0 && innerHeight > 0 {
			i = (a.Width*math.Pow(a.Height, 3) - innerWidth*math.Pow(innerHeight, 3)) / 12
		} else {
			i = a.Width * math.Pow(a.Height, 3) / 12
		}
	} else {
		i = a.Width * math.Pow(a.Height, 3) / 12
	}
	
	stress := moment * c / i // Pa
	return stress / 1e6 // Convert to MPa
}

// Specifications returns a map of specifications
func (a *Arm) Specifications() map[string]interface{} {
	return map[string]interface{}{
		"type":           "Quadcopter Arm",
		"length":         a.Length,
		"width":          a.Width,
		"height":         a.Height,
		"wall_thickness": a.WallThickness,
		"hollow":         a.Hollow,
		"material":       a.Material.Name,
		"weight_g":       a.Weight(),
		"volume_cm3":     a.Volume() / 1000.0,
	}
}
