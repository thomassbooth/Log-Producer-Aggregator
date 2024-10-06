package internal

import (
	"fmt"
	"log-aggregator/shared"
	"sync/atomic"
	"time"
)

type Worker struct {
	id     int
	jobs   <-chan shared.LogMessage
	quit   chan struct{}
	active *int32
}

type WorkerPool struct {
	jobs        chan shared.LogMessage
	quit        chan struct{}
	workers     []*Worker
	activeCount int32
}

func NewWorkerPool(numWorkers int) *WorkerPool {

	jobs := make(chan shared.LogMessage, 100) // Buffer to hold incoming jobs

	pool := &WorkerPool{
		jobs:        jobs,
		quit:        make(chan struct{}),
		workers:     make([]*Worker, numWorkers),
		activeCount: 0,
	}

	// Setup workers and put them in the pool
	for i := 0; i < numWorkers; i++ {
		worker := Worker{
			id:     i,
			jobs:   jobs,
			quit:   make(chan struct{}),
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

	for {
		select {
		case logMsg := <-w.jobs:
			// Process the log message
			time.Sleep(10 * time.Second)
			fmt.Printf("Worker %d processing log: %s\n", w.id, logMsg.Message)
			// Here you would add code to save logMsg to MongoDB or other storage
		case <-w.quit:
			fmt.Printf("Worker %d stopping\n", w.id)
			return // Exit the worker's loop to stop it
		}
	}
}

// Stop stops the worker by sending a signal to its quit channel.
func (w *Worker) Stop() {
	close(w.quit) // Close the quit channel to signal the worker to stop
}

func (wp *WorkerPool) AddJob(logMsg shared.LogMessage) {
	wp.jobs <- logMsg
}

// Stop stops all workers in the pool.
func (wp *WorkerPool) Stop() {
	for _, worker := range wp.workers {
		worker.Stop() // Call each worker's Stop method
	}
	close(wp.quit) // Close the main quit channel if needed (optional)
}

// ActiveWorkers returns the number of active workers.
func (wp *WorkerPool) ActiveWorkers() int {
	return int(atomic.LoadInt32(&wp.activeCount))
}

// QueuedTasks returns the number of queued tasks.
func (wp *WorkerPool) QueuedTasks() int {
	return len(wp.jobs)
}
