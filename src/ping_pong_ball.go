package main

import (
    "context"
    "fmt"
    "time"
)

func player(name string, table chan string, ctx context.Context) {
    for {
        select {
        case ball := <-table:
            fmt.Println(name, "received", ball)
            // Sleep to simulate the time taken to hit the ball back
            time.Sleep(100 * time.Millisecond)
            table <- ball
        case <-ctx.Done():
            fmt.Println(name, "wins the game!")
            return
        }
    }
}

func main() {
    // Creating a channel for passing the ping pong ball
    table := make(chan string)

    // Context to handle the timeout
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
    defer cancel() // Ensure resources are released after the context is no longer needed

    // Start two player goroutines
    go player("Player 1", table, ctx)
    go player("Player 2", table, ctx)

    // Start the game by serving the ball to the table
    table <- "ball"

    // Waiting for the game to finish
    <-ctx.Done()

    // Close the channel after the context deadline is exceeded
    close(table)

    // Wait for a moment to see the final print statements
    time.Sleep(1 * time.Second)
    fmt.Println("Game over.")
}
