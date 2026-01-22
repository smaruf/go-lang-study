package main

import (
	"fmt"
	"sync"
)

// Pipeline demonstrates the pipeline pattern for concurrent data processing
// Data flows through multiple stages, each stage processes and passes to next

// Stage represents a processing stage in the pipeline
type Stage func(<-chan interface{}) <-chan interface{}

// Pipeline builds and executes a pipeline of stages
func Pipeline(stages ...Stage) Stage {
	return func(in <-chan interface{}) <-chan interface{} {
		out := in
		for _, stage := range stages {
			out = stage(out)
		}
		return out
	}
}

// Generator creates a channel and sends values to it
func generator(nums ...int) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

// Square multiplies each number by itself
func square(in <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)
		for n := range in {
			num := n.(int)
			out <- num * num
		}
	}()
	return out
}

// Double multiplies each number by 2
func double(in <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)
		for n := range in {
			num := n.(int)
			out <- num * 2
		}
	}()
	return out
}

// AddTen adds 10 to each number
func addTen(in <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)
		for n := range in {
			num := n.(int)
			out <- num + 10
		}
	}()
	return out
}

// FanOut creates multiple workers to process input concurrently
func fanOut(in <-chan interface{}, workers int, processor Stage) []<-chan interface{} {
	channels := make([]<-chan interface{}, workers)
	for i := 0; i < workers; i++ {
		channels[i] = processor(in)
	}
	return channels
}

// FanIn merges multiple channels into one
func fanIn(channels ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	var wg sync.WaitGroup

	multiplex := func(c <-chan interface{}) {
		defer wg.Done()
		for n := range c {
			out <- n
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	fmt.Println("=== Simple Pipeline Example ===")
	// Create a simple pipeline: generator -> square -> double -> addTen
	input := generator(1, 2, 3, 4, 5)
	pipeline := Pipeline(square, double, addTen)
	output := pipeline(input)

	for result := range output {
		fmt.Printf("%v ", result)
	}
	fmt.Println()

	fmt.Println("\n=== Fan-Out/Fan-In Pattern Example ===")
	// Demonstrate fan-out/fan-in for parallel processing
	input2 := generator(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	
	// Fan-out: Create 3 workers to square numbers concurrently
	workers := fanOut(input2, 3, square)
	
	// Fan-in: Merge results from all workers
	merged := fanIn(workers...)
	
	// Further process merged results
	final := double(merged)

	results := make([]int, 0)
	for result := range final {
		results = append(results, result.(int))
	}
	fmt.Printf("Results: %v\n", results)

	fmt.Println("\n=== Complex Multi-Stage Pipeline ===")
	// Build a complex pipeline with multiple stages
	complexInput := generator(1, 2, 3, 4, 5)
	stage1 := square(complexInput)
	stage2 := addTen(stage1)
	stage3 := double(stage2)

	fmt.Print("Processing: ")
	for result := range stage3 {
		fmt.Printf("%v ", result)
	}
	fmt.Println()
}
