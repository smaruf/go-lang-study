package main

import (
	"machine"
	"time"
)

// Temperature Sensor Example for Raspberry Pi Pico
// Uses the built-in temperature sensor (ADC channel 4)
// No external hardware required

func main() {
	// Configure ADC for temperature sensor
	machine.InitADC()
	
	// The Pico has a built-in temperature sensor on ADC4
	sensor := machine.ADC{Pin: machine.ADC4}
	sensor.Configure(machine.ADCConfig{})

	for {
		// Read ADC value
		raw := sensor.Get()
		
		// Convert to voltage (3.3V reference, 12-bit ADC)
		voltage := float32(raw) * 3.3 / 65536.0

		// Convert to temperature (Celsius)
		// Formula from Pico datasheet: T = 27 - (ADC_voltage - 0.706) / 0.001721
		temperature := 27.0 - (voltage-0.706)/0.001721

		// Print temperature
		println("Temperature:", int(temperature), "°C")
		println("Raw ADC:", raw)
		println("Voltage:", int(voltage*1000), "mV")
		println("---")

		time.Sleep(2 * time.Second)
	}
}
