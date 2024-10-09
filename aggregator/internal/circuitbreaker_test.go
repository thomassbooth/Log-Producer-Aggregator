package internal_test

import (
	"fmt"
	"log-aggregator/aggregator/internal"
	"testing"
	"time"
)

// TestCircuitBreaker_InitialState tests the initial state of the CircuitBreaker.
func TestCircuitBreaker_InitialState(t *testing.T) {
	cb := internal.NewCircuitBreaker(3, 100*time.Millisecond)

	if cb.State() != "closed" {
		t.Errorf("Expected circuit breaker state to be 'closed', got '%s'", cb.State())
	}
}

// TestCircuitBreaker_OpenState tests the transition to the open state after failures.
func TestCircuitBreaker_OpenState(t *testing.T) {
	cb := internal.NewCircuitBreaker(3, 100*time.Millisecond)

	// Trigger failures to open the circuit breaker
	for i := 0; i < 3; i++ {
		err := cb.Call(func() error {
			return fmt.Errorf("simulated error")
		})
		if err == nil {
			t.Errorf("Expected an error, got none on attempt %d", i+1)
		}
	}

	// Check that the circuit breaker is now open
	if cb.State() != "open" {
		t.Errorf("Expected circuit breaker state to be 'open', got '%s'", cb.State())
	}
}

// TestCircuitBreaker_HalfOpenToClosed tests the transition from half-open to closed state.
func TestCircuitBreaker_HalfOpenToClosed(t *testing.T) {
	cb := internal.NewCircuitBreaker(3, 100*time.Millisecond)

	// Trigger failures to open the circuit breaker
	for i := 0; i < 3; i++ {
		err := cb.Call(func() error {
			return fmt.Errorf("simulated error")
		})
		if err == nil {
			t.Errorf("Expected an error, got none on attempt %d", i+1)
		}
	}

	// Check that the circuit breaker is now open
	if cb.State() != "open" {
		t.Errorf("Expected circuit breaker state to be 'open', got '%s'", cb.State())
	}

	// Wait for the reset timeout to pass
	time.Sleep(150 * time.Millisecond)

	// Test a successful call when the circuit breaker is half-open
	err := cb.Call(func() error {
		return nil // simulate successful operation
	})
	if err != nil {
		t.Errorf("Expected no error after timeout, got %v", err)
	}

	// The circuit breaker should return to closed state after a successful call
	if cb.State() != "closed" {
		t.Errorf("Expected circuit breaker state to be 'closed', got '%s'", cb.State())
	}
}

// TestCircuitBreaker_ResetsAfterSuccess tests that the circuit breaker resets after a successful call.
func TestCircuitBreaker_ResetsAfterSuccess(t *testing.T) {
	cb := internal.NewCircuitBreaker(3, 100*time.Millisecond)

	// Trigger failures to open the circuit breaker
	for i := 0; i < 3; i++ {
		err := cb.Call(func() error {
			return fmt.Errorf("simulated error")
		})
		if err == nil {
			t.Errorf("Expected an error, got none on attempt %d", i+1)
		}
	}

	// Check that the circuit breaker is now open
	if cb.State() != "open" {
		t.Errorf("Expected circuit breaker state to be 'open', got '%s'", cb.State())
	}

	// Wait for the reset timeout to pass
	time.Sleep(150 * time.Millisecond)

	// Test a successful call when the circuit breaker is half-open
	err := cb.Call(func() error {
		return nil // simulate successful operation
	})
	if err != nil {
		t.Errorf("Expected no error after timeout, got %v", err)
	}

	// Check the state after the successful call
	if cb.State() != "closed" {
		t.Errorf("Expected circuit breaker state to be 'closed', got '%s'", cb.State())
	}
}

// TestCircuitBreaker_OpenStateAfterFailure tests that the circuit breaker opens again after a failure in closed state.
func TestCircuitBreaker_OpenStateAfterFailure(t *testing.T) {
	cb := internal.NewCircuitBreaker(3, 100*time.Millisecond)

	// Trigger failures to open the circuit breaker
	for i := 0; i < 3; i++ {
		err := cb.Call(func() error {
			return fmt.Errorf("simulated error")
		})
		if err == nil {
			t.Errorf("Expected an error, got none on attempt %d", i+1)
		}
	}

	// Check that the circuit breaker is now open
	if cb.State() != "open" {
		t.Errorf("Expected circuit breaker state to be 'open', got '%s'", cb.State())
	}

	// Test a failed call to ensure it remains open
	err := cb.Call(func() error {
		return fmt.Errorf("simulated error")
	})
	if err == nil {
		t.Error("Expected an error, got none")
	}

	// Check the state remains open
	if cb.State() != "open" {
		t.Errorf("Expected circuit breaker state to be 'open', got '%s'", cb.State())
	}
}
