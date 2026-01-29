package main

import (
	"machine"
	"time"
)

var period uint64 = 1e9 / 500

func main() {
    // This program is specific to the Raspberry Pi Pico.
    pin := machine.LED
	pwm := machine.PWM4 // Pin 25 (LED on pico) corresponds to PWM4.

	// Configure the PWM with the given period.
	pwm.Configure(machine.PWMConfig{
		Period: period,
	})

	ch, err := pwm.Channel(pin)
	if err != nil {
		println(err.Error())
		return
	}

	for { 
		for i := 1; i < 255; i++ {
            // This performs a stylish fade-out blink
			pwm.Set(ch, pwm.Top()/uint32(i))
			time.Sleep(time.Millisecond * 5)
		}
	}
}
