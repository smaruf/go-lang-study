package main

import (
    "context"
    "fmt"
    "time"
)

func operation(ctx context.Context) {
    select {
    case <-time.After(5 * time.Second): // Simulates a 5-second task
        fmt.Println("operation completed")
    case <-ctx.Done():
        fmt.Println("operation canceled:", ctx.Err())
    }
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    go operation(ctx)

    time.Sleep(3 * time.Second) // Simulate other work
}
