package main

import (
	"machine"
	"time"
)

// PWM LED Dimming Example for Raspberry Pi Pico
// Hardware:
//   - External LED connected to GPIO 15 with 220Ω resistor
//   - Or use built-in LED on GPIO 25

func main() {
	// Use GPIO 15 for PWM LED (or machine.LED for built-in)
	ledPin := machine.GP15
	ledPin.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// Configure PWM
	pwm := machine.PWM4 // PWM slice for GP15
	err := pwm.Configure(machine.PWMConfig{
		Period: 16384, // PWM period
	})
	if err != nil {
		println("PWM configuration error:", err.Error())
		return
	}

	// Get PWM channel for the pin
	channel, err := pwm.Channel(ledPin)
	if err != nil {
		println("PWM channel error:", err.Error())
		return
	}

	brightness := uint32(0)
	increasing := true

	for {
		// Set PWM duty cycle (brightness)
		pwm.Set(channel, brightness)

		// Fade in and out
		if increasing {
			brightness += 100
			if brightness >= 16384 {
				brightness = 16384
				increasing = false
			}
		} else {
			if brightness >= 100 {
				brightness -= 100
			} else {
				brightness = 0
				increasing = true
			}
		}

		time.Sleep(10 * time.Millisecond)
	}
}
