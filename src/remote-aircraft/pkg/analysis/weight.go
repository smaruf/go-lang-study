package analysis

import "math"

// Point3D represents a 3D point
type Point3D struct {
	X float64 // mm
	Y float64 // mm
	Z float64 // mm
}

// Component represents a component with weight and position
type Component struct {
	Name     string  // Component name
	Weight   float64 // Weight in grams
	Position Point3D // Position in mm
}

// WeightAnalysis performs weight and center of gravity calculations
type WeightAnalysis struct {
	Components []Component
}

// NewWeightAnalysis creates a new weight analysis
func NewWeightAnalysis() *WeightAnalysis {
	return &WeightAnalysis{
		Components: make([]Component, 0),
	}
}

// AddComponent adds a component to the analysis
func (w *WeightAnalysis) AddComponent(name string, weight float64, x, y, z float64) {
	w.Components = append(w.Components, Component{
		Name:   name,
		Weight: weight,
		Position: Point3D{
			X: x,
			Y: y,
			Z: z,
		},
	})
}

// TotalWeight calculates the total weight of all components
func (w *WeightAnalysis) TotalWeight() float64 {
	total := 0.0
	for _, comp := range w.Components {
		total += comp.Weight
	}
	return total
}

// CenterOfGravity calculates the 3D center of gravity
func (w *WeightAnalysis) CenterOfGravity() Point3D {
	if len(w.Components) == 0 {
		return Point3D{}
	}
	
	totalWeight := w.TotalWeight()
	if totalWeight == 0 {
		return Point3D{}
	}
	
	cg := Point3D{}
	
	for _, comp := range w.Components {
		cg.X += comp.Position.X * comp.Weight
		cg.Y += comp.Position.Y * comp.Weight
		cg.Z += comp.Position.Z * comp.Weight
	}
	
	cg.X /= totalWeight
	cg.Y /= totalWeight
	cg.Z /= totalWeight
	
	return cg
}

// WeightDistribution calculates weight distribution percentages
func (w *WeightAnalysis) WeightDistribution() map[string]float64 {
	totalWeight := w.TotalWeight()
	if totalWeight == 0 {
		return make(map[string]float64)
	}
	
	distribution := make(map[string]float64)
	for _, comp := range w.Components {
		distribution[comp.Name] = (comp.Weight / totalWeight) * 100
	}
	
	return distribution
}

// IsBalanced checks if the aircraft is balanced (CG near origin)
func (w *WeightAnalysis) IsBalanced(tolerance float64) bool {
	cg := w.CenterOfGravity()
	distance := math.Sqrt(cg.X*cg.X + cg.Y*cg.Y)
	return distance <= tolerance
}

// Summary returns a summary of the weight analysis
func (w *WeightAnalysis) Summary() map[string]interface{} {
	cg := w.CenterOfGravity()
	
	return map[string]interface{}{
		"total_weight_g":   w.TotalWeight(),
		"num_components":   len(w.Components),
		"cg_x_mm":          cg.X,
		"cg_y_mm":          cg.Y,
		"cg_z_mm":          cg.Z,
		"is_balanced_10mm": w.IsBalanced(10.0),
		"distribution":     w.WeightDistribution(),
	}
}
