package retrygo

import (
	"math/rand"
	"time"
)

// Constant returns a RetryPolicy that always returns the same interval
// between retries.
//
// Formula: interval
func Constant(interval time.Duration) RetryPolicy {
	return func(ri RetryInfo) (bool, time.Duration) {
		return true, interval
	}
}

// Linear returns a RetryPolicy that increases the interval between retries
// linearly.
//
// Formula: interval * fails
func Linear(interval time.Duration) RetryPolicy {
	return func(ri RetryInfo) (bool, time.Duration) {
		return true, interval * time.Duration(ri.Fails)
	}
}

// Exponential returns a RetryPolicy that increases the interval between
// retries exponentially.
//
// Formula: interval * 2^fails
func Exponential(interval time.Duration) RetryPolicy {
	return func(ri RetryInfo) (bool, time.Duration) {
		return true, interval * time.Duration(1<<uint(ri.Fails))
	}
}

// Jitter returns a RetryPolicy that adds a random non-negative jitter to the interval
// between retries.
//
// Formula: interval + rand.Int63n(interval)
func Jitter(interval time.Duration) RetryPolicy {
	return func(ri RetryInfo) (bool, time.Duration) {
		return true, interval + time.Duration(rand.Int63n(int64(interval)))
	}
}

// LimitCount returns a RetryPolicy that limits the number of retries.
func LimitCount(count int) RetryPolicy {
	return func(ri RetryInfo) (bool, time.Duration) {
		return ri.Fails < count, 0
	}
}

// LimitTime returns a RetryPolicy that limits the total time spent on retries.
//
// WARNING: Use context.WithTimeout instead of this function if you can!
func LimitTime(limit time.Duration) RetryPolicy {
	return func(ri RetryInfo) (bool, time.Duration) {
		return time.Since(ri.Since) < limit, 0
	}
}

// Combine returns a RetryPolicy that combines multiple RetryPolicies.
// The resulting RetryPolicy will return false if any of the policies return false.
// But the resulting RetryPolicy will return the sum of the sleep times of all policies.
func Combine(policies ...RetryPolicy) RetryPolicy {
	return func(ri RetryInfo) (bool, time.Duration) {
		comulativeSleep := 0 * time.Second
		for _, policy := range policies {
			continueRetry, sleep := policy(ri)
			if !continueRetry {
				return false, 0
			}
			comulativeSleep += sleep
		}
		return true, comulativeSleep
	}
}
