package main

import (
	"fmt"
	"time"
	
	"github.com/smaruf/go-lang-study/src/tiny/sr71sim/avionics"
	"github.com/smaruf/go-lang-study/src/tiny/sr71sim/engine"
	"github.com/smaruf/go-lang-study/src/tiny/sr71sim/flying"
	"github.com/smaruf/go-lang-study/src/tiny/sr71sim/fueling"
)

// Location represents a named location with coordinates
type Location struct {
	Name      string
	Latitude  float64
	Longitude float64
}

func main() {
	fmt.Println("=== SR-71 Blackbird World Tour Flight Simulation ===")
	fmt.Println("Route: Florida → Atlantic → Moscow → Asia → Pacific → Virginia")
	fmt.Println()
	
	// Define waypoints for the world tour
	locations := []Location{
		{"Florida (Start)", 28.5383, -81.3792},           // Kennedy Space Center, Florida
		{"Mid-Atlantic", 40.0, -40.0},                    // Over Atlantic Ocean
		{"Moscow, Russia", 55.7558, 37.6173},             // Moscow
		{"Siberia", 60.0, 100.0},                         // Over Siberia
		{"Pacific Ocean", 35.0, -160.0},                  // Mid-Pacific
		{"California Coast", 36.7783, -119.4179},         // California
		{"Virginia (End)", 37.4316, -78.6569},            // Virginia
	}
	
	// Initialize aircraft and systems
	aircraft := flying.NewSR71()
	eng := engine.New()
	avio := avionics.New()
	fuel := fueling.NewFuelSystem()
	
	// Set up waypoints
	waypoints := make([]avionics.Coordinates, len(locations))
	for i, loc := range locations {
		waypoints[i] = avionics.Coordinates{Latitude: loc.Latitude, Longitude: loc.Longitude}
	}
	avio.SetWaypoints(waypoints)
	
	// Start at Florida
	fmt.Println("📍 Starting Position: Florida")
	avio.SetPosition(locations[0].Latitude, locations[0].Longitude)
	fmt.Printf("   Coordinates: %.4f°N, %.4f°W\n", locations[0].Latitude, -locations[0].Longitude)
	fmt.Println()
	
	// Pre-flight status
	fmt.Println("Pre-Flight Status:")
	fmt.Println(fuel.GetStatus())
	fmt.Println("Navigation System: GPS")
	fmt.Println()
	
	// Takeoff sequence
	fmt.Println("=== Takeoff from Florida ===")
	fmt.Println("Rolling down runway...")
	time.Sleep(500 * time.Millisecond)
	
	aircraft.Accelerate(200)
	eng.SetSpeed(0.3) // Mach 0.3
	
	fmt.Println("Wheels up! Climbing...")
	aircraft.FlyAtHeight(10000)
	avio.SetAltitude(10000)
	avio.SetSpeed(200)
	time.Sleep(500 * time.Millisecond)
	
	// Climb to cruise altitude
	fmt.Println("\n=== Climbing to Cruise Altitude ===")
	for alt := 20000; alt <= 80000; alt += 20000 {
		aircraft.FlyAtHeight(alt)
		avio.SetAltitude(float64(alt))
		
		speed := 500 + (alt / 200)
		aircraft.Velocity = speed
		avio.SetSpeed(float64(speed))
		mach := float64(speed) / 767.0
		eng.SetSpeed(mach)
		
		fuel.UpdateEngineType(mach)
		fuel.ConsumeFuel(1 * time.Minute)
		
		fmt.Printf("Altitude: %d ft | Speed: %d mph (Mach %.2f) | Engine: %s\n",
			alt, speed, mach, eng.Mode())
		
		time.Sleep(300 * time.Millisecond)
	}
	
	// Cruise phase - Navigate through waypoints
	fmt.Println("\n=== Cruising at Mach 3+ - World Tour Navigation ===")
	aircraft.AdjustVelocityForMission(flying.MissionReconnaissance)
	mach := aircraft.GetMachNumber()
	eng.SetSpeed(mach)
	eng.SetAltitude(80000)
	avio.SetSpeed(float64(aircraft.Velocity))
	avio.EnableAutopilot()
	avio.Update()
	
	fmt.Printf("Cruise Speed: %d mph (Mach %.2f)\n", aircraft.Velocity, mach)
	fmt.Printf("External Heat: %.0f°F\n", avio.GetState().ExternalHeat)
	fmt.Printf("Cabin Pressure: %.1f psi\n", avio.GetState().CabinPressure)
	fmt.Println()
	
	// Navigate through each waypoint
	for wpIdx := 1; wpIdx < len(locations); wpIdx++ {
		targetLoc := locations[wpIdx]
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("Waypoint %d: %s\n", wpIdx, targetLoc.Name)
		fmt.Printf("Target: %.4f°, %.4f°\n", targetLoc.Latitude, targetLoc.Longitude)
		
		currentPos := avio.GetPosition()
		distance := avionics.CalculateDistance(currentPos, 
			avionics.Coordinates{Latitude: targetLoc.Latitude, Longitude: targetLoc.Longitude})
		fmt.Printf("Distance: %.0f nm\n", distance)
		
		// Environmental variations based on location
		if wpIdx == 1 {
			// Over Atlantic - turbulence, switch to INS
			fmt.Println("\n⚠️  Entering Atlantic Ocean - High Winds Detected")
			fmt.Println("Switching navigation: GPS → INS (Inertial)")
			avio.SwitchNavigationSystem("INS")
			fmt.Println("Engine performance: Nominal in maritime environment")
		} else if wpIdx == 2 {
			// Approaching Moscow - switching back to GPS
			fmt.Println("\n📡 Approaching Moscow Airspace")
			fmt.Println("Switching navigation: INS → GPS")
			avio.SwitchNavigationSystem("GPS")
			fmt.Println("Temperature variation: -40°F ground temp, 600°F exterior")
		} else if wpIdx == 3 {
			// Over Siberia - extreme cold
			fmt.Println("\n❄️  Over Siberia - Extreme Cold Environment")
			fmt.Println("Ground temperature: -60°F")
			fmt.Println("Fuel system: Compensating for thermal differential")
			fuel.ConsumeFuel(5 * time.Minute) // Extra fuel burn in extreme conditions
		} else if wpIdx == 4 {
			// Pacific Ocean - switching to INS again
			fmt.Println("\n🌊 Entering Pacific Ocean")
			fmt.Println("Switching navigation: GPS → INS")
			avio.SwitchNavigationSystem("INS")
			fmt.Println("No ground references - Pure inertial navigation")
		} else if wpIdx == 5 {
			// California coast - back to GPS
			fmt.Println("\n🏔️  Approaching California Coast")
			fmt.Println("Switching navigation: INS → GPS")
			avio.SwitchNavigationSystem("GPS")
		}
		
		// Simulate flight to waypoint
		etaMinutes := distance / (float64(aircraft.Velocity) * 0.868976 / 60.0)
		fmt.Printf("ETA: %.1f minutes\n", etaMinutes)
		fmt.Println()
		
		// Navigate in segments
		segmentTime := 2.0 // minutes per update
		segmentCount := int(etaMinutes / segmentTime)
		if segmentCount < 1 {
			segmentCount = 1
		}
		
		for seg := 0; seg < segmentCount; seg++ {
			reached := avio.NavigateToWaypoint(float64(aircraft.Velocity), segmentTime)
			
			currentPos = avio.GetPosition()
			remaining := avio.DistanceToWaypoint()
			
			fmt.Printf("  Position: %.4f°, %.4f° | Remaining: %.0f nm\n",
				currentPos.Latitude, currentPos.Longitude, remaining)
			
			// Consume fuel
			fuel.ConsumeFuel(time.Duration(segmentTime) * time.Minute)
			
			time.Sleep(200 * time.Millisecond)
			
			if reached {
				break
			}
		}
		
		// Ensure we reach the waypoint
		avio.SetPosition(targetLoc.Latitude, targetLoc.Longitude)
		avio.NextWaypoint()
		
		fmt.Printf("\n✓ Reached: %s\n", targetLoc.Name)
		fmt.Printf("  Current Fuel: %.1f%%\n", fuel.GetFuelLevel())
		
		// Refuel if needed (simulate aerial refueling)
		if fuel.NeedsRefueling() && wpIdx < len(locations)-1 {
			fmt.Println("\n⚠️  Fuel level low - Aerial refueling required")
			fuel.StartRefueling()
			
			fmt.Println("Rendezvous with tanker aircraft...")
			aircraft.Decelerate(700) // Slow down for refueling
			avio.SetSpeed(float64(aircraft.Velocity))
			time.Sleep(500 * time.Millisecond)
			
			fmt.Println("Connected to refueling boom")
			for i := 0; i < 3; i++ {
				fuel.Refuel(2 * time.Minute)
				fmt.Printf("Refueling... %.1f%%\n", fuel.GetFuelLevel())
				time.Sleep(300 * time.Millisecond)
			}
			
			fuel.StopRefueling()
			aircraft.Accelerate(700) // Resume cruise speed
			avio.SetSpeed(float64(aircraft.Velocity))
			fmt.Println("Refueling complete - Resuming cruise speed")
		}
		
		fmt.Println()
		time.Sleep(300 * time.Millisecond)
	}
	
	// Descent and landing in Virginia
	fmt.Println("\n=== Beginning Descent to Virginia ===")
	avio.DisableAutopilot()
	
	for alt := 70000; alt >= 10000; alt -= 15000 {
		aircraft.FlyAtHeight(alt)
		avio.SetAltitude(float64(alt))
		
		speed := 300 + (alt / 200)
		aircraft.Velocity = speed
		avio.SetSpeed(float64(speed))
		mach := float64(speed) / 767.0
		eng.SetSpeed(mach)
		
		fuel.UpdateEngineType(mach)
		
		fmt.Printf("Descending: %d ft | Speed: %d mph (Mach %.2f)\n",
			alt, speed, mach)
		
		time.Sleep(300 * time.Millisecond)
	}
	
	// Final approach
	fmt.Println("\n=== Final Approach - Virginia ===")
	aircraft.FlyAtHeight(5000)
	aircraft.Decelerate(1000)
	avio.SetSpeed(float64(aircraft.Velocity))
	
	fmt.Println("Landing gear down")
	fmt.Println("On final approach...")
	time.Sleep(500 * time.Millisecond)
	
	aircraft.FlyAtHeight(0)
	aircraft.Decelerate(aircraft.Velocity)
	
	fmt.Println("✓ Touchdown in Virginia!")
	fmt.Println("Deploying drag chute...")
	fmt.Println("Aircraft stopped")
	
	// Post-flight summary
	fmt.Println("\n=== World Tour Flight Complete! ===")
	finalPos := avio.GetPosition()
	fmt.Printf("Final Position: %.4f°, %.4f°\n", finalPos.Latitude, finalPos.Longitude)
	fmt.Printf("Total Distance: ~%.0f nautical miles\n", 
		avionics.CalculateDistance(
			avionics.Coordinates{Latitude: locations[0].Latitude, Longitude: locations[0].Longitude},
			finalPos) * float64(len(locations)-1) / 2)
	fmt.Println()
	fmt.Println(fuel.GetStatus())
	fmt.Println()
	fmt.Println("Route Summary:")
	for i, loc := range locations {
		if i == 0 {
			fmt.Printf("  ✈️  %s (Takeoff)\n", loc.Name)
		} else if i == len(locations)-1 {
			fmt.Printf("  🛬 %s (Landing)\n", loc.Name)
		} else {
			fmt.Printf("  📍 %s\n", loc.Name)
		}
	}
	fmt.Println("\nFlight complete! 🌍✨")
}
