# Remote Aircraft: FPV Drone & Fixed-Wing Design System (Go Edition)

**Professional parametric aircraft design system for designing, building, and flying FPV drones and fixed-wing aircraft implemented in Go.**

This is a complete, production-ready repository combining:
- 🐹 Go-based parametric design calculations
- 🎓 Hands-on multirotor and fixed-wing design tools
- 📐 Engineering analysis tools (structural & aerodynamic)
- 🖨️ 3D printing calculations and specifications
- ✈️ Flight-ready designs for multirotors and fixed-wing UAVs

---

## 🚀 Quick Start

### 1. Installation

```bash
# Clone or navigate to repository
cd src/remote-aircraft

# Install dependencies
go mod tidy
```

### 2. Use the CLI Designer

```bash
# Launch the Aircraft Designer CLI
go run cmd/designer/main.go

# Or build and run
go build -o aircraft-designer cmd/designer/main.go
./aircraft-designer
```

This command-line interface allows you to:
- Design Fixed Wing Aircraft or Gliders
- Enter custom parameters
- Generate design specifications
- Calculate performance metrics

### 3. Generate Design Programmatically

```bash
# Run analysis examples
go run examples/weight_calc/main.go
go run examples/stress_analysis/main.go
go run examples/fixed_wing_analysis/main.go

# Generate parts specifications
go run examples/generate_motor_mounts/main.go
go run examples/generate_fixed_wing/main.go
```

---

## 📦 Repository Structure

```
remote-aircraft/
├── README.md                    # This file
├── go.mod                       # Go module definition
├── cmd/                         # Command-line applications
│   └── designer/                # Interactive aircraft designer
│       └── main.go
│
├── pkg/                         # Reusable packages
│   ├── materials/               # Material properties database
│   │   └── materials.go
│   ├── parts/                   # Parametric component designs
│   │   ├── motor_mount.go       # Motor mounting specifications
│   │   ├── arm.go               # Quadcopter arms
│   │   ├── camera_mount.go      # FPV camera mounts
│   │   └── battery_tray.go      # Battery holders
│   ├── frames/                  # Complete frame assemblies
│   │   └── quad_frame.go        # Quadcopter frame generator
│   ├── analysis/                # Engineering calculations
│   │   ├── weight.go            # Weight calculations
│   │   ├── cg.go                # Center of gravity
│   │   └── stress.go            # Stress analysis
│   └── fixedwing/               # Fixed-Wing Aircraft Design
│       ├── wing.go              # Parametric wing design
│       ├── spar.go              # Spar design & load calcs
│       ├── fuselage.go          # Fuselage sections
│       ├── tail.go              # Tail components
│       └── loads.go             # Aerodynamic loads
│
├── examples/                    # Runnable examples
│   ├── weight_calc/             # Weight & CG calculator
│   ├── stress_analysis/         # Stress analysis
│   ├── generate_motor_mounts/   # Custom motor mounts
│   ├── fixed_wing_analysis/     # Fixed-wing load analysis
│   └── generate_fixed_wing/     # Generate fixed-wing specs
│
└── output/                      # Generated specifications
    └── *.json
```

---

## ✨ Features

### 🖥️ CLI Aircraft Designer
- **Interactive Design**: User-friendly command-line interface
- **Fixed Wing Aircraft**: Complete parametric design with motor
- **Gliders**: Optimized for unpowered flight performance
- **Material Selection**: Choose from PLA, PETG, Nylon, or CF-Nylon
- **Design Summary**: Automatic performance calculations and recommendations
- **JSON Output**: Export designs in JSON format for further processing

### Parametric Design Calculations

#### Multirotor Components
- **Motor Mounts**: Customizable for any motor size (1507 to 2810+)
- **Quadcopter Arms**: Various lengths and cross-sections
- **Camera Mounts**: Adjustable tilt angles
- **Battery Trays**: Sized for different battery capacities
- **Complete Frames**: Full quadcopter assemblies

#### Fixed-Wing Components ✈️
- **Wing Design**: Parametric airfoil calculations (Clark-Y, symmetric)
- **Fuselage Sections**: Modular design specifications
- **Wing Mount Specifications**: Reinforced connection components
- **Tail Components**: Horizontal & vertical stabilizers
- **Tail Boom Mounting**: Structural calculations

#### Material Database
Pre-configured materials with density values:
- **PLA**: 1.24 g/cm³ - Easy printing, good strength
- **PETG**: 1.27 g/cm³ - Weather resistant, durable
- **Nylon**: 1.15 g/cm³ - High strength, flexible
- **CF-Nylon**: 1.20 g/cm³ - Carbon fiber reinforced, maximum strength

### Engineering Analysis Tools

#### Weight & Balance
- Component weight calculations
- Center of gravity (CG) computation
- Balance point determination
- Weight distribution analysis

#### Structural Analysis
- Beam stress calculations
- Factor of safety determination
- Load capacity estimates
- Material strength verification

#### Aerodynamic Analysis (Fixed-Wing)
- Lift calculations based on wing area and airspeed
- Drag estimation
- Glide ratio prediction
- Stall speed calculation
- Required thrust estimation

---

## 🎯 Usage Examples

### Basic Motor Mount Design

```go
package main

import (
    "fmt"
    "github.com/smaruf/go-lang-study/src/remote-aircraft/pkg/parts"
    "github.com/smaruf/go-lang-study/src/remote-aircraft/pkg/materials"
)

func main() {
    mount := parts.MotorMount{
        MotorDiameter: 28.0,  // mm
        HoleSpacing:   19.0,  // mm
        Thickness:     3.0,   // mm
        Material:      materials.PETG,
    }
    
    weight := mount.CalculateWeight()
    fmt.Printf("Motor mount weight: %.2f g\n", weight)
    
    stress := mount.CalculateStress(500.0) // 500g thrust
    fmt.Printf("Maximum stress: %.2f MPa\n", stress)
}
```

### Fixed-Wing Aircraft Analysis

```go
package main

import (
    "fmt"
    "github.com/smaruf/go-lang-study/src/remote-aircraft/pkg/fixedwing"
)

func main() {
    wing := fixedwing.Wing{
        Span:      1000.0, // mm
        Chord:     200.0,  // mm
        Thickness: 12.0,   // %
        Dihedral:  3.0,    // degrees
    }
    
    // Calculate wing area and loading
    area := wing.Area()
    fmt.Printf("Wing area: %.2f cm²\n", area)
    
    // Estimate lift at given airspeed
    airspeed := 15.0 // m/s
    lift := wing.CalculateLift(airspeed)
    fmt.Printf("Lift at %.1f m/s: %.2f g\n", airspeed, lift)
}
```

---

## 🧮 Engineering Calculations

### Supported Calculations
- **Weight**: Component mass based on material density and volume
- **Center of Gravity**: 3D CG position for complete aircraft
- **Stress**: Beam bending stress under load
- **Lift**: Wing lift based on area, airspeed, and coefficient
- **Drag**: Drag estimation for performance prediction
- **Glide Ratio**: L/D ratio for gliders and powered aircraft
- **Stall Speed**: Minimum flying speed calculation

---

## 🏗️ Building & Testing

### Build the CLI Designer

```bash
go build -o aircraft-designer cmd/designer/main.go
```

### Run Tests

```bash
go test ./...
```

### Run with Coverage

```bash
go test -cover ./...
```

---

## 📚 Design Guidelines

### Multirotor Design
- **Motor Selection**: Match motor KV to battery voltage and prop size
- **Arm Length**: 120-300mm typical for mini to racing quads
- **Weight Distribution**: Keep CG near geometric center
- **Prop Clearance**: Minimum 10mm ground clearance

### Fixed-Wing Design
- **Wing Loading**: 20-40 g/dm² for trainers, 40-80 g/dm² for sport
- **Aspect Ratio**: 6-10 for efficiency, 4-6 for aerobatics
- **CG Position**: 25-33% of mean aerodynamic chord
- **Tail Volume**: Follow standard coefficients for stability

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit pull requests or open issues for bugs and feature requests.

---

## 📄 License

This project is part of the go-lang-study repository and follows the same license.

---

**Start Designing**: Run the CLI designer and create your custom aircraft! 🚀✈️

For questions or suggestions, please open an issue or contribute to the discussion.
