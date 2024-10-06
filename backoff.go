package retrygo

import (
	"math/rand"
	"time"
)

// Constant returns a RetryPolicy that always returns the same interval
// between retries.
//
// Sleep formula: interval
func Constant(interval time.Duration) RetryPolicy {
	return func(ri RetryInfo) (bool, time.Duration) {
		return true, interval
	}
}

// Linear returns a RetryPolicy that increases the interval between retries
// linearly.
//
// Sleep formula: interval * fails
func Linear(interval time.Duration) RetryPolicy {
	return func(ri RetryInfo) (bool, time.Duration) {
		return true, interval * time.Duration(ri.Fails)
	}
}

// Exponential returns a RetryPolicy that increases the interval between
// retries exponentially.
//
// Sleep formula: interval * 2^fails
func Exponential(interval time.Duration) RetryPolicy {
	return func(ri RetryInfo) (bool, time.Duration) {
		return true, interval * time.Duration(1<<uint(ri.Fails-1))
	}
}

// Jitter returns a RetryPolicy that adds a random non-negative jitter to the interval
// between retries.
//
// Sleep formula: interval + rand.Int63n(interval)
func Jitter(interval time.Duration) RetryPolicy {
	return func(ri RetryInfo) (bool, time.Duration) {
		return true, interval + time.Duration(rand.Int63n(int64(interval)))
	}
}
