package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	startingRoutines := runtime.NumGoroutine()

	// won't capture returned channel
	funcWithLeakingGoRoutine()

	// sleep to allow some time for the goroutine to exit (it won't exit since it's blocked by the channel)
	time.Sleep(5 * time.Second)

	endingRoutines := runtime.NumGoroutine()

	fmt.Println("goroutines on start:", startingRoutines)
	fmt.Println("goroutines at end:", endingRoutines)
	fmt.Println("Leaking goroutines:", endingRoutines-startingRoutines)
}

func funcWithLeakingGoRoutine() <-chan int {
	ch := make(chan int)

	go func() {
		ch <- 5
		fmt.Println("sent value via channel")
	}()

	return ch
}
