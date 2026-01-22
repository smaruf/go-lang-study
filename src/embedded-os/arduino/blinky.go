package main

import (
	"machine"
	"time"
)

// Blinky - Classic LED blink example for Arduino
// Hardware: Uses built-in LED on Pin 13
// No external components required

func main() {
	// Configure the built-in LED pin (Pin 13)
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	println("Arduino Blinky Started")

	// Blink forever
	for {
		led.High() // Turn LED on
		time.Sleep(1000 * time.Millisecond)
		
		led.Low() // Turn LED off
		time.Sleep(1000 * time.Millisecond)
	}
}
