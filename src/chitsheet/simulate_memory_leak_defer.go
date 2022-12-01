func main() {
	printMemUsage("at start of program")
	startingRoutines := runtime.NumGoroutine()

	// won't capture returned channel
	funcWithLeakingGoRoutine()
	runtime.GC() // enforce GC cycle
  
  // pause main thread to allow the program to settle 
	time.Sleep(5 * time.Second)
  
	endingRoutines := runtime.NumGoroutine()

	fmt.Println("goroutines on start:", startingRoutines)
	fmt.Println("goroutines at end:", endingRoutines)
	fmt.Println("Leaking goroutines:", endingRoutines-startingRoutines)

	printMemUsage("just before program exit")
}

func funcWithLeakingGoRoutine() <-chan int {
	ch := make(chan int)

	randomString := createRandomStringOfSize(1 << 20) // 1MB
	printMemUsage("after allocating big chunk of memory")
  
	go func() {
		// refer memory intensive resource in a deferred closure inside the goroutine
    defer func() {
		  fmt.Printf("length of big string: %d\n", len(randomString))
		  
		  ch <- 5
		  fmt.Println("sent value via channel")
    }()
	}()

	return ch
}

/*
Outputs the following (running w/ older golang version):
----------------------
[at start of program] mem usage: 103688
[after allocating big chunk of memory] mem usage: 2202048
length of big string: 1048576
goroutines on start: 1
goroutines at end: 2
Leaking goroutines: 1
[just before program exit] mem usage: 97544
*/

/*
Outputs the following (running w/ newer golang version):
----------------------
[at start of program] mem usage: 108464
[after allocating big chunk of memory] mem usage: 2206800
length of big string: 1048576
goroutines on start: 1
goroutines at end: 2
Leaking goroutines: 1
[just before program exit] mem usage: 1150608
*/
