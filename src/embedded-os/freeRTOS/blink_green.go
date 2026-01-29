package main

import (
    "time"
    "machine"
)

const (
    LED_PIN  = machine.LED
    RED_LED  = machine.GPIO0 // Assuming GPIO 0 for red LED
)

func GreenLEDTask() {
    for {
        LED_PIN.High()
        time.Sleep(1 * time.Second)
        LED_PIN.Low()
        time.Sleep(1 * time.Second)
    }
}

func RedLEDTask() {
    for {
        RED_LED.High()
        time.Sleep(100 * time.Millisecond)
        RED_LED.Low()
        time.Sleep(100 * time.Millisecond)
    }
}

func main() {
    LED_PIN.Configure(machine.PinConfig{Mode: machine.PinOutput})
    RED_LED.Configure(machine.PinConfig{Mode: machine.PinOutput})

    go GreenLEDTask()
    // go RedLEDTask() // Uncomment if needed

    // Keep the main function running
    select {}
}
