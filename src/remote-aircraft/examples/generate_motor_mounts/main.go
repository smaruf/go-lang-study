package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/smaruf/go-lang-study/src/remote-aircraft/pkg/materials"
	"github.com/smaruf/go-lang-study/src/remote-aircraft/pkg/parts"
)

func main() {
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println("    Motor Mount Generator")
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println()

	// Common motor sizes (diameter in mm)
	motorSizes := []struct {
		name     string
		diameter float64
	}{
		{"1507", 15.0},
		{"1806", 18.0},
		{"2204", 22.0},
		{"2207", 22.0},
		{"2306", 23.0},
		{"2407", 24.0},
		{"2806", 28.0},
	}

	allMounts := make(map[string]interface{})

	for _, motor := range motorSizes {
		fmt.Printf("Generating mount for %s motor (%.1fmm):\n", motor.name, motor.diameter)

		mount := parts.NewMotorMount(motor.diameter, materials.PETG)
		specs := mount.Specifications()

		fmt.Printf("  Dimensions: %.1f x %.1f x %.1f mm\n",
			specs["width"], specs["height"], specs["thickness"])
		fmt.Printf("  Weight: %.2f g\n", specs["weight_g"])
		fmt.Printf("  Material: %s\n", specs["material"])

		// Test stress at 500g thrust
		stress := mount.CalculateStress(500.0)
		fmt.Printf("  Stress at 500g thrust: %.2f MPa\n", stress)
		fmt.Println()

		allMounts[motor.name] = specs
	}

	// Save all specifications to JSON
	jsonData, err := json.MarshalIndent(allMounts, "", "  ")
	if err != nil {
		fmt.Printf("Error creating JSON: %v\n", err)
		return
	}

	outputDir := "../../output"
	os.MkdirAll(outputDir, 0755)
	outputFile := outputDir + "/motor_mounts.json"
	
	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}

	fmt.Printf("Motor mount specifications saved to: %s\n", outputFile)
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════")
}
