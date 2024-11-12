package main

import (
    "context"
    "fmt"
    "math/rand"
    "time"
)

func player(name string, table chan string, ctx context.Context) {
    for {
        select {
        case ball := <-table:
            fmt.Println(name, "received", ball)
            // Simulate random wait time between 0 to 10 seconds
            waitTime := time.Duration(rand.Intn(10)) * time.Second
            time.Sleep(waitTime)

            select {
            case table <- "ball":
                fmt.Println(name, "hits the ball back")
            case <-ctx.Done():
                // Context is done, so stop the game
                fmt.Println(name, "holds the ball. Game over!")
                return
            }
        case <-ctx.Done():
            // If the game times out and this player doesn't have the ball
            return
        }
    }
}

func main() {
    rand.Seed(time.Now().UnixNano()) // Seed random number generator

    // Creating a channel for passing the ping pong ball
    table := make(chan string)

    // Set up a context that will be cancelled after 2 minutes
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
    defer cancel()

    // Start two player goroutines
    go player("Player 1", table, ctx)
    go player("Player 2", table, ctx)

    // Start the game by serving the ball to the table
    table <- "ball"

    // Block until the context is done
    <-ctx.Done()

    // Close the channel after the context deadline is exceeded
    close(table)

    // Just wait a bit for final messages to print before main exits
    time.Sleep(1 * time.Second)
    fmt.Println("Ping pong game finished.")
}
