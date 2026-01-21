package main

import (
	"machine"
	"time"
)

// Button Input Example for Arduino
// Hardware:
//   - Button connected to Pin 2
//   - 10kΩ pull-up resistor to 5V
//   - Button to GND
//   - LED on Pin 13 (built-in)

func main() {
	// Configure LED
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// Configure button with pull-up
	button := machine.D2
	button.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	println("Button Input Example Started")

	// Track LED state
	ledState := false

	// Previous button state for edge detection
	prevButtonState := button.Get()

	for {
		// Read current button state
		currentButtonState := button.Get()

		// Detect button press (transition from high to low)
		if prevButtonState && !currentButtonState {
			// Button pressed - toggle LED
			ledState = !ledState
			if ledState {
				led.High()
				println("LED ON")
			} else {
				led.Low()
				println("LED OFF")
			}
			
			// Debounce delay
			time.Sleep(50 * time.Millisecond)
		}

		prevButtonState = currentButtonState
		time.Sleep(10 * time.Millisecond)
	}
}
