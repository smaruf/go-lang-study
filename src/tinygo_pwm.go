package main

import (
	"machine"
	"time"
)

func main() {
    pin := machine.LED
    // This program is specific to the Raspberry Pi Pico.
	pin.Configure(machine.PinConfig{Mode: machine.PinPWM})
	pwm, err := machine.NewPWM(pin)
	if err != nil {
		println(err.Error())
	}
	var period uint64 = 1e9 / 500
	err = pwm.Configure(machine.PWMConfig{Period: period})
	if err != nil {
		println(err.Error())
	}
	ch, err := pwm.Channel(pin)
	if err != nil {
		println(err.Error())
	}
	for { 
		for i := 1; i < 255; i++ {
            // This performs a stylish fade-out blink
			pwm.Set(ch, pwm.Top()/uint32(i))
			time.Sleep(time.Millisecond * 5)
		}
	}
}

// flash: tinygo flash -target=pico
