package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/smaruf/go-lang-study/src/remote-aircraft/pkg/fixedwing"
)

func main() {
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println("    Fixed-Wing Aircraft Generator")
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println()

	// Generate multiple aircraft configurations
	configs := []struct {
		name      string
		wingspan  float64
		chord     float64
		weight    float64
	}{
		{"Trainer Small", 800.0, 150.0, 180.0},
		{"Trainer Medium", 1000.0, 200.0, 250.0},
		{"Sport", 1200.0, 180.0, 300.0},
		{"Glider", 1500.0, 150.0, 200.0},
	}

	allAircraft := make(map[string]interface{})

	for _, config := range configs {
		fmt.Printf("Generating %s:\n", config.name)

		// Create wing
		wing := fixedwing.NewWing(config.wingspan, config.chord)

		// Create fuselage (proportional to wing)
		fuselageLength := config.wingspan * 0.8
		fuselage := fixedwing.NewFuselage(fuselageLength, 60.0, 80.0)

		// Create tail
		tail := fixedwing.NewTail(wing.Span, wing.Chord)

		// Calculate performance
		cruiseSpeed := 15.0
		lift := wing.CalculateLift(cruiseSpeed)
		drag := wing.CalculateDrag(cruiseSpeed)
		stallSpeed := wing.StallSpeed(config.weight)

		fmt.Printf("  Wing Span: %.1f mm\n", wing.Span)
		fmt.Printf("  Wing Area: %.1f cm²\n", wing.AreaCm2())
		fmt.Printf("  Glide Ratio: %.1f:1\n", wing.GlideRatio())
		fmt.Printf("  Stall Speed: %.1f m/s\n", stallSpeed)
		fmt.Printf("  Lift at %.1f m/s: %.1f g\n", cruiseSpeed, lift)
		fmt.Println()

		// Compile specifications
		aircraft := map[string]interface{}{
			"name":     config.name,
			"weight_g": config.weight,
			"wing":     wing.Specifications(),
			"fuselage": fuselage.Specifications(),
			"tail":     tail.Specifications(),
			"performance": map[string]interface{}{
				"cruise_speed_m/s":   cruiseSpeed,
				"lift_at_cruise_g":   lift,
				"drag_at_cruise_g":   drag,
				"stall_speed_m/s":    stallSpeed,
				"wing_loading_g/dm2": wing.WingLoading(config.weight),
			},
		}

		allAircraft[config.name] = aircraft
	}

	// Save all specifications to JSON
	jsonData, err := json.MarshalIndent(allAircraft, "", "  ")
	if err != nil {
		fmt.Printf("Error creating JSON: %v\n", err)
		return
	}

	outputDir := "../../output"
	os.MkdirAll(outputDir, 0755)
	outputFile := outputDir + "/fixed_wing_designs.json"
	
	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}

	fmt.Printf("Fixed-wing designs saved to: %s\n", outputFile)
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════")
}
