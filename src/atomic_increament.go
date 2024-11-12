package main

import (
    "fmt"
    "sync"
    "sync/atomic"
)

func main() {
    var count int32
    var wg sync.WaitGroup
    increment := func() {
        defer wg.Done()
        for i := 0; i < 1000; i++ {
            atomic.AddInt32(&count, 1)
        }
    }

    // Start several goroutines to perform increments
    wg.Add(4)
    go increment()
    go increment()
    go increment()
    go increment()
    wg.Wait()
    
    fmt.Println("Count:", count)  // Output: Count: 4000
}
