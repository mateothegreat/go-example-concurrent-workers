package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"golang.org/x/exp/rand"
)

var (
	// maxConcurrent is the maximum number of concurrent workers.
	maxConcurrent = 10
	// totalJobs is the total number of jobs to be processed.
	totalJobs = 10
)

// Result is the result of a job that is gathered from the results channel.
type Result struct {
	Message string
}

// worker is a worker that processes jobs from the jobs
// channel and sends the results to the results channel.
// It uses a semaphore to limit the number of concurrent workers.
func worker(id int, jobs <-chan string, results chan<- Result, resultsWg *sync.WaitGroup) {
	log.Printf("Worker %d started, waiting for jobs...", id)
	for job := range jobs {
		// Process the job (dummy processing here)
		time.Sleep(time.Duration(1+rand.Intn(3000)) * time.Millisecond)
		results <- Result{Message: fmt.Sprintf("Worker %d processed job: %s", id, job)}
		resultsWg.Done() // Mark the result as done
	}
}

// TryEnqueue tries to enqueue a job to the jobs channel.
// If the channel is full, it returns false so that the caller can retry.
func TryEnqueue(job string, jobs chan<- string) bool {
	select {
	case jobs <- job:
		return true
	default:
		return false
	}
}

func main() {
	jobs := make(chan string, totalJobs)
	results := make(chan Result, totalJobs)
	sem := make(chan struct{}, maxConcurrent)

	// First, wait for job scheduling to finish.
	var jobWg sync.WaitGroup

	// Second, wait for all results to be processed.
	var resultsWg sync.WaitGroup

	// Create and start workers.
	// This is a fan-out pattern where we create a pool of workers
	// and let them process jobs concurrently as they receive them
	// from the jobs channel.
	for i := 0; i < maxConcurrent; i++ {
		go worker(i, jobs, results, &resultsWg)
	}

	// Schedule jobs to be processed.
	// This is a fan-in pattern where we create a pool of jobs
	// and let the workers process them concurrently.
	for i := 0; i < totalJobs; i++ {
		job := time.Now().Format("15:04:05")

		// Tell the wait group that we're adding a new job and
		// to wait for it to finish.
		jobWg.Add(1)
		resultsWg.Add(1) // Add to results wait group

		// Acquire a slot from the semaphore.
		sem <- struct{}{}

		// Spawn a new goroutine to process each job.
		go func(job string) {
			// This deferred closure is used to release the semaphore
			// when the goroutine is done or when the goroutine is cancelled.
			defer func() {
				// Release the slot when the goroutine is done.
				<-sem
				// Tell the wait group that the job is done.
				jobWg.Done()
			}()

			// Try to enqueue the job to the jobs channel.
			// If the channel is full, it will block and the goroutine will be blocked.
			for !TryEnqueue(job, jobs) {
				log.Println("Queue is full, waiting...")
				time.Sleep(time.Second)
			}
		}(job)
	}

	// Wait for all jobs to be enqueued
	jobWg.Wait()

	// Close the jobs channel to signal to the workers that there are no more jobs.
	close(jobs)

	// Collect results concurrently in a separate goroutine
	// to avoid blocking the main goroutine. This is a trickle-down
	// pattern where we create a pool of workers to process the results
	// concurrently and then collect the results as they come in.
	go func() {
		for result := range results {
			log.Printf("Result received: %s", result.Message)
		}
	}()

	// Wait for all results to be processed
	resultsWg.Wait()
	close(results) // Close the results channel to signal that no more results will be sent

	log.Println("All jobs processed. Shutting down...")
}
