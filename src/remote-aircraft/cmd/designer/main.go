package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/smaruf/go-lang-study/src/remote-aircraft/pkg/fixedwing"
	"github.com/smaruf/go-lang-study/src/remote-aircraft/pkg/materials"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	printHeader()

	for {
		choice := showMainMenu(reader)

		switch choice {
		case 1:
			designFixedWing(reader)
		case 2:
			designGlider(reader)
		case 3:
			fmt.Println("\nThank you for using Aircraft Designer!")
			return
		default:
			fmt.Println("\nInvalid choice. Please try again.")
		}
	}
}

func printHeader() {
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println("           Aircraft Designer - CLI")
	fmt.Println("    Parametric CAD for Remote Aircraft")
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println()
}

func showMainMenu(reader *bufio.Reader) int {
	fmt.Println("\nSelect Aircraft Type:")
	fmt.Println("  1. ✈️  Fixed Wing Aircraft")
	fmt.Println("  2. 🪂  Glider")
	fmt.Println("  3. Exit")
	fmt.Print("\nYour choice: ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	choice, _ := strconv.Atoi(input)

	return choice
}

func designFixedWing(reader *bufio.Reader) {
	fmt.Println("\n═══════════════════════════════════════════════════")
	fmt.Println("       Fixed Wing Aircraft Designer")
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println()

	// Get parameters from user
	wingSpan := readFloat(reader, "Wing Span (mm)", 1000.0)
	wingChord := readFloat(reader, "Wing Chord (mm)", 200.0)
	wingThickness := readFloat(reader, "Wing Thickness (%)", 12.0)
	dihedral := readFloat(reader, "Dihedral Angle (degrees)", 3.0)

	fuseLength := readFloat(reader, "Fuselage Length (mm)", 800.0)
	fuseWidth := readFloat(reader, "Fuselage Width (mm)", 60.0)
	fuseHeight := readFloat(reader, "Fuselage Height (mm)", 80.0)

	hStabSpan := readFloat(reader, "Horizontal Stabilizer Span (mm)", 400.0)
	hStabChord := readFloat(reader, "Horizontal Stabilizer Chord (mm)", 100.0)
	vStabHeight := readFloat(reader, "Vertical Stabilizer Height (mm)", 150.0)
	vStabChord := readFloat(reader, "Vertical Stabilizer Chord (mm)", 120.0)

	motorDiameter := readFloat(reader, "Motor Diameter (mm)", 28.0)
	propDiameter := readFloat(reader, "Propeller Diameter (inches)", 9.0)

	materialName := readString(reader, "Material (PLA/PETG/Nylon/CF-Nylon)", "PETG")
	material := materials.GetMaterialByName(materialName)

	estimatedWeight := readFloat(reader, "Estimated Weight (g)", 250.0)

	// Create components
	wing := &fixedwing.Wing{
		Span:      wingSpan,
		Chord:     wingChord,
		Thickness: wingThickness,
		Dihedral:  dihedral,
		Airfoil:   "Clark-Y",
	}

	fuselage := fixedwing.NewFuselage(fuseLength, fuseWidth, fuseHeight)

	tail := &fixedwing.Tail{
		HorizontalSpan:  hStabSpan,
		HorizontalChord: hStabChord,
		VerticalHeight:  vStabHeight,
		VerticalChord:   vStabChord,
	}

	// Calculate performance
	cruiseSpeed := 15.0 // m/s default
	loads := fixedwing.NewLoads(estimatedWeight, cruiseSpeed, wing)

	// Display results
	fmt.Println("\n═══════════════════════════════════════════════════")
	fmt.Println("         Design Summary")
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println()

	fmt.Println("Wing Specifications:")
	fmt.Printf("  Span: %.1f mm\n", wing.Span)
	fmt.Printf("  Chord: %.1f mm\n", wing.Chord)
	fmt.Printf("  Area: %.1f cm²\n", wing.AreaCm2())
	fmt.Printf("  Aspect Ratio: %.2f\n", wing.AspectRatio())
	fmt.Printf("  Glide Ratio: %.1f:1\n", wing.GlideRatio())
	fmt.Println()

	fmt.Println("Performance:")
	fmt.Printf("  Stall Speed: %.1f m/s\n", wing.StallSpeed(estimatedWeight))
	fmt.Printf("  Wing Loading: %.1f g/dm²\n", wing.WingLoading(estimatedWeight))
	fmt.Printf("  Lift at %.1f m/s: %.1f g\n", cruiseSpeed, wing.CalculateLift(cruiseSpeed))
	fmt.Println()

	loadsAnalysis := loads.Analyze()
	fmt.Println("Propulsion:")
	fmt.Printf("  Required Thrust: %.1f g\n", loadsAnalysis["required_thrust_g"])
	fmt.Printf("  Power Required: %.1f W\n", loadsAnalysis["power_required_W"])
	fmt.Printf("  Motor: %.1f mm diameter\n", motorDiameter)
	fmt.Printf("  Propeller: %.1f inches\n", propDiameter)
	fmt.Println()

	fmt.Println("Build:")
	fmt.Printf("  Material: %s\n", material.Name)
	fmt.Printf("  Density: %.2f g/cm³\n", material.Density)
	fmt.Println()

	// Compile complete design
	design := map[string]interface{}{
		"type":     "Fixed Wing Aircraft",
		"wing":     wing.Specifications(),
		"fuselage": fuselage.Specifications(),
		"tail":     tail.Specifications(),
		"propulsion": map[string]interface{}{
			"motor_diameter_mm": motorDiameter,
			"prop_diameter_in":  propDiameter,
		},
		"build": map[string]interface{}{
			"material": material.Name,
			"density":  material.Density,
		},
		"performance": loadsAnalysis,
	}

	// Save to file
	saveDesign(design, "fixed_wing_design")
}

func designGlider(reader *bufio.Reader) {
	fmt.Println("\n═══════════════════════════════════════════════════")
	fmt.Println("             Glider Designer")
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println()

	// Get parameters from user
	wingSpan := readFloat(reader, "Wing Span (mm)", 1500.0)
	wingChord := readFloat(reader, "Wing Chord (mm)", 150.0)
	wingThickness := readFloat(reader, "Wing Thickness (%)", 12.0)
	dihedral := readFloat(reader, "Dihedral Angle (degrees)", 5.0)

	fuseLength := readFloat(reader, "Fuselage Length (mm)", 1000.0)
	fuseWidth := readFloat(reader, "Fuselage Width (mm)", 50.0)
	fuseHeight := readFloat(reader, "Fuselage Height (mm)", 60.0)

	materialName := readString(reader, "Material (PLA/PETG/Nylon/CF-Nylon)", "PLA")
	material := materials.GetMaterialByName(materialName)

	estimatedWeight := readFloat(reader, "Estimated Weight (g)", 150.0)

	// Create components
	wing := &fixedwing.Wing{
		Span:      wingSpan,
		Chord:     wingChord,
		Thickness: wingThickness,
		Dihedral:  dihedral,
		Airfoil:   "Clark-Y",
	}

	fuselage := fixedwing.NewFuselage(fuseLength, fuseWidth, fuseHeight)
	tail := fixedwing.NewTail(wing.Span, wing.Chord)

	// Display results
	fmt.Println("\n═══════════════════════════════════════════════════")
	fmt.Println("         Glider Design Summary")
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println()

	fmt.Println("Wing Specifications:")
	fmt.Printf("  Span: %.1f mm\n", wing.Span)
	fmt.Printf("  Chord: %.1f mm\n", wing.Chord)
	fmt.Printf("  Area: %.1f cm²\n", wing.AreaCm2())
	fmt.Printf("  Aspect Ratio: %.2f\n", wing.AspectRatio())
	fmt.Printf("  Glide Ratio: %.1f:1\n", wing.GlideRatio())
	fmt.Println()

	fmt.Println("Performance:")
	fmt.Printf("  Stall Speed: %.1f m/s\n", wing.StallSpeed(estimatedWeight))
	fmt.Printf("  Wing Loading: %.1f g/dm²\n", wing.WingLoading(estimatedWeight))
	fmt.Printf("  Best Glide Speed: %.1f m/s (estimated)\n", wing.StallSpeed(estimatedWeight)*1.4)
	fmt.Println()

	fmt.Println("Build:")
	fmt.Printf("  Material: %s\n", material.Name)
	fmt.Printf("  Density: %.2f g/cm³\n", material.Density)
	fmt.Println()

	// Compile complete design
	design := map[string]interface{}{
		"type":     "Glider",
		"wing":     wing.Specifications(),
		"fuselage": fuselage.Specifications(),
		"tail":     tail.Specifications(),
		"build": map[string]interface{}{
			"material": material.Name,
			"density":  material.Density,
		},
		"performance": map[string]interface{}{
			"weight_g":           estimatedWeight,
			"stall_speed_m/s":    wing.StallSpeed(estimatedWeight),
			"wing_loading_g/dm2": wing.WingLoading(estimatedWeight),
			"glide_ratio":        wing.GlideRatio(),
		},
	}

	// Save to file
	saveDesign(design, "glider_design")
}

func readFloat(reader *bufio.Reader, prompt string, defaultValue float64) float64 {
	fmt.Printf("%s [%.1f]: ", prompt, defaultValue)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultValue
	}

	value, err := strconv.ParseFloat(input, 64)
	if err != nil {
		fmt.Printf("Invalid input, using default: %.1f\n", defaultValue)
		return defaultValue
	}

	return value
}

func readString(reader *bufio.Reader, prompt string, defaultValue string) string {
	fmt.Printf("%s [%s]: ", prompt, defaultValue)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultValue
	}

	return input
}

func saveDesign(design map[string]interface{}, filename string) {
	jsonData, err := json.MarshalIndent(design, "", "  ")
	if err != nil {
		fmt.Printf("Error creating JSON: %v\n", err)
		return
	}

	outputDir := "../../output"
	os.MkdirAll(outputDir, 0755)
	outputFile := fmt.Sprintf("%s/%s.json", outputDir, filename)

	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		fmt.Printf("Error saving design: %v\n", err)
		return
	}

	fmt.Printf("✓ Design saved to: %s\n", outputFile)
}
