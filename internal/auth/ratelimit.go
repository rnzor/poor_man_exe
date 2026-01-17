package auth

import (
	"sync"
	"time"
)

// RateLimiter implements a simple token bucket for IP-based rate limiting
type RateLimiter struct {
	mu       sync.Mutex
	buckets  map[string]*bucket
	rate     float64 // tokens per second
	capacity float64
}

type bucket struct {
	tokens     float64
	lastUpdate time.Time
}

func NewRateLimiter(rate float64, capacity float64) *RateLimiter {
	return &RateLimiter{
		buckets:  make(map[string]*bucket),
		rate:     rate,
		capacity: capacity,
	}
}

func (r *RateLimiter) Allow(ip string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	b, ok := r.buckets[ip]
	if !ok {
		b = &bucket{
			tokens:     r.capacity,
			lastUpdate: time.Now(),
		}
		r.buckets[ip] = b
	}

	now := time.Now()
	elapsed := now.Sub(b.lastUpdate).Seconds()
	b.tokens += elapsed * r.rate
	if b.tokens > r.capacity {
		b.tokens = r.capacity
	}
	b.lastUpdate = now

	if b.tokens >= 1.0 {
		b.tokens -= 1.0
		return true
	}

	return false
}

// Cleanup removes old buckets to save memory
func (r *RateLimiter) Cleanup() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for ip, b := range r.buckets {
		if now.Sub(b.lastUpdate) > 10*time.Minute {
			delete(r.buckets, ip)
		}
	}
}
