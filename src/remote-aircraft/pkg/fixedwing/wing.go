package fixedwing

import "math"

// Wing represents a fixed-wing aircraft wing
type Wing struct {
	Span      float64 // Wing span in mm
	Chord     float64 // Wing chord in mm
	Thickness float64 // Wing thickness as percentage of chord
	Dihedral  float64 // Dihedral angle in degrees
	Airfoil   string  // Airfoil type (e.g., "Clark-Y", "Symmetric")
}

// NewWing creates a new wing with default values
func NewWing(span, chord float64) *Wing {
	return &Wing{
		Span:      span,
		Chord:     chord,
		Thickness: 12.0, // 12% thickness
		Dihedral:  3.0,  // 3 degrees dihedral
		Airfoil:   "Clark-Y",
	}
}

// Area calculates wing area in mm²
func (w *Wing) Area() float64 {
	return w.Span * w.Chord
}

// AreaCm2 calculates wing area in cm²
func (w *Wing) AreaCm2() float64 {
	return w.Area() / 100.0
}

// AspectRatio calculates the wing aspect ratio
func (w *Wing) AspectRatio() float64 {
	return w.Span / w.Chord
}

// CalculateLift estimates lift in grams at given airspeed (m/s)
func (w *Wing) CalculateLift(airspeed float64) float64 {
	// Simplified lift equation: L = 0.5 * ρ * V² * S * CL
	// ρ = air density (1.225 kg/m³ at sea level)
	// V = airspeed in m/s
	// S = wing area in m²
	// CL = lift coefficient (typical 0.8 for Clark-Y at moderate angle)
	
	rho := 1.225 // kg/m³
	cl := 0.8    // Dimensionless
	
	// Convert area to m²
	areaM2 := w.Area() / 1000000.0
	
	// Calculate lift in Newtons
	liftN := 0.5 * rho * airspeed * airspeed * areaM2 * cl
	
	// Convert to grams
	return liftN * 1000.0 / 9.81
}

// CalculateDrag estimates drag in grams at given airspeed (m/s)
func (w *Wing) CalculateDrag(airspeed float64) float64 {
	// Simplified drag equation: D = 0.5 * ρ * V² * S * CD
	// CD = drag coefficient (typical 0.04 for well-designed wing)
	
	rho := 1.225 // kg/m³
	cd := 0.04   // Dimensionless
	
	// Convert area to m²
	areaM2 := w.Area() / 1000000.0
	
	// Calculate drag in Newtons
	dragN := 0.5 * rho * airspeed * airspeed * areaM2 * cd
	
	// Convert to grams
	return dragN * 1000.0 / 9.81
}

// GlideRatio calculates the lift-to-drag ratio
func (w *Wing) GlideRatio() float64 {
	// Simplified L/D calculation based on aspect ratio
	// Higher aspect ratio = better glide ratio
	ar := w.AspectRatio()
	
	// Empirical formula: L/D ≈ AR * efficiency
	efficiency := 0.8 // Typical efficiency factor
	return ar * efficiency
}

// StallSpeed calculates stall speed in m/s for given weight
func (w *Wing) StallSpeed(weightGrams float64) float64 {
	// Stall speed: V_stall = sqrt(2 * W / (ρ * S * CL_max))
	// CL_max for Clark-Y is approximately 1.5
	
	rho := 1.225    // kg/m³
	clMax := 1.5    // Dimensionless
	
	// Convert weight to Newtons
	weightN := weightGrams * 9.81 / 1000.0
	
	// Convert area to m²
	areaM2 := w.Area() / 1000000.0
	
	// Calculate stall speed
	return math.Sqrt(2 * weightN / (rho * areaM2 * clMax))
}

// WingLoading calculates wing loading in g/dm²
func (w *Wing) WingLoading(weightGrams float64) float64 {
	// Wing loading = weight / area
	areaDm2 := w.Area() / 10000.0 // Convert mm² to dm²
	return weightGrams / areaDm2
}

// Specifications returns a map of specifications
func (w *Wing) Specifications() map[string]interface{} {
	return map[string]interface{}{
		"type":         "Fixed Wing",
		"span_mm":      w.Span,
		"chord_mm":     w.Chord,
		"thickness_%":  w.Thickness,
		"dihedral_deg": w.Dihedral,
		"airfoil":      w.Airfoil,
		"area_cm2":     w.AreaCm2(),
		"aspect_ratio": w.AspectRatio(),
		"glide_ratio":  w.GlideRatio(),
	}
}
