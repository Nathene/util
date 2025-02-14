package ratelimiter

import (
	"sync"
	"time"
)

// RateLimiter controls how many requests can be processed per time unit.
type RateLimiter struct {
	mu        sync.Mutex
	tokens    int
	capacity  int
	rate      time.Duration
	lastCheck time.Time
}

// NewRateLimiter creates a new token bucket rate limiter.
func New(capacity int, refillRate time.Duration) *RateLimiter {
	return &RateLimiter{
		tokens:    capacity,
		capacity:  capacity,
		rate:      refillRate,
		lastCheck: time.Now(),
	}
}

// Allow checks if a request is allowed under the rate limit.
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastCheck)
	rl.lastCheck = now

	// Refill tokens
	rl.tokens += int(elapsed / rl.rate)
	if rl.tokens > rl.capacity {
		rl.tokens = rl.capacity
	}

	// Consume token if available
	if rl.tokens > 0 {
		rl.tokens--
		return true
	}
	return false
}

