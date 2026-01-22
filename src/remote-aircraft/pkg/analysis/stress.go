package analysis

import "math"

// StressAnalysis performs structural stress calculations
type StressAnalysis struct {
	MaterialName     string  // Material name
	YieldStrength    float64 // Yield strength in MPa
	SafetyFactor     float64 // Safety factor
}

// NewStressAnalysis creates a new stress analysis
func NewStressAnalysis(materialName string, yieldStrength float64) *StressAnalysis {
	return &StressAnalysis{
		MaterialName:  materialName,
		YieldStrength: yieldStrength,
		SafetyFactor:  2.0, // Default safety factor
	}
}

// CalculateBeamStress calculates bending stress in a beam
// length: beam length in mm
// width: beam width in mm
// height: beam height in mm
// load: load in grams
func (s *StressAnalysis) CalculateBeamStress(length, width, height, load float64) float64 {
	// Convert load to Newtons
	loadN := load * 9.81 / 1000.0
	
	// Bending moment for cantilever beam: M = F * L
	moment := loadN * length
	
	// Distance from neutral axis
	c := height / 2
	
	// Second moment of area for rectangular cross-section: I = w * h³ / 12
	i := width * math.Pow(height, 3) / 12
	
	// Bending stress: σ = M * c / I
	stress := moment * c / i // Pa
	
	return stress / 1e6 // Convert to MPa
}

// IsSafe checks if the stress is within safe limits
func (s *StressAnalysis) IsSafe(stress float64) bool {
	allowableStress := s.YieldStrength / s.SafetyFactor
	return stress <= allowableStress
}

// MaxSafeLoad calculates the maximum safe load for a beam
func (s *StressAnalysis) MaxSafeLoad(length, width, height float64) float64 {
	allowableStress := s.YieldStrength / s.SafetyFactor
	allowableStressPa := allowableStress * 1e6 // Convert to Pa
	
	c := height / 2
	i := width * math.Pow(height, 3) / 12
	
	// M = σ * I / c
	maxMoment := allowableStressPa * i / c
	
	// F = M / L
	maxLoadN := maxMoment / length
	
	// Convert to grams
	return maxLoadN * 1000.0 / 9.81
}

// Analyze performs a complete stress analysis
func (s *StressAnalysis) Analyze(length, width, height, load float64) map[string]interface{} {
	stress := s.CalculateBeamStress(length, width, height, load)
	maxLoad := s.MaxSafeLoad(length, width, height)
	allowableStress := s.YieldStrength / s.SafetyFactor
	
	return map[string]interface{}{
		"material":            s.MaterialName,
		"yield_strength_MPa":  s.YieldStrength,
		"safety_factor":       s.SafetyFactor,
		"allowable_stress_MPa": allowableStress,
		"applied_stress_MPa":  stress,
		"is_safe":             s.IsSafe(stress),
		"max_safe_load_g":     maxLoad,
		"utilization_%":       (stress / allowableStress) * 100,
	}
}

// Common material yield strengths in MPa
const (
	PLAYieldStrength     = 50.0  // PLA yield strength
	PETGYieldStrength    = 53.0  // PETG yield strength
	NylonYieldStrength   = 75.0  // Nylon yield strength
	CFNylonYieldStrength = 150.0 // Carbon fiber nylon yield strength
)

// GetYieldStrength returns yield strength for common materials
func GetYieldStrength(materialName string) float64 {
	switch materialName {
	case "PLA":
		return PLAYieldStrength
	case "PETG":
		return PETGYieldStrength
	case "Nylon":
		return NylonYieldStrength
	case "CF-Nylon", "CFNylon":
		return CFNylonYieldStrength
	default:
		return PETGYieldStrength
	}
}
