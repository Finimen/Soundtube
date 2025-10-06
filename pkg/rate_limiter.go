package pkg

import (
	"sort"
	"soundtube/pkg/config"
	"sync"
	"time"
)

type RateLimiter struct {
	mu          sync.Mutex
	requests    map[string][]time.Time
	maxRequests int
	window      time.Duration
}

func NewRateLimiter(cgf *config.RateLimiter) *RateLimiter {
	return &RateLimiter{
		requests:    make(map[string][]time.Time),
		maxRequests: cgf.MaxRequests,
		window:      time.Duration(cgf.Window),
	}
}

func (r *RateLimiter) Allow(ip string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-r.window)

	timestamps, exists := r.requests[ip]
	if !exists {
		timestamps = []time.Time{}
	}

	firstValidIndex := sort.Search(len(timestamps), func(i int) bool {
		return timestamps[i].After(windowStart)
	})
	validTimestamps := timestamps[firstValidIndex:]

	if len(validTimestamps) >= r.maxRequests {
		return false
	}

	validTimestamps = append(validTimestamps, now)
	r.requests[ip] = validTimestamps

	return true
}
