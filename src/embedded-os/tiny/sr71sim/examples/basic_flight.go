package main

import (
	"fmt"
	"time"
	
	"github.com/smaruf/go-lang-study/src/embedded-os/tiny/sr71sim/avionics"
	"github.com/smaruf/go-lang-study/src/embedded-os/tiny/sr71sim/engine"
	"github.com/smaruf/go-lang-study/src/embedded-os/tiny/sr71sim/flying"
	"github.com/smaruf/go-lang-study/src/embedded-os/tiny/sr71sim/fueling"
)

func main() {
	fmt.Println("=== SR-71 Blackbird Basic Flight Simulation ===")
	fmt.Println()
	
	// Initialize aircraft
	aircraft := flying.NewSR71()
	eng := engine.New()
	avio := avionics.New()
	fuel := fueling.NewFuelSystem()
	
	// Pre-flight status
	fmt.Println("Pre-Flight Status:")
	fmt.Println(fuel.GetStatus())
	fmt.Println()
	
	// Takeoff sequence
	fmt.Println("=== Takeoff Sequence ===")
	fmt.Println("Rolling down runway...")
	time.Sleep(1 * time.Second)
	
	aircraft.Accelerate(200)
	eng.SetSpeed(0.3) // Mach 0.3
	
	fmt.Println("Wheels up! Climbing...")
	aircraft.FlyAtHeight(10000)
	avio.SetAltitude(10000)
	avio.SetSpeed(200)
	time.Sleep(1 * time.Second)
	
	// Climb to cruise altitude
	fmt.Println("\n=== Climbing to Cruise Altitude ===")
	for alt := 20000; alt <= 80000; alt += 10000 {
		aircraft.FlyAtHeight(alt)
		avio.SetAltitude(float64(alt))
		
		// Increase speed with altitude
		speed := 500 + (alt / 200)
		aircraft.Velocity = speed
		avio.SetSpeed(float64(speed))
		mach := float64(speed) / 767.0
		eng.SetSpeed(mach)
		
		// Update fuel system
		fuel.UpdateEngineType(mach)
		fuel.ConsumeFuel(2 * time.Minute)
		
		fmt.Printf("Altitude: %d ft | Speed: %d mph (Mach %.2f) | Engine: %s | Fuel: %.1f%%\n",
			alt, speed, mach, eng.Mode(), fuel.GetFuelLevel())
		
		time.Sleep(500 * time.Millisecond)
	}
	
	// Cruise phase
	fmt.Println("\n=== Cruising at Mach 3+ ===")
	aircraft.AdjustVelocityForMission(flying.MissionReconnaissance)
	mach := aircraft.GetMachNumber()
	eng.SetSpeed(mach)
	eng.SetAltitude(80000)
	avio.SetSpeed(float64(aircraft.Velocity))
	avio.Update()
	
	fmt.Println("\nCruise Status:")
	fmt.Println(aircraft.GetStatus())
	fmt.Println("\nEngine:", eng.GetState().Mode)
	fmt.Printf("External Heat: %.0f°F\n", avio.GetState().ExternalHeat)
	fmt.Printf("Cabin Pressure: %.1f psi\n", avio.GetState().CabinPressure)
	fmt.Println()
	
	// Simulate cruise for 10 minutes
	fmt.Println("Cruising for 10 minutes...")
	fuel.ConsumeFuel(10 * time.Minute)
	fmt.Println(fuel.GetStatus())
	fmt.Println()
	
	// Check if refueling needed
	if fuel.NeedsRefueling() {
		fmt.Println("⚠️  Fuel level low - initiating aerial refueling")
		fuel.StartRefueling()
		
		fmt.Println("Contacting tanker aircraft...")
		time.Sleep(1 * time.Second)
		fmt.Println("Connected to refueling boom")
		
		for i := 0; i < 5; i++ {
			fuel.Refuel(1 * time.Minute)
			fmt.Printf("Refueling... Fuel level: %.1f%%\n", fuel.GetFuelLevel())
			time.Sleep(500 * time.Millisecond)
		}
		
		fuel.StopRefueling()
		fmt.Println()
	}
	
	// Descent and landing
	fmt.Println("=== Beginning Descent ===")
	for alt := 70000; alt >= 10000; alt -= 10000 {
		aircraft.FlyAtHeight(alt)
		avio.SetAltitude(float64(alt))
		
		// Decrease speed with altitude
		speed := 300 + (alt / 200)
		aircraft.Velocity = speed
		avio.SetSpeed(float64(speed))
		mach := float64(speed) / 767.0
		eng.SetSpeed(mach)
		
		fuel.UpdateEngineType(mach)
		
		fmt.Printf("Descending: %d ft | Speed: %d mph (Mach %.2f) | Engine: %s\n",
			alt, speed, mach, eng.Mode())
		
		time.Sleep(500 * time.Millisecond)
	}
	
	// Final approach
	fmt.Println("\n=== Final Approach ===")
	aircraft.FlyAtHeight(5000)
	aircraft.Decelerate(1000)
	
	fmt.Println("Landing gear down")
	fmt.Println("On final approach...")
	time.Sleep(1 * time.Second)
	
	aircraft.FlyAtHeight(0)
	aircraft.Decelerate(aircraft.Velocity)
	
	fmt.Println("✓ Touchdown!")
	fmt.Println("Deploying drag chute...")
	fmt.Println("Aircraft stopped")
	
	// Post-flight status
	fmt.Println("\n=== Post-Flight Summary ===")
	fmt.Println(aircraft.GetStatus())
	fmt.Println()
	fmt.Println(fuel.GetStatus())
	fmt.Println()
	fmt.Println("Flight complete! 🎉")
}
