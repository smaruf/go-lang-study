package fixedwing

// Tail represents tail surfaces (horizontal and vertical stabilizers)
type Tail struct {
	HorizontalSpan  float64 // Horizontal stabilizer span in mm
	HorizontalChord float64 // Horizontal stabilizer chord in mm
	VerticalHeight  float64 // Vertical stabilizer height in mm
	VerticalChord   float64 // Vertical stabilizer chord in mm
}

// NewTail creates a new tail with default proportions
func NewTail(wingSpan, wingChord float64) *Tail {
	// Typical tail sizing based on main wing
	return &Tail{
		HorizontalSpan:  wingSpan * 0.4,  // 40% of wing span
		HorizontalChord: wingChord * 0.5, // 50% of wing chord
		VerticalHeight:  wingSpan * 0.15, // 15% of wing span
		VerticalChord:   wingChord * 0.6, // 60% of wing chord
	}
}

// HorizontalArea calculates horizontal stabilizer area in mm²
func (t *Tail) HorizontalArea() float64 {
	return t.HorizontalSpan * t.HorizontalChord
}

// VerticalArea calculates vertical stabilizer area in mm²
func (t *Tail) VerticalArea() float64 {
	return t.VerticalHeight * t.VerticalChord
}

// Specifications returns a map of specifications
func (t *Tail) Specifications() map[string]interface{} {
	return map[string]interface{}{
		"type":                  "Tail Surfaces",
		"horizontal_span_mm":    t.HorizontalSpan,
		"horizontal_chord_mm":   t.HorizontalChord,
		"horizontal_area_cm2":   t.HorizontalArea() / 100.0,
		"vertical_height_mm":    t.VerticalHeight,
		"vertical_chord_mm":     t.VerticalChord,
		"vertical_area_cm2":     t.VerticalArea() / 100.0,
	}
}
