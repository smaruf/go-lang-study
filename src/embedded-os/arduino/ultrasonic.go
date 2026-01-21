package main

import (
	"machine"
	"time"
)

// Ultrasonic Distance Sensor Example for Arduino
// Hardware:
//   - HC-SR04 Ultrasonic sensor
//   - Trigger pin to D7
//   - Echo pin to D8
//   - VCC to 5V
//   - GND to GND

const (
	SOUND_SPEED = 0.0343 // Speed of sound in cm/microsecond
)

var (
	trigPin = machine.D7
	echoPin = machine.D8
)

func main() {
	// Configure pins
	trigPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	echoPin.Configure(machine.PinConfig{Mode: machine.PinInput})

	println("Ultrasonic Distance Sensor")
	println("Measuring distance...")

	for {
		distance := measureDistance()
		
		if distance > 0 && distance < 400 {
			println("Distance:", int(distance), "cm")
		} else {
			println("Out of range")
		}
		
		time.Sleep(500 * time.Millisecond)
	}
}

func measureDistance() float32 {
	// Send 10µs pulse
	trigPin.Low()
	time.Sleep(2 * time.Microsecond)
	trigPin.High()
	time.Sleep(10 * time.Microsecond)
	trigPin.Low()

	// Wait for echo pin to go high
	timeout := time.Now().Add(30 * time.Millisecond)
	for echoPin.Get() == false {
		if time.Now().After(timeout) {
			return -1 // Timeout
		}
	}

	// Measure pulse width
	startTime := time.Now()
	
	timeout = time.Now().Add(30 * time.Millisecond)
	for echoPin.Get() == true {
		if time.Now().After(timeout) {
			return -1 // Timeout
		}
	}
	
	duration := time.Since(startTime)

	// Calculate distance
	// Distance = (Time × Speed of Sound) / 2
	distance := float32(duration.Microseconds()) * SOUND_SPEED / 2.0

	return distance
}
