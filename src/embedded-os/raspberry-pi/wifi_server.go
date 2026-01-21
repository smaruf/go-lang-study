package main

import (
	"machine"
	"time"
)

// WiFi Web Server Example for Raspberry Pi Pico W
// Hardware: Raspberry Pi Pico W (with WiFi)
// Note: This is a template - actual WiFi implementation requires additional setup

// Replace these with your WiFi credentials
const (
	WIFI_SSID     = "your-wifi-ssid"
	WIFI_PASSWORD = "your-wifi-password"
)

func main() {
	// Configure LED for status indication
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	println("Starting WiFi Web Server...")
	
	// Blink to show we're starting
	for i := 0; i < 5; i++ {
		led.Toggle()
		time.Sleep(200 * time.Millisecond)
	}

	// Note: Full WiFi implementation requires CYW43 driver
	// This is a simplified example showing the structure
	
	println("WiFi: Connecting to", WIFI_SSID)
	
	// TODO: Initialize WiFi driver
	// TODO: Connect to WiFi network
	// TODO: Get IP address
	// TODO: Start HTTP server
	
	println("Server running at http://192.168.1.100/")
	println("LED control at http://192.168.1.100/led/on")
	println("LED control at http://192.168.1.100/led/off")
	
	// Main loop
	ledState := false
	for {
		// Blink LED to show server is running
		led.Set(ledState)
		ledState = !ledState
		time.Sleep(1 * time.Second)
		
		// TODO: Handle HTTP requests
		// Example:
		// - GET /led/on  -> Turn LED on
		// - GET /led/off -> Turn LED off
		// - GET /status  -> Return JSON with LED state
	}
}

// HTTP handler functions (to be implemented with WiFi driver)

func handleLEDOn() string {
	// led.High()
	return `{"status":"ok","led":"on"}`
}

func handleLEDOff() string {
	// led.Low()
	return `{"status":"ok","led":"off"}`
}

func handleStatus() string {
	// ledState := led.Get()
	return `{"led":"unknown","uptime":0}`
}
