package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)

func worker(id int, ch chan<- int, wg *sync.WaitGroup) {
    defer wg.Done() // Signal that this worker is done

    // Simulate some work by sleeping
    sleepDuration := rand.Intn(5) + 1
    fmt.Printf("Worker %d is working for %d seconds\n", id, sleepDuration)
    time.Sleep(time.Duration(sleepDuration) * time.Second)
    
    // Send the result to the channel
    ch <- sleepDuration
}

func main() {
    // Create a channel to communicate the results
    results := make(chan int)

    // Use a WaitGroup to wait for all goroutines to finish
    var wg sync.WaitGroup

    // Number of workers/goroutines
    numWorkers := 5
    wg.Add(numWorkers) // Set the number of goroutines we'll wait for

    // Start multiple workers
    for i := 1; i <= numWorkers; i++ {
        go worker(i, results, &wg)
    }

    // Close the channel in the background once all goroutines are done
    go func() {
        wg.Wait()
        close(results)
    }()

    // Collect all results from the channel
    for result := range results {
        fmt.Printf("Received result: %d seconds\n", result)
    }

    fmt.Println("All workers finished.")
}
