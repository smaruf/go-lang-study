package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/smaruf/go-lang-study/src/remote-aircraft/pkg/analysis"
)

func main() {
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println("    Structural Stress Analysis")
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println()

	// Test different materials
	materials := []struct {
		name     string
		yieldMPa float64
	}{
		{"PLA", analysis.PLAYieldStrength},
		{"PETG", analysis.PETGYieldStrength},
		{"Nylon", analysis.NylonYieldStrength},
		{"CF-Nylon", analysis.CFNylonYieldStrength},
	}

	// Arm dimensions (typical quadcopter arm)
	length := 150.0  // mm
	width := 12.0    // mm
	height := 8.0    // mm
	load := 500.0    // grams (crash load)

	fmt.Printf("Test Configuration:\n")
	fmt.Printf("  Arm Length: %.1f mm\n", length)
	fmt.Printf("  Arm Width: %.1f mm\n", width)
	fmt.Printf("  Arm Height: %.1f mm\n", height)
	fmt.Printf("  Applied Load: %.1f g\n", load)
	fmt.Println()

	results := make(map[string]interface{})

	for _, mat := range materials {
		fmt.Printf("Testing %s:\n", mat.name)
		fmt.Printf("  Yield Strength: %.1f MPa\n", mat.yieldMPa)

		sa := analysis.NewStressAnalysis(mat.name, mat.yieldMPa)
		analysisResults := sa.Analyze(length, width, height, load)

		stress := analysisResults["applied_stress_MPa"].(float64)
		maxLoad := analysisResults["max_safe_load_g"].(float64)
		isSafe := analysisResults["is_safe"].(bool)
		utilization := analysisResults["utilization_%"].(float64)

		fmt.Printf("  Applied Stress: %.2f MPa\n", stress)
		fmt.Printf("  Max Safe Load: %.1f g\n", maxLoad)
		fmt.Printf("  Is Safe: %v\n", isSafe)
		fmt.Printf("  Utilization: %.1f%%\n", utilization)
		fmt.Println()

		results[mat.name] = analysisResults
	}

	// Save results to JSON
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Printf("Error creating JSON: %v\n", err)
		return
	}

	outputDir := "../../output"
	os.MkdirAll(outputDir, 0755)
	outputFile := outputDir + "/stress_analysis.json"
	
	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}

	fmt.Printf("Analysis saved to: %s\n", outputFile)
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════")
}
