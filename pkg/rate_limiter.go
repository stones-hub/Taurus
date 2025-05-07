package pkg

import (
	"sync"
	"time"
)

// RateLimiter is a token bucket rate limiter.
type RateLimiter struct {
	capacity      int           // Maximum number of tokens
	tokens        int           // Current number of tokens
	fillInterval  time.Duration // Interval to add a token
	lastTokenTime time.Time     // Last time a token was added
	mutex         sync.Mutex    // Mutex to protect shared state
}

// NewRateLimiter creates a new RateLimiter.
func NewRateLimiter(capacity int, fillInterval time.Duration) *RateLimiter {
	return &RateLimiter{
		capacity:      capacity,
		tokens:        capacity,
		fillInterval:  fillInterval,
		lastTokenTime: time.Now(),
	}
}

// Allow checks if a request can proceed. It returns true if allowed, false otherwise.
func (rl *RateLimiter) Allow() bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastTokenTime)

	// Add tokens based on elapsed time
	tokensToAdd := int(elapsed / rl.fillInterval)
	if tokensToAdd > 0 {
		rl.tokens = min(rl.capacity, rl.tokens+tokensToAdd)
		rl.lastTokenTime = now
	}

	if rl.tokens > 0 {
		rl.tokens--
		return true
	}

	return false
}

// min returns the smaller of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
