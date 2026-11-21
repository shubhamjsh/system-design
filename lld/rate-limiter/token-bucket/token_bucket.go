// Package tokenbucket implements an in-process token bucket rate limiter.
package main

import (
	"fmt"
	"sync"
	"time"
)

// RateLimiter limits the rate of operations using a token bucket. It is safe
// for concurrent use. A token is added every refillInterval until the bucket
// reaches capacity.
type RateLimiter struct {
	mu             sync.Mutex
	capacity       int
	tokens         int
	refillInterval time.Duration
	lastRefill     time.Time
	now            func() time.Time
}

// New creates a limiter with a full bucket. capacity must be positive and
// refillInterval must be positive; invalid configurations create a limiter
// that rejects every request.
func New(capacity int, refillInterval time.Duration) *RateLimiter {
	if capacity <= 0 || refillInterval <= 0 {
		return &RateLimiter{}
	}

	now := time.Now()
	return &RateLimiter{
		capacity:       capacity,
		tokens:         capacity,
		refillInterval: refillInterval,
		lastRefill:     now,
		now:            time.Now,
	}
}

// Allow reports whether one request can proceed, consuming one token if so.
func (l *RateLimiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.capacity <= 0 || l.refillInterval <= 0 {
		return false
	}

	now := l.now()
	elapsed := now.Sub(l.lastRefill)
	if elapsed >= l.refillInterval {
		added := int(elapsed / l.refillInterval)
		if added > l.capacity-l.tokens {
			added = l.capacity - l.tokens
		}
		l.tokens += added
		l.lastRefill = l.lastRefill.Add(time.Duration(int64(elapsed/l.refillInterval)) * l.refillInterval)
	}

	if l.tokens == 0 {
		return false
	}
	l.tokens--
	return true
}

func main() {
	limiter := New(5, time.Second)

	for i := range 10 {
		fmt.Printf("Request %d: %v\n", i+1, limiter.Allow())
	}
}
