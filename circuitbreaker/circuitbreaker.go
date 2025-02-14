package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

const (
	open = "open"
	closed = "closed"
)

type CircuitBreaker struct {
	mu         sync.Mutex
	failures   int
	threshold  int
	timeout    time.Duration
	lastFail   time.Time
	state      string // "closed", "open", "half-open"
}

func New(threshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{threshold: threshold, timeout: timeout, state: closed}
}

func (cb *CircuitBreaker) Call(fn func() error) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.state == open && time.Since(cb.lastFail) < cb.timeout {
		return errors.New("circuit breaker is open")
	}

	err := fn()
	if err != nil {
		cb.failures++
		if cb.failures >= cb.threshold {
			cb.state = open
			cb.lastFail = time.Now()
		}
		return err
	}

	cb.state = closed
	cb.failures = 0
	return nil
}
