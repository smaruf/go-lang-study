package main

import (
	"machine"
	"time"
)

// Blinky - Classic LED blink example for Raspberry Pi Pico
// Hardware: Uses built-in LED on GPIO 25
// No external components required

func main() {
	// Configure the built-in LED pin
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// Blink forever
	for {
		led.High() // Turn LED on
		time.Sleep(500 * time.Millisecond)
		
		led.Low() // Turn LED off
		time.Sleep(500 * time.Millisecond)
	}
}
