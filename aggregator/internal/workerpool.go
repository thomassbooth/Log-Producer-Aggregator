package internal

import (
	"fmt"
	"sync/atomic"
	"time"
)

type LogMessage struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
}
type Worker struct {
	id     int
	jobs   <-chan LogMessage
	quit   <-chan struct{}
	active *int32
}

type WorkerPool struct {
	jobs        chan LogMessage
	quit        chan struct{}
	workers     []*Worker
	activeCount int32
}

func NewWorkerPool(numWorkers int) *WorkerPool {

	jobs := make(chan LogMessage, 100) // Buffer to hold incoming jobs
	quit := make(chan struct{})        // Channel to signal worker to stop
	pool := &WorkerPool{
		jobs:        jobs,
		quit:        quit,
		workers:     make([]*Worker, numWorkers),
		activeCount: 0,
	}

	// Setup workers and put them in the pool
	for i := 0; i < numWorkers; i++ {
		worker := Worker{
			id:     i,
			jobs:   jobs,
			quit:   quit,
			active: &pool.activeCount,
		}
		pool.workers[i] = &worker
		// Start each worker in a new goroutine
		go worker.start()
	}

	return pool
}

// Start starts the worker's job processing loop.
func (w *Worker) start() {
	atomic.AddInt32(w.active, 1)        // Increment active worker count
	defer atomic.AddInt32(w.active, -1) // Decrement active worker count when done

free:
	for {
		select {
		case logMsg := <-w.jobs:
			// Process the log message
			time.Sleep(10 * time.Second)
			fmt.Printf("Worker %d processing log: %s\n", w.id, logMsg.Message)
			// Here you would add code to save logMsg to MongoDB or other storage
		case <-w.quit:
			fmt.Printf("Worker %d stopping\n", w.id)
			w.Stop()
			break free // I want to break freeeee
		}
	}
}

// Stop stops the worker by sending a signal to its quit channel.
func (w *Worker) Stop() {
}

func (wp *WorkerPool) AddJob(logMsg LogMessage) {
	wp.jobs <- logMsg
}

// Stop stops all workers in the pool.
func (wp *WorkerPool) Stop() {
	// Signal all workers to stop
	for range wp.workers {
		wp.quit <- struct{}{} // Send stop signal to worker
	}

	// Optionally close the jobs channel to prevent further job submissions
	close(wp.jobs) // This is optional
}

// ActiveWorkers returns the number of active workers.
func (wp *WorkerPool) ActiveWorkers() int {
	return int(atomic.LoadInt32(&wp.activeCount))
}

// QueuedTasks returns the number of queued tasks.
func (wp *WorkerPool) QueuedTasks() int {
	return len(wp.jobs)
}
