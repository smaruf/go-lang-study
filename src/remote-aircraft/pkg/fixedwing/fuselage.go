package fixedwing

// Fuselage represents an aircraft fuselage
type Fuselage struct {
	Length float64 // Fuselage length in mm
	Width  float64 // Fuselage width in mm
	Height float64 // Fuselage height in mm
}

// NewFuselage creates a new fuselage
func NewFuselage(length, width, height float64) *Fuselage {
	return &Fuselage{
		Length: length,
		Width:  width,
		Height: height,
	}
}

// Volume calculates approximate fuselage volume in mm³
func (f *Fuselage) Volume() float64 {
	// Simplified as ellipsoid volume: V ≈ 4/3 * π * a * b * c
	// where a = length/2, b = width/2, c = height/2
	return (4.0 / 3.0) * 3.14159 * (f.Length / 2) * (f.Width / 2) * (f.Height / 2)
}

// Specifications returns a map of specifications
func (f *Fuselage) Specifications() map[string]interface{} {
	return map[string]interface{}{
		"type":       "Fuselage",
		"length_mm":  f.Length,
		"width_mm":   f.Width,
		"height_mm":  f.Height,
		"volume_cm3": f.Volume() / 1000.0,
	}
}
