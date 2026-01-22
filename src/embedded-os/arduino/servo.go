package main

import (
	"machine"
	"time"
)

// Servo Motor Control Example for Arduino
// Hardware:
//   - Servo signal wire to Pin 9
//   - Servo power (red) to 5V
//   - Servo ground (brown/black) to GND

func main() {
	// Configure PWM for servo control
	servoPin := machine.D9
	servoPin.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// Configure PWM
	pwm := machine.PWM{Pin: servoPin}
	err := pwm.Configure(machine.PWMConfig{
		Period: 20000000, // 20ms period for servo (50Hz)
	})
	if err != nil {
		println("PWM configuration error:", err.Error())
		return
	}

	println("Servo Control Started")

	// Servo positions (in nanoseconds for 20ms period)
	// 1ms (1000000ns) = 0°
	// 1.5ms (1500000ns) = 90°
	// 2ms (2000000ns) = 180°

	positions := []uint32{
		1000000,  // 0 degrees
		1250000,  // 45 degrees
		1500000,  // 90 degrees
		1750000,  // 135 degrees
		2000000,  // 180 degrees
	}

	posNames := []string{"0°", "45°", "90°", "135°", "180°"}

	for {
		for i, pos := range positions {
			println("Moving to", posNames[i])
			pwm.Set(servoPin, pos)
			time.Sleep(1 * time.Second)
		}

		// Reverse direction
		for i := len(positions) - 2; i > 0; i-- {
			println("Moving to", posNames[i])
			pwm.Set(servoPin, positions[i])
			time.Sleep(1 * time.Second)
		}
	}
}
