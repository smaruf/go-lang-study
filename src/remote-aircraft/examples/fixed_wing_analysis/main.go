package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/smaruf/go-lang-study/src/remote-aircraft/pkg/fixedwing"
)

func main() {
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println("    Fixed-Wing Aircraft Analysis")
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println()

	// Define wing parameters
	wing := fixedwing.NewWing(1000.0, 200.0) // 1000mm span, 200mm chord
	wing.Thickness = 12.0
	wing.Dihedral = 3.0
	wing.Airfoil = "Clark-Y"

	// Define fuselage
	fuselage := fixedwing.NewFuselage(800.0, 60.0, 80.0)

	// Define tail
	tail := fixedwing.NewTail(wing.Span, wing.Chord)

	// Aircraft parameters
	aircraftWeight := 250.0 // grams
	cruiseSpeed := 15.0     // m/s

	// Display wing specifications
	fmt.Println("Wing Specifications:")
	fmt.Printf("  Span: %.1f mm\n", wing.Span)
	fmt.Printf("  Chord: %.1f mm\n", wing.Chord)
	fmt.Printf("  Thickness: %.1f%%\n", wing.Thickness)
	fmt.Printf("  Dihedral: %.1f°\n", wing.Dihedral)
	fmt.Printf("  Airfoil: %s\n", wing.Airfoil)
	fmt.Printf("  Area: %.1f cm²\n", wing.AreaCm2())
	fmt.Printf("  Aspect Ratio: %.2f\n", wing.AspectRatio())
	fmt.Printf("  Glide Ratio: %.1f:1\n", wing.GlideRatio())
	fmt.Println()

	// Performance calculations
	fmt.Println("Performance Analysis:")
	fmt.Printf("  Aircraft Weight: %.1f g\n", aircraftWeight)
	fmt.Printf("  Cruise Speed: %.1f m/s\n", cruiseSpeed)
	fmt.Println()

	lift := wing.CalculateLift(cruiseSpeed)
	drag := wing.CalculateDrag(cruiseSpeed)
	stallSpeed := wing.StallSpeed(aircraftWeight)
	wingLoading := wing.WingLoading(aircraftWeight)

	fmt.Printf("  Lift at Cruise: %.1f g\n", lift)
	fmt.Printf("  Drag at Cruise: %.1f g\n", drag)
	fmt.Printf("  Stall Speed: %.1f m/s\n", stallSpeed)
	fmt.Printf("  Wing Loading: %.1f g/dm²\n", wingLoading)
	fmt.Println()

	// Load analysis
	loads := fixedwing.NewLoads(aircraftWeight, cruiseSpeed, wing)
	loadsAnalysis := loads.Analyze()

	fmt.Println("Detailed Load Analysis:")
	fmt.Printf("  Required Thrust: %.1f g\n", loadsAnalysis["required_thrust_g"])
	fmt.Printf("  Power Required: %.1f W\n", loadsAnalysis["power_required_W"])
	fmt.Printf("  Excess Lift: %.1f g\n", loadsAnalysis["excess_lift_g"])
	fmt.Printf("  Can Fly: %v\n", loadsAnalysis["can_fly"])
	fmt.Println()

	// Test at different speeds
	fmt.Println("Performance at Different Speeds:")
	speeds := []float64{10.0, 12.0, 15.0, 18.0, 20.0}
	for _, speed := range speeds {
		l := wing.CalculateLift(speed)
		d := wing.CalculateDrag(speed)
		fmt.Printf("  %.1f m/s: Lift=%.1fg, Drag=%.1fg\n", speed, l, d)
	}
	fmt.Println()

	// Compile complete specifications
	completeSpec := map[string]interface{}{
		"wing":     wing.Specifications(),
		"fuselage": fuselage.Specifications(),
		"tail":     tail.Specifications(),
		"performance": map[string]interface{}{
			"weight_g":           aircraftWeight,
			"cruise_speed_m/s":   cruiseSpeed,
			"lift_g":             lift,
			"drag_g":             drag,
			"stall_speed_m/s":    stallSpeed,
			"wing_loading_g/dm2": wingLoading,
		},
		"loads": loadsAnalysis,
	}

	// Save to JSON
	jsonData, err := json.MarshalIndent(completeSpec, "", "  ")
	if err != nil {
		fmt.Printf("Error creating JSON: %v\n", err)
		return
	}

	outputDir := "../../output"
	os.MkdirAll(outputDir, 0755)
	outputFile := outputDir + "/fixed_wing_analysis.json"
	
	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}

	fmt.Printf("Analysis saved to: %s\n", outputFile)
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════")
}
