package main

import (
    "fmt"
    "sync"
    "sync/atomic"
    "time"
)

func main() {
    var counter int32 = 0
    var wg sync.WaitGroup
    wg.Add(2)

    increment := func() {
        defer wg.Done()
        for i := 0; i < 1000; i++ {
            atomic.AddInt32(&counter, 1)
            time.Sleep(time.Microsecond)  // simulate some processing time
        }
    }

    // Start two goroutines to increase the counter
    go increment()
    go increment()
    wg.Wait()

    fmt.Println("Final Counter:", counter)
}
