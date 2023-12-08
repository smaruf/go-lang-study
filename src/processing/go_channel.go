package main

import (
 "fmt"
 "sync"
)

func process(data int, wg *sync.WaitGroup, ch chan int) {
 defer wg.Done()
 result := data * 2 // some processing
 ch <- result
}

func main() {
 var wg sync.WaitGroup
 ch := make(chan int, 5) // Channel to communicate results

 for i := 0; i < 5; i++ {
 wg.Add(1)
 go process(i, &wg, ch) // Starting a goroutine
 }

 wg.Wait()
 close(ch)

 for result := range ch {
 fmt.Println(result)
 }
}
