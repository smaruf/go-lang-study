package main

import (
	"machine"
	"time"
)

// Button Input Example for Raspberry Pi Pico
// Hardware:
//   - Button connected to GPIO 15
//   - 10kΩ pull-up resistor to 3.3V
//   - Button to GND
//   - LED on GPIO 25 (built-in)

func main() {
	// Configure LED
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// Configure button with pull-up
	button := machine.GP15
	button.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

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
			} else {
				led.Low()
			}
			
			// Debounce delay
			time.Sleep(50 * time.Millisecond)
		}

		prevButtonState = currentButtonState
		time.Sleep(10 * time.Millisecond)
	}
}
