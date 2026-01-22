package parts

import (
	"github.com/smaruf/go-lang-study/src/remote-aircraft/pkg/materials"
)

// BatteryTray represents a battery holder/tray
type BatteryTray struct {
	BatteryLength float64            // Battery length in mm
	BatteryWidth  float64            // Battery width in mm
	BatteryHeight float64            // Battery height in mm
	WallThickness float64            // Tray wall thickness in mm
	BaseThickness float64            // Tray base thickness in mm
	Material      materials.Material // Material type
}

// NewBatteryTray creates a battery tray for standard battery sizes
func NewBatteryTray(batteryLength, batteryWidth, batteryHeight float64, material materials.Material) *BatteryTray {
	return &BatteryTray{
		BatteryLength: batteryLength,
		BatteryWidth:  batteryWidth,
		BatteryHeight: batteryHeight,
		WallThickness: 2.0,
		BaseThickness: 2.5,
		Material:      material,
	}
}

// Volume calculates the volume in mm³
func (b *BatteryTray) Volume() float64 {
	// Calculate outer dimensions
	outerLength := b.BatteryLength + 2*b.WallThickness + 4 // Add clearance
	outerWidth := b.BatteryWidth + 2*b.WallThickness + 4
	outerHeight := b.BatteryHeight * 0.7 // Tray goes 70% up the battery
	
	// Total volume including walls and base
	totalVolume := outerLength * outerWidth * (outerHeight + b.BaseThickness)
	
	// Subtract inner cavity
	innerVolume := b.BatteryLength * b.BatteryWidth * outerHeight
	
	return totalVolume - innerVolume
}

// Weight calculates the weight in grams
func (b *BatteryTray) Weight() float64 {
	volumeCm3 := b.Volume() / 1000.0 // Convert mm³ to cm³
	return volumeCm3 * b.Material.Density
}

// Specifications returns a map of specifications
func (b *BatteryTray) Specifications() map[string]interface{} {
	outerLength := b.BatteryLength + 2*b.WallThickness + 4
	outerWidth := b.BatteryWidth + 2*b.WallThickness + 4
	outerHeight := b.BatteryHeight*0.7 + b.BaseThickness
	
	return map[string]interface{}{
		"type":            "Battery Tray",
		"battery_length":  b.BatteryLength,
		"battery_width":   b.BatteryWidth,
		"battery_height":  b.BatteryHeight,
		"outer_length":    outerLength,
		"outer_width":     outerWidth,
		"outer_height":    outerHeight,
		"wall_thickness":  b.WallThickness,
		"base_thickness":  b.BaseThickness,
		"material":        b.Material.Name,
		"weight_g":        b.Weight(),
		"volume_cm3":      b.Volume() / 1000.0,
	}
}
