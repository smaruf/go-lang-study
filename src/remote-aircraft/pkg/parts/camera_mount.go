package parts

import (
	"github.com/smaruf/go-lang-study/src/remote-aircraft/pkg/materials"
)

// CameraMount represents an FPV camera mount
type CameraMount struct {
	CameraWidth   float64            // Camera width in mm
	CameraHeight  float64            // Camera height in mm
	TiltAngle     float64            // Camera tilt angle in degrees
	WallThickness float64            // Mount wall thickness in mm
	Material      materials.Material // Material type
}

// NewCameraMount creates a camera mount for standard FPV cameras
func NewCameraMount(tiltAngle float64, material materials.Material) *CameraMount {
	return &CameraMount{
		CameraWidth:   28.0, // Standard micro camera width
		CameraHeight:  28.0, // Standard micro camera height
		TiltAngle:     tiltAngle,
		WallThickness: 2.0,
		Material:      material,
	}
}

// Volume calculates the volume in mm³
func (c *CameraMount) Volume() float64 {
	// Mount dimensions
	mountWidth := c.CameraWidth + 2*c.WallThickness + 2
	mountHeight := c.CameraHeight + 2*c.WallThickness + 2
	mountDepth := 15.0 // Typical depth
	
	// Calculate total volume
	totalVolume := mountWidth * mountHeight * mountDepth
	
	// Subtract camera cavity
	cavityDepth := 10.0
	cavityVolume := c.CameraWidth * c.CameraHeight * cavityDepth
	
	return totalVolume - cavityVolume
}

// Weight calculates the weight in grams
func (c *CameraMount) Weight() float64 {
	volumeCm3 := c.Volume() / 1000.0 // Convert mm³ to cm³
	return volumeCm3 * c.Material.Density
}

// EffectiveTiltAngle returns the effective tilt angle considering mounting
func (c *CameraMount) EffectiveTiltAngle() float64 {
	return c.TiltAngle
}

// Specifications returns a map of specifications
func (c *CameraMount) Specifications() map[string]interface{} {
	mountWidth := c.CameraWidth + 2*c.WallThickness + 2
	mountHeight := c.CameraHeight + 2*c.WallThickness + 2
	
	return map[string]interface{}{
		"type":           "Camera Mount",
		"camera_width":   c.CameraWidth,
		"camera_height":  c.CameraHeight,
		"tilt_angle":     c.TiltAngle,
		"mount_width":    mountWidth,
		"mount_height":   mountHeight,
		"wall_thickness": c.WallThickness,
		"material":       c.Material.Name,
		"weight_g":       c.Weight(),
		"volume_cm3":     c.Volume() / 1000.0,
	}
}
