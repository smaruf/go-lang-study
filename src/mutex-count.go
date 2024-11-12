package main

import (
    "fmt"
    "sync"
    "time"
)

// A safeCounter is safe to use concurrently.
type safeCounter struct {
    val   int
    mutex sync.Mutex
}

// Inc increments the counter safely.
func (c *safeCounter) Inc() {
    c.mutex.Lock()         // Lock the mutex before modifying the value
    c.val++                // Increment the counter
    c.mutex.Unlock()       // Unlock the mutex after modifying the value
}

// Value returns the current value of the counter safely.
func (c *safeCounter) Value() int {
    c.mutex.Lock()         // Lock the mutex before reading the value
    defer c.mutex.Unlock() // Unlock the mutex after reading the value
    return c.val
}

func main() {
    var c safeCounter

    // Start several goroutines to increment the counter.
    for i := 0; i < 10; i++ {
        go func() {
            for j := 0; j < 1000; j++ {
                c.Inc()
            }
        }()
    }

    // Wait a bit for all goroutines to finish
    time.Sleep(1 * time.Second)

    // Output the final counter value.
    fmt.Println("Final counter value:", c.Value())
}
