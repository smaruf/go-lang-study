package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Job represents a unit of work
type Job struct {
	ID      int
	Payload string
}

// Result represents the result of a job
type Result struct {
	JobID  int
	Output string
	Err    error
}

// WorkerPool manages a pool of workers
type WorkerPool struct {
	numWorkers int
	jobs       chan Job
	results    chan Result
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(numWorkers int, jobBufferSize int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		numWorkers: numWorkers,
		jobs:       make(chan Job, jobBufferSize),
		results:    make(chan Result, jobBufferSize),
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Start initializes and starts all workers
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// worker processes jobs from the jobs channel
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	fmt.Printf("Worker %d started\n", id)

	for {
		select {
		case <-wp.ctx.Done():
			fmt.Printf("Worker %d shutting down\n", id)
			return
		case job, ok := <-wp.jobs:
			if !ok {
				fmt.Printf("Worker %d: jobs channel closed\n", id)
				return
			}
			result := wp.processJob(id, job)
			wp.results <- result
		}
	}
}

// processJob simulates job processing
func (wp *WorkerPool) processJob(workerID int, job Job) Result {
	fmt.Printf("Worker %d processing job %d: %s\n", workerID, job.ID, job.Payload)
	
	// Simulate work
	time.Sleep(time.Millisecond * 100)
	
	return Result{
		JobID:  job.ID,
		Output: fmt.Sprintf("Processed by worker %d: %s", workerID, job.Payload),
		Err:    nil,
	}
}

// Submit adds a job to the pool
func (wp *WorkerPool) Submit(job Job) {
	wp.jobs <- job
}

// Close gracefully shuts down the worker pool
func (wp *WorkerPool) Close() {
	close(wp.jobs)
	wp.wg.Wait()
	close(wp.results)
}

// Shutdown forcefully stops all workers
func (wp *WorkerPool) Shutdown() {
	wp.cancel()
	wp.wg.Wait()
	close(wp.results)
}

// Results returns the results channel
func (wp *WorkerPool) Results() <-chan Result {
	return wp.results
}

func main() {
	// Create a worker pool with 3 workers and buffer size of 10
	pool := NewWorkerPool(3, 10)
	pool.Start()

	// Submit jobs
	numJobs := 10
	go func() {
		for i := 1; i <= numJobs; i++ {
			job := Job{
				ID:      i,
				Payload: fmt.Sprintf("Task-%d", i),
			}
			pool.Submit(job)
		}
		pool.Close() // Close when all jobs are submitted
	}()

	// Collect results
	resultsReceived := 0
	for result := range pool.Results() {
		if result.Err != nil {
			fmt.Printf("Job %d failed: %v\n", result.JobID, result.Err)
		} else {
			fmt.Printf("Job %d completed: %s\n", result.JobID, result.Output)
		}
		resultsReceived++
	}

	fmt.Printf("\nAll jobs completed! Processed %d jobs\n", resultsReceived)
}
