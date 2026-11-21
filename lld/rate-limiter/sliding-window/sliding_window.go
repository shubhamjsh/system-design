// Package slidingwindow implements an in-process exact sliding window rate limiter.
package main

import (
	"fmt"
	"sync"
	"time"
)

// RateLimiter allows at most limit requests in any trailing window. It is safe
// for concurrent use.
type RateLimiter struct {
	mu      sync.Mutex
	limit   int
	window  time.Duration
	allowed []time.Time
	now     func() time.Time
}

// New creates a limiter that allows up to limit requests in each trailing
// window. A non-positive limit or window creates a limiter that rejects every
// request.
func New(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		limit:  limit,
		window: window,
		now:    time.Now,
	}
}

// Allow reports whether a request can proceed, recording its time if allowed.
func (l *RateLimiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.limit <= 0 || l.window <= 0 {
		return false
	}

	now := l.now()
	cutoff := now.Add(-l.window)
	firstValid := 0
	for firstValid < len(l.allowed) && !l.allowed[firstValid].After(cutoff) {
		firstValid++
	}
	if firstValid > 0 {
		l.allowed = append(l.allowed[:0], l.allowed[firstValid:]...)
	}

	if len(l.allowed) >= l.limit {
		return false
	}
	l.allowed = append(l.allowed, now)
	return true
}

func main() {
	limiter := New(5, time.Second)

	for i := range 10 {
		if i > 7 && i < 9 {
			time.Sleep(time.Second)
		}
		fmt.Printf("Request %d: %v\n", i+1, limiter.Allow())
	}
}
