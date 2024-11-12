package main

import (
    "context"
    "errors"
    "fmt"
    "math/rand"
    "time"
)

// Simulated operation that fails about half the time.
func flakyOperation() error {
    if rand.Intn(2) == 0 {
        return errors.New("operation failed")
    }
    return nil
}

// retry executes the provided function f until it succeeds or the context is done.
func retry(ctx context.Context, maxRetries int, f func() error) error {
    for i := 0; i < maxRetries; i++ {
        err := f()
        if err == nil {
            return nil
        }
        fmt.Println("Attempt", i+1, "failed; retrying...")

        // Wait for a second or until the context is cancelled/timed out
        select {
        case <-time.After(1 * time.Second):
        case <-ctx.Done():
            return ctx.Err()
        }
    }
    return errors.New("max retries exceeded")
}

func main() {
    // Create a context that times out after 5 seconds
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    err := retry(ctx, 5, flakyOperation)
    if err != nil {
        fmt.Println("Operation failed:", err)
    } else {
        fmt.Println("Operation succeeded.")
    }
}
