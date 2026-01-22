package fixedwing

import "math"

// Loads performs aerodynamic load calculations for fixed-wing aircraft
type Loads struct {
	Weight   float64 // Aircraft weight in grams
	Airspeed float64 // Airspeed in m/s
	Wing     *Wing   // Wing reference
}

// NewLoads creates a new loads calculator
func NewLoads(weight, airspeed float64, wing *Wing) *Loads {
	return &Loads{
		Weight:   weight,
		Airspeed: airspeed,
		Wing:     wing,
	}
}

// RequiredLift calculates required lift to maintain level flight
func (l *Loads) RequiredLift() float64 {
	// In level flight, lift must equal weight
	return l.Weight
}

// RequiredThrust calculates required thrust for level flight
func (l *Loads) RequiredThrust() float64 {
	// In level flight, thrust must equal drag
	return l.Wing.CalculateDrag(l.Airspeed)
}

// PowerRequired calculates power required in watts
func (l *Loads) PowerRequired() float64 {
	// Power = Thrust * Velocity
	thrust := l.RequiredThrust()
	
	// Convert thrust from grams to Newtons
	thrustN := thrust * 9.81 / 1000.0
	
	return thrustN * l.Airspeed // Watts
}

// LoadFactor calculates the load factor (g-force) for a given lift
func (l *Loads) LoadFactor(lift float64) float64 {
	return lift / l.Weight
}

// TurnRadius calculates turn radius in meters for a given bank angle (degrees)
func (l *Loads) TurnRadius(bankAngle float64) float64 {
	// Turn radius = V² / (g * tan(bank_angle))
	g := 9.81 // m/s²
	
	// Convert bank angle to radians
	bankRad := bankAngle * math.Pi / 180.0
	
	// Calculate turn radius
	return (l.Airspeed * l.Airspeed) / (g * math.Tan(bankRad))
}

// Analyze performs a complete load analysis
func (l *Loads) Analyze() map[string]interface{} {
	lift := l.Wing.CalculateLift(l.Airspeed)
	drag := l.Wing.CalculateDrag(l.Airspeed)
	thrust := l.RequiredThrust()
	power := l.PowerRequired()
	
	return map[string]interface{}{
		"weight_g":            l.Weight,
		"airspeed_m/s":        l.Airspeed,
		"lift_at_speed_g":     lift,
		"drag_g":              drag,
		"required_thrust_g":   thrust,
		"power_required_W":    power,
		"excess_lift_g":       lift - l.Weight,
		"can_fly":             lift >= l.Weight,
		"stall_speed_m/s":     l.Wing.StallSpeed(l.Weight),
		"wing_loading_g/dm2":  l.Wing.WingLoading(l.Weight),
	}
}
