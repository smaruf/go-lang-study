package main

import (
	"machine"
	"time"
)

// Serial Communication Example for Arduino
// Hardware: Uses built-in serial (USB)
// Monitor with: tinygo monitor

func main() {
	println("Arduino Serial Communication Example")
	println("Sending messages every second...")

	counter := 0

	for {
		counter++
		println("Message", counter, "- Time:", time.Now().Unix())
		println("Temperature: 25°C")
		println("Status: OK")
		println("---")
		
		time.Sleep(1 * time.Second)
	}
}
