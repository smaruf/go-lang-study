package main

import (
    "fmt"
    "time"
)

// Function that prints numbers
func printNumbers(ch chan int) {
    defer close(ch)  // Ensures the channel is closed when the function returns
    for i := 1; i <= 5; i++ {
        ch <- i // Send numbers to the channel
        time.Sleep(1 * time.Second) // Simulate work
    }
}

func main() {
    ch := make(chan int)  // Create a channel that transports integers

    go printNumbers(ch) // Start goroutine

    for num := range ch { // Read from channel until it's closed
        fmt.Println("Received:", num)
    }

    fmt.Println("Finished receiving!")
}
