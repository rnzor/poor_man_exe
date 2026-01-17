package auth

import (
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	// 1 token per second, capacity 2
	rl := NewRateLimiter(1.0, 2.0)
	ip := "1.2.3.4"

	// Should allow first 2 requests immediately
	if !rl.Allow(ip) {
		t.Error("Should allow 1st request")
	}
	if !rl.Allow(ip) {
		t.Error("Should allow 2nd request")
	}

	// Should block 3rd request
	if rl.Allow(ip) {
		t.Error("Should block 3rd request")
	}

	// Wait 1 second (at least part of it)
	time.Sleep(1100 * time.Millisecond)

	// Should allow 1 more request
	if !rl.Allow(ip) {
		t.Error("Should allow request after recharge")
	}
}

func TestRateLimiterCleanup(t *testing.T) {
	rl := NewRateLimiter(1.0, 1.0)
	ip := "1.1.1.1"

	rl.Allow(ip)

	if len(rl.buckets) != 1 {
		t.Errorf("Expected 1 bucket, got %d", len(rl.buckets))
	}

	// Manually backdate the bucket to simulate old age
	rl.mu.Lock()
	rl.buckets[ip].lastUpdate = time.Now().Add(-11 * time.Minute)
	rl.mu.Unlock()

	rl.Cleanup()

	if len(rl.buckets) != 0 {
		t.Errorf("Expected 0 buckets after cleanup, got %d", len(rl.buckets))
	}
}
