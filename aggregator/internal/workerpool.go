package internal

import (
	"fmt"
	"log-aggregator/aggregator/storage"
	"log-aggregator/aggregator/utils"
	"sync/atomic"
)

type Worker struct {
	id     int
	jobs   <-chan utils.Job
	quit   <-chan struct{}
	active *int32
	store  *storage.Storage
}

type WorkerPool struct {
	jobs        chan utils.Job
	quit        chan struct{}
	workers     []*Worker
	activeCount int32
}

func NewWorkerPool(numWorkers int, store *storage.Storage) *WorkerPool {

	jobs := make(chan utils.Job, 100) // Buffer to hold incoming jobs
	quit := make(chan struct{})       // Channel to signal worker to stop
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
			store:  store,
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
		case job := <-w.jobs:
			switch job.Type {
			case utils.FetchJob: // Specify the log level
				// Fetch logs from the store
				fetchedLogs, err := w.store.GetLogMessages(job.StartTime, job.EndTime, job.LogLevel)
				// Send the fetched logs back via the Result channel
				if err != nil {
					fmt.Println(err)
					job.Result <- nil
					continue
				}

				job.Result <- fetchedLogs

			case utils.StoreJob:
				w.store.InsertLogMessages(job.Logs)
			}

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

func (wp *WorkerPool) AddJob(job utils.Job) {
	wp.jobs <- job
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
