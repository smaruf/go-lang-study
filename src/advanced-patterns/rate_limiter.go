package main

import (
	"context"
	"fmt"
	"time"
)

// RateLimiter implements a token bucket rate limiter
type RateLimiter struct {
	tokens       chan struct{}
	fillInterval time.Duration
	capacity     int
	ctx          context.Context
	cancel       context.CancelFunc
}

// NewRateLimiter creates a new rate limiter
// rate: number of tokens per second
// burst: maximum number of tokens (bucket capacity)
func NewRateLimiter(rate int, burst int) *RateLimiter {
	ctx, cancel := context.WithCancel(context.Background())
	rl := &RateLimiter{
		tokens:       make(chan struct{}, burst),
		fillInterval: time.Second / time.Duration(rate),
		capacity:     burst,
		ctx:          ctx,
		cancel:       cancel,
	}

	// Fill initial tokens
	for i := 0; i < burst; i++ {
		rl.tokens <- struct{}{}
	}

	// Start token refiller
	go rl.refiller()

	return rl
}

// refiller continuously adds tokens to the bucket
func (rl *RateLimiter) refiller() {
	ticker := time.NewTicker(rl.fillInterval)
	defer ticker.Stop()

	for {
		select {
		case <-rl.ctx.Done():
			return
		case <-ticker.C:
			select {
			case rl.tokens <- struct{}{}:
				// Token added
			default:
				// Bucket is full, skip
			}
		}
	}
}

// Allow returns true if the request is allowed
func (rl *RateLimiter) Allow() bool {
	select {
	case <-rl.tokens:
		return true
	default:
		return false
	}
}

// Wait blocks until a token is available or context is cancelled
func (rl *RateLimiter) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-rl.tokens:
		return nil
	}
}

// Stop shuts down the rate limiter
func (rl *RateLimiter) Stop() {
	rl.cancel()
}

// APIClient demonstrates rate-limited API client
type APIClient struct {
	limiter *RateLimiter
}

// NewAPIClient creates a new API client with rate limiting
func NewAPIClient(requestsPerSecond int, burst int) *APIClient {
	return &APIClient{
		limiter: NewRateLimiter(requestsPerSecond, burst),
	}
}

// MakeRequest simulates making an API request
func (c *APIClient) MakeRequest(requestID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Printf("Request %d: Waiting for rate limiter...\n", requestID)
	start := time.Now()

	if err := c.limiter.Wait(ctx); err != nil {
		return fmt.Errorf("rate limiter error: %w", err)
	}

	elapsed := time.Since(start)
	fmt.Printf("Request %d: Allowed after %v, making API call...\n", requestID, elapsed)

	// Simulate API call
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("Request %d: Completed\n", requestID)

	return nil
}

// Close shuts down the API client
func (c *APIClient) Close() {
	c.limiter.Stop()
}

func main() {
	fmt.Println("=== Basic Rate Limiter Demo ===")
	
	// Create a rate limiter: 2 requests per second, burst of 3
	limiter := NewRateLimiter(2, 3)
	defer limiter.Stop()

	// Try 10 requests quickly
	fmt.Println("\nAttempting 10 rapid requests (2/sec limit, burst 3):")
	for i := 1; i <= 10; i++ {
		allowed := limiter.Allow()
		if allowed {
			fmt.Printf("Request %d: ✓ Allowed\n", i)
		} else {
			fmt.Printf("Request %d: ✗ Rate limited\n", i)
		}
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("\n=== API Client with Rate Limiting ===")
	
	// Create API client with rate limiting
	client := NewAPIClient(3, 5) // 3 requests per second, burst of 5
	defer client.Close()

	// Make multiple requests
	numRequests := 10
	for i := 1; i <= numRequests; i++ {
		go func(id int) {
			if err := client.MakeRequest(id); err != nil {
				fmt.Printf("Request %d failed: %v\n", id, err)
			}
		}(i)
	}

	// Wait for all requests to complete
	time.Sleep(5 * time.Second)

	fmt.Println("\n=== Concurrent Rate Limiting Test ===")
	
	// Test concurrent access
	concurrentLimiter := NewRateLimiter(5, 10)
	defer concurrentLimiter.Stop()

	start := time.Now()
	const concurrentRequests = 20
	done := make(chan bool, concurrentRequests)

	for i := 1; i <= concurrentRequests; i++ {
		go func(id int) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			if err := concurrentLimiter.Wait(ctx); err != nil {
				fmt.Printf("Request %d: Failed - %v\n", id, err)
			} else {
				fmt.Printf("Request %d: Success at %v\n", id, time.Since(start))
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to finish
	for i := 0; i < concurrentRequests; i++ {
		<-done
	}

	fmt.Printf("\nCompleted %d requests in %v\n", concurrentRequests, time.Since(start))
}
