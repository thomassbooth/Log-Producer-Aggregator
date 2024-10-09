package internal

import (
	"fmt"
	"sync"
	"time"
)

// CircuitBreaker struct
type CircuitBreaker struct {
	mu               sync.Mutex
	failureCount     int
	failureThreshold int
	state            string // "closed", "open", "half-open"
	resetTimeout     time.Duration
	lastFailureTime  time.Time
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(threshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		failureThreshold: threshold,
		state:            "closed",
		resetTimeout:     timeout,
	}
}

// Call executes the given function and handles the circuit breaker logic
func (cb *CircuitBreaker) Call(fn func() error) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case "open":
		if time.Since(cb.lastFailureTime) < cb.resetTimeout {
			fmt.Println("Circuit breaker is open, rejecting call")
			return fmt.Errorf("circuit breaker is open")
		}
		// Switch to half-open state if timeout has passed
		cb.state = "half-open"
	case "half-open":
	}

	// Execute our function and check for errors
	if err := fn(); err != nil {
		cb.failureCount++
		cb.lastFailureTime = time.Now()
		fmt.Printf("Call failed, failure count: %d\n", cb.failureCount)
		if cb.failureCount >= cb.failureThreshold {
			cb.state = "open"
		}
		return err
	}

	// Successful call resets the circuit breaker
	cb.failureCount = 0
	if cb.state == "half-open" {
		cb.state = "closed" // Return to closed state
	}
	return nil
}

// State returns the current state of the circuit breaker
func (cb *CircuitBreaker) State() string {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}
