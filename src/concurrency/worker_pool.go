package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

// Job represents work to be done
type Job struct {
	ID       int
	Duration time.Duration
	Name     string
}

// Result represents the result of processing a job
type Result struct {
	JobID     int
	Success   bool
	Message   string
	Duration  time.Duration
	ProcessedAt time.Time
}

// Worker represents a worker goroutine
type Worker struct {
	ID         int
	jobChannel chan Job
	resultCh   chan Result
	quit       chan bool
}

// WorkerPool manages a pool of workers
type WorkerPool struct {
	workers    []*Worker
	jobQueue   chan Job
	resultCh   chan Result
	numWorkers int
	wg         sync.WaitGroup
}

// NewWorker creates a new worker
func NewWorker(id int, jobCh chan Job, resultCh chan Result) *Worker {
	return &Worker{
		ID:         id,
		jobChannel: jobCh,
		resultCh:   resultCh,
		quit:       make(chan bool),
	}
}

// Start begins the worker's processing loop
func (w *Worker) Start(ctx context.Context) {
	go func() {
		defer func() {
			log.Printf("Worker %d stopped", w.ID)
		}()
		
		for {
			select {
			case job := <-w.jobChannel:
				start := time.Now()
				log.Printf("Worker %d started job %d (%s)", w.ID, job.ID, job.Name)
				
				// Simulate work with context cancellation support
				select {
				case <-time.After(job.Duration):
					w.resultCh <- Result{
						JobID:       job.ID,
						Success:     true,
						Message:     fmt.Sprintf("Job %d completed by worker %d", job.ID, w.ID),
						Duration:    time.Since(start),
						ProcessedAt: time.Now(),
					}
					log.Printf("Worker %d completed job %d", w.ID, job.ID)
				case <-ctx.Done():
					w.resultCh <- Result{
						JobID:       job.ID,
						Success:     false,
						Message:     fmt.Sprintf("Job %d cancelled", job.ID),
						Duration:    time.Since(start),
						ProcessedAt: time.Now(),
					}
					log.Printf("Worker %d cancelled job %d", w.ID, job.ID)
					return
				}
			case <-ctx.Done():
				return
			case <-w.quit:
				return
			}
		}
	}()
}

// Stop stops the worker
func (w *Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(numWorkers int) *WorkerPool {
	return &WorkerPool{
		workers:    make([]*Worker, numWorkers),
		jobQueue:   make(chan Job, 100),
		resultCh:   make(chan Result, 100),
		numWorkers: numWorkers,
	}
}

// Start initializes and starts all workers
func (wp *WorkerPool) Start(ctx context.Context) {
	log.Printf("Starting worker pool with %d workers", wp.numWorkers)
	
	for i := 0; i < wp.numWorkers; i++ {
		worker := NewWorker(i+1, wp.jobQueue, wp.resultCh)
		wp.workers[i] = worker
		worker.Start(ctx)
	}
	
	// Start result processor
	go wp.processResults(ctx)
}

// processResults handles the results from workers
func (wp *WorkerPool) processResults(ctx context.Context) {
	for {
		select {
		case result := <-wp.resultCh:
			status := "SUCCESS"
			if !result.Success {
				status = "FAILED"
			}
			log.Printf("[%s] Job %d: %s (took %v)", 
				status, result.JobID, result.Message, result.Duration)
			wp.wg.Done()
		case <-ctx.Done():
			log.Println("Result processor stopped")
			return
		}
	}
}

// SubmitJob adds a job to the job queue
func (wp *WorkerPool) SubmitJob(job Job) {
	wp.wg.Add(1)
	go func() {
		wp.jobQueue <- job
	}()
}

// Wait waits for all jobs to complete
func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}

// Stop gracefully stops all workers
func (wp *WorkerPool) Stop() {
	log.Println("Stopping worker pool...")
	for _, worker := range wp.workers {
		worker.Stop()
	}
}

func main() {
	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create and start worker pool
	pool := NewWorkerPool(3)
	pool.Start(ctx)

	// Submit jobs
	jobs := []Job{
		{ID: 1, Duration: 2 * time.Second, Name: "Process Data"},
		{ID: 2, Duration: 1 * time.Second, Name: "Send Email"},
		{ID: 3, Duration: 3 * time.Second, Name: "Generate Report"},
		{ID: 4, Duration: 500 * time.Millisecond, Name: "Update Cache"},
		{ID: 5, Duration: 1500 * time.Millisecond, Name: "Backup Database"},
		{ID: 6, Duration: 800 * time.Millisecond, Name: "Clean Logs"},
		{ID: 7, Duration: 2500 * time.Millisecond, Name: "Sync Files"},
		{ID: 8, Duration: 600 * time.Millisecond, Name: "Index Search"},
	}

	log.Printf("Submitting %d jobs...", len(jobs))
	
	// Add some randomness to job duration
	rand.Seed(time.Now().UnixNano())
	for _, job := range jobs {
		// Add random variation to duration (±50%)
		variation := time.Duration(rand.Intn(int(job.Duration/2)))
		if rand.Intn(2) == 0 {
			job.Duration += variation
		} else {
			job.Duration -= variation
		}
		pool.SubmitJob(job)
	}

	// Wait for all jobs to complete
	log.Println("Waiting for jobs to complete...")
	pool.Wait()

	// Graceful shutdown
	pool.Stop()
	log.Println("All jobs completed!")
}