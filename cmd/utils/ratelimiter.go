package utils

import (
	"sync"
	"time"
)

type RateLimiter struct {
	requests map[string]time.Time
	mu       sync.Mutex
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		requests: make(map[string]time.Time),
	}
}

func (rl *RateLimiter) Allow(domain string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if lastReq, exists := rl.requests[domain]; exists {
		if time.Since(lastReq) < time.Second {
			return false
		}
	}

	rl.requests[domain] = time.Now()
	return true
}
