package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/smaruf/go-lang-study/src/remote-aircraft/pkg/analysis"
	"github.com/smaruf/go-lang-study/src/remote-aircraft/pkg/frames"
	"github.com/smaruf/go-lang-study/src/remote-aircraft/pkg/materials"
)

func main() {
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println("    Quadcopter Weight & CG Analysis")
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println()

	// Create a quadcopter frame
	frame := frames.NewQuadFrame(150.0, materials.PETG)

	// Create weight analysis
	wa := analysis.NewWeightAnalysis()

	// Add frame components (positioned around center)
	wa.AddComponent("Frame", frame.TotalWeight(), 0, 0, 0)

	// Add motors (positioned at arm ends)
	armLen := frame.ArmLength
	motorWeight := 28.0 // grams each
	wa.AddComponent("Motor FL", motorWeight, armLen, armLen, 0)
	wa.AddComponent("Motor FR", motorWeight, armLen, -armLen, 0)
	wa.AddComponent("Motor RL", motorWeight, -armLen, armLen, 0)
	wa.AddComponent("Motor RR", motorWeight, -armLen, -armLen, 0)

	// Add flight controller
	wa.AddComponent("Flight Controller", 15.0, 0, 0, 5)

	// Add battery
	wa.AddComponent("Battery", 180.0, 0, 0, -10)

	// Add camera
	wa.AddComponent("Camera", 25.0, 30, 0, 5)

	// Add receiver
	wa.AddComponent("Receiver", 8.0, -20, 0, 5)

	// Display results
	fmt.Printf("Total Weight: %.2f g\n", wa.TotalWeight())
	fmt.Println()

	cg := wa.CenterOfGravity()
	fmt.Printf("Center of Gravity:\n")
	fmt.Printf("  X: %.2f mm\n", cg.X)
	fmt.Printf("  Y: %.2f mm\n", cg.Y)
	fmt.Printf("  Z: %.2f mm\n", cg.Z)
	fmt.Println()

	balanced := wa.IsBalanced(10.0)
	fmt.Printf("Is Balanced (±10mm): %v\n", balanced)
	fmt.Println()

	// Weight distribution
	fmt.Println("Weight Distribution:")
	distribution := wa.WeightDistribution()
	for name, percent := range distribution {
		fmt.Printf("  %-20s: %5.1f%%\n", name, percent)
	}
	fmt.Println()

	// Export to JSON
	summary := wa.Summary()
	jsonData, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		fmt.Printf("Error creating JSON: %v\n", err)
		return
	}

	// Save to file
	outputDir := "../../output"
	os.MkdirAll(outputDir, 0755)
	outputFile := outputDir + "/weight_analysis.json"
	
	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}

	fmt.Printf("Analysis saved to: %s\n", outputFile)
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════")
}
