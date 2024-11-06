package ratelimiter

import (
	"fmt"
	"sync"
	"time"
)

// RateLimitData now only contains fields related to individual IP's rate limiting
type RateLimitData struct {
	IP              string
	LastAllowedTime time.Time
	RequestCount    int
	Blocked         bool
	BlockedUntil    time.Time
}

// InMemoryRateLimiter now contains MaxCount and Interval at the instance level
type InMemoryRateLimiter struct {
	store    map[string]*RateLimitData
	mu       sync.RWMutex  // Added mutex for thread safety
	MaxCount int           // Max allowed requests in the interval
	Interval time.Duration // Time window for rate limiting
}

// NewInMemoryRateLimiter creates a new instance of InMemoryRateLimiter with specified MaxCount and Interval
func NewInMemoryRateLimiter(maxCount int, interval time.Duration) *InMemoryRateLimiter {
	return &InMemoryRateLimiter{
		store:    make(map[string]*RateLimitData),
		MaxCount: maxCount,
		Interval: interval,
	}
}

// IsAllowed is the main entry point for rate limiting
func (rl *InMemoryRateLimiter) IsAllowed(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Get or initialize rate limit data
	data := rl.getRateLimitData(ip)
	return rl.processRequest(data)
}

// getRateLimitData handles first-time initialization more explicitly
func (rl *InMemoryRateLimiter) getRateLimitData(ip string) *RateLimitData {
	// Check if IP exists in store
	data, exists := rl.store[ip]
	if exists {
		return data
	}

	// First time seeing this IP - create new rate limit data
	now := time.Now()
	data = &RateLimitData{
		IP:              ip,
		LastAllowedTime: now,
		RequestCount:    0,
		Blocked:         false,
	}

	// Store the new data
	rl.store[ip] = data

	// Log first-time initialization
	fmt.Printf("New IP initialized: %s at %v\n", ip, now)

	return data
}

// processRequest handles subsequent requests
func (rl *InMemoryRateLimiter) processRequest(data *RateLimitData) bool {
	now := time.Now()
	// println(data)
	// Handle blocked IPs
	if data.Blocked {
		if now.Before(data.BlockedUntil) {
			fmt.Printf("IP %s is blocked until %v\n", data.IP, data.BlockedUntil)
			return false
		}
		data.Blocked = false
		fmt.Printf("IP %s block expired, resetting counts\n", data.IP)
	}

	// Check if we should reset the window
	if time.Since(data.LastAllowedTime) > rl.Interval {
		fmt.Printf("Resetting window for IP %s\n", data.IP)
		data.RequestCount = 0
		data.LastAllowedTime = now
	}

	// Increment request count
	data.RequestCount++

	// Check if limit exceeded
	if data.RequestCount > rl.MaxCount {
		data.Blocked = true
		data.BlockedUntil = now.Add(rl.Interval)
		fmt.Printf("IP %s exceeded limit (%d/%d). Blocked until %v\n",
			data.IP, data.RequestCount, rl.MaxCount, data.BlockedUntil)
		return false
	}

	fmt.Printf("Request allowed for IP %s (%d/%d)\n",
		data.IP, data.RequestCount, rl.MaxCount)
	return true
}

// GetStats returns statistics for an IP
func (rl *InMemoryRateLimiter) GetStats(ip string) string {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	data, exists := rl.store[ip]
	if !exists {
		return fmt.Sprintf("No data for IP: %s", ip)
	}

	return fmt.Sprintf(`
IP: %s
Request Count: %d/%d
Blocked: %v
Block Expires: %v
Time in Window: %v`,
		data.IP,
		data.RequestCount,
		rl.MaxCount,
		data.Blocked,
		data.BlockedUntil,
		time.Since(data.LastAllowedTime))
}
