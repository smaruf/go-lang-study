package materials

// Material represents a 3D printing material with its properties
type Material struct {
	Name    string  // Material name
	Density float64 // Density in g/cm³
}

// Common 3D printing materials
var (
	// PLA - Easy printing, good strength
	PLA = Material{
		Name:    "PLA",
		Density: 1.24, // g/cm³
	}

	// PETG - Weather resistant, durable
	PETG = Material{
		Name:    "PETG",
		Density: 1.27, // g/cm³
	}

	// Nylon - High strength, flexible
	Nylon = Material{
		Name:    "Nylon",
		Density: 1.15, // g/cm³
	}

	// CFNylon - Carbon fiber reinforced, maximum strength
	CFNylon = Material{
		Name:    "CF-Nylon",
		Density: 1.20, // g/cm³
	}
)

// AllMaterials returns a slice of all available materials
func AllMaterials() []Material {
	return []Material{PLA, PETG, Nylon, CFNylon}
}

// GetMaterialByName returns a material by its name
func GetMaterialByName(name string) Material {
	switch name {
	case "PLA":
		return PLA
	case "PETG":
		return PETG
	case "Nylon":
		return Nylon
	case "CF-Nylon", "CFNylon":
		return CFNylon
	default:
		return PETG // Default to PETG
	}
}
