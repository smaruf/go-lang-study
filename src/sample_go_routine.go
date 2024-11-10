v-package main

import (
	"fmt"
	"time"
)

func main() {
	// Create a channel to communicate with the goroutine
	dataChannel := make(chan int)

	// Start a goroutine to send data
	go func() {
		defer close(dataChannel) // Ensure the channel is closed when done

		// Simulate some work
		for i := 1; i <= 5; i++ {
			dataChannel <- i
			time.Sleep(time.Millisecond * 500) // Sleep to simulate processing time
		}
		// The channel will be closed here after the loop ends
	}()

	// Receive data from the channel
	for data := range dataChannel {
		fmt.Println("Received:", data)
	}

	fmt.Println("All data received, channel closed.")
}
