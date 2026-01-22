# Usage Guide - Remote Aircraft

This guide provides detailed usage examples for the Remote Aircraft design system.

## Table of Contents
- [Quick Start](#quick-start)
- [CLI Designer](#cli-designer)
- [Programmatic Usage](#programmatic-usage)
- [Examples](#examples)
- [Output Files](#output-files)

---

## Quick Start

### Install and Build

```bash
cd src/remote-aircraft
go mod tidy
go build -o aircraft-designer cmd/designer/main.go
```

### Run Examples

```bash
# Weight and center of gravity analysis
go run examples/weight_calc/main.go

# Structural stress analysis
go run examples/stress_analysis/main.go

# Fixed-wing performance analysis
go run examples/fixed_wing_analysis/main.go

# Generate motor mount specifications
go run examples/generate_motor_mounts/main.go

# Generate fixed-wing aircraft designs
go run examples/generate_fixed_wing/main.go
```

---

## CLI Designer

The interactive CLI designer allows you to create custom aircraft designs.

### Launch the Designer

```bash
./aircraft-designer
# or
go run cmd/designer/main.go
```

### Design Fixed-Wing Aircraft

1. Select option 1 for "Fixed Wing Aircraft"
2. Enter parameters (or press Enter for defaults):
   - Wing Span: 1000 mm
   - Wing Chord: 200 mm
   - Wing Thickness: 12%
   - Dihedral Angle: 3°
   - Fuselage dimensions
   - Tail dimensions
   - Motor and propeller specs
   - Material selection
   - Estimated weight

3. Review the design summary with:
   - Wing specifications
   - Performance metrics
   - Propulsion requirements
   - Build materials

4. Design is automatically saved to `output/fixed_wing_design.json`

### Design Glider

1. Select option 2 for "Glider"
2. Enter parameters focused on unpowered flight
3. Review glider-specific performance metrics
4. Design saved to `output/glider_design.json`

---

## Programmatic Usage

### Weight and Balance Analysis

```go
package main

import (
    "fmt"
    "github.com/smaruf/go-lang-study/src/remote-aircraft/pkg/analysis"
)

func main() {
    // Create weight analysis
    wa := analysis.NewWeightAnalysis()
    
    // Add components with positions
    wa.AddComponent("Frame", 68.0, 0, 0, 0)
    wa.AddComponent("Battery", 180.0, 0, 0, -10)
    wa.AddComponent("Motor FL", 28.0, 150, 150, 0)
    
    // Calculate total weight and CG
    totalWeight := wa.TotalWeight()
    cg := wa.CenterOfGravity()
    
    fmt.Printf("Total: %.1fg, CG: (%.1f, %.1f, %.1f)\n",
        totalWeight, cg.X, cg.Y, cg.Z)
}
```

### Stress Analysis

```go
package main

import (
    "fmt"
    "github.com/smaruf/go-lang-study/src/remote-aircraft/pkg/analysis"
)

func main() {
    // Create stress analyzer for PETG
    sa := analysis.NewStressAnalysis("PETG", analysis.PETGYieldStrength)
    
    // Analyze arm under crash load
    results := sa.Analyze(150.0, 12.0, 8.0, 500.0)
    
    fmt.Printf("Stress: %.2f MPa\n", results["applied_stress_MPa"])
    fmt.Printf("Safe: %v\n", results["is_safe"])
}
```

### Wing Design

```go
package main

import (
    "fmt"
    "github.com/smaruf/go-lang-study/src/remote-aircraft/pkg/fixedwing"
)

func main() {
    // Create wing
    wing := fixedwing.NewWing(1000.0, 200.0)
    
    // Calculate performance
    area := wing.AreaCm2()
    lift := wing.CalculateLift(15.0)  // at 15 m/s
    stallSpeed := wing.StallSpeed(250.0)  // for 250g aircraft
    
    fmt.Printf("Area: %.1f cm²\n", area)
    fmt.Printf("Lift: %.1f g\n", lift)
    fmt.Printf("Stall: %.1f m/s\n", stallSpeed)
}
```

### Complete Aircraft

```go
package main

import (
    "fmt"
    "github.com/smaruf/go-lang-study/src/remote-aircraft/pkg/fixedwing"
)

func main() {
    // Create aircraft components
    wing := fixedwing.NewWing(1000.0, 200.0)
    fuselage := fixedwing.NewFuselage(800.0, 60.0, 80.0)
    tail := fixedwing.NewTail(wing.Span, wing.Chord)
    
    // Perform load analysis
    loads := fixedwing.NewLoads(250.0, 15.0, wing)
    analysis := loads.Analyze()
    
    fmt.Printf("Required thrust: %.1f g\n", 
        analysis["required_thrust_g"])
    fmt.Printf("Power required: %.1f W\n", 
        analysis["power_required_W"])
}
```

### Motor Mounts

```go
package main

import (
    "fmt"
    "github.com/smaruf/go-lang-study/src/remote-aircraft/pkg/parts"
    "github.com/smaruf/go-lang-study/src/remote-aircraft/pkg/materials"
)

func main() {
    // Create motor mount for 2207 motor
    mount := parts.NewMotorMount(22.0, materials.PETG)
    
    weight := mount.Weight()
    stress := mount.CalculateStress(500.0)
    
    fmt.Printf("Weight: %.2f g\n", weight)
    fmt.Printf("Stress: %.2f MPa\n", stress)
}
```

### Quadcopter Frame

```go
package main

import (
    "fmt"
    "github.com/smaruf/go-lang-study/src/remote-aircraft/pkg/frames"
    "github.com/smaruf/go-lang-study/src/remote-aircraft/pkg/materials"
)

func main() {
    // Create 150mm frame
    frame := frames.NewQuadFrame(150.0, materials.PETG)
    
    weight := frame.TotalWeight()
    wheelbase := frame.Wheelbase()
    
    fmt.Printf("Frame weight: %.1f g\n", weight)
    fmt.Printf("Wheelbase: %.1f mm\n", wheelbase)
}
```

---

## Examples

### Example 1: Weight Calculator
**File:** `examples/weight_calc/main.go`

Demonstrates:
- Creating a complete quadcopter weight analysis
- Adding components with 3D positions
- Calculating center of gravity
- Checking balance

**Output:** `output/weight_analysis.json`

### Example 2: Stress Analysis
**File:** `examples/stress_analysis/main.go`

Demonstrates:
- Testing different materials (PLA, PETG, Nylon, CF-Nylon)
- Calculating stress under load
- Determining safety factors
- Finding maximum safe loads

**Output:** `output/stress_analysis.json`

### Example 3: Fixed-Wing Analysis
**File:** `examples/fixed_wing_analysis/main.go`

Demonstrates:
- Complete fixed-wing aircraft design
- Aerodynamic calculations
- Performance at multiple airspeeds
- Load analysis

**Output:** `output/fixed_wing_analysis.json`

### Example 4: Motor Mount Generator
**File:** `examples/generate_motor_mounts/main.go`

Demonstrates:
- Generating mounts for common motor sizes
- Weight and stress calculations
- JSON export of specifications

**Output:** `output/motor_mounts.json`

### Example 5: Fixed-Wing Generator
**File:** `examples/generate_fixed_wing/main.go`

Demonstrates:
- Multiple aircraft configurations
- Trainer, sport, and glider designs
- Performance comparison
- Complete specifications export

**Output:** `output/fixed_wing_designs.json`

---

## Output Files

All examples and the CLI designer generate JSON files in the `output/` directory.

### JSON Structure

```json
{
  "type": "Fixed Wing Aircraft",
  "wing": {
    "span_mm": 1000.0,
    "chord_mm": 200.0,
    "area_cm2": 2000.0,
    "aspect_ratio": 5.0
  },
  "performance": {
    "stall_speed_m/s": 3.7,
    "wing_loading_g/dm2": 12.5,
    "required_thrust_g": 112.4
  }
}
```

### Using Output Files

The JSON files can be:
- Imported into other tools
- Used for documentation
- Compared across designs
- Shared with team members

---

## Tips and Best Practices

### Design Tips

1. **Wing Loading**: Keep between 20-40 g/dm² for trainers
2. **Aspect Ratio**: Higher = better glide performance
3. **Dihedral**: 3-5° for stability
4. **Material Selection**: 
   - PLA: Easy to print, good for prototypes
   - PETG: Weather resistant, good for final builds
   - CF-Nylon: Maximum strength for high-stress parts

### Performance Optimization

1. **Minimize Weight**: Every gram counts
2. **Balance CG**: Keep within ±10mm of center
3. **Adequate Thrust**: 2:1 thrust-to-weight ratio recommended
4. **Stall Margin**: Cruise speed should be 2-3x stall speed

### Common Issues

**High stall speed?**
- Increase wing area
- Reduce weight
- Check wing loading

**Excessive stress?**
- Increase wall thickness
- Use stronger material
- Redesign load path

**Poor glide ratio?**
- Increase aspect ratio
- Reduce drag
- Optimize airfoil

---

## Next Steps

1. Try the CLI designer with your own parameters
2. Modify the examples for your specific needs
3. Create custom components by extending the packages
4. Share your designs with the community

For questions or issues, please open an issue on GitHub.
