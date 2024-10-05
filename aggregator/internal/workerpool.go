// In log-aggregator/aggregator/internal/workerpool.go
package internal

import (
	"fmt"
)

type LogMessage struct {
	Timestamp string
	Level     string
	Message   string
}

type Worker struct {
	id   int
	jobs <-chan LogMessage
	quit chan struct{}
}

type WorkerPool struct {
	jobs    chan LogMessage
	quit    chan struct{}
	workers []*Worker
}

func NewWorkerPool(numWorkers int) *WorkerPool {

	jobs := make(chan LogMessage, 100) // Buffer to hold incoming jobs

	pool := &WorkerPool{
		jobs:    jobs, // Buffer to hold incoming jobs
		quit:    make(chan struct{}),
		workers: make([]*Worker, numWorkers),
	}

	// setup workers and put them in the pool
	for i := 0; i < numWorkers; i++ {
		worker := Worker{
			id:   i,
			jobs: jobs,
			quit: make(chan struct{}),
		}
		pool.workers[i] = &worker
		//start our workers on new threads
		go worker.start()
	}

	return pool
}

// Start starts the worker's job processing loop.
func (w *Worker) start() {
	for {
		select {
		case logMsg := <-w.jobs:

			// Process the log message

			// shared.Log("Worker %d processing log: %s", w.id, logMsg.Message)
			fmt.Printf("Worker %d processing log: %s\n", w.id, logMsg.Message)
			// Here you would add code to save logMsg to MongoDB
		case <-w.quit:
			return
		}
	}
}

func (wp *WorkerPool) AddJob(logMsg LogMessage) {
	wp.jobs <- logMsg
}

// Stop stops all workers in the pool.
func (wp *WorkerPool) Stop() {
	for _, worker := range wp.workers {
		close(worker.quit) // Signal each worker to stop
	}
	close(wp.quit) // Close the main quit channel
}
