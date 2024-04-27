package retrygo

import (
	"time"
)

// RetryInfo contains information about the retry
type RetryInfo struct {
	Fails int       // Fails is the number of retries
	Err   error     // Err is the error returned by the function
	Since time.Time // Since is the time when the retry started
}

// RetryPolicy is a function that returns a retry strategy based on the RetryInfo
type RetryPolicy func(RetryInfo) (continueRetry bool, sleep time.Duration)

// LimitCount returns a RetryPolicy that limits the number of retries.
//
// Sleep formula: 0
func LimitCount(count int) RetryPolicy {
	return func(ri RetryInfo) (bool, time.Duration) {
		return ri.Fails < count, 0
	}
}

// LimitTime returns a RetryPolicy that limits the total time spent on retries.
//
// Sleep formula: 0
//
// WARNING: Use context.WithTimeout instead of this function if you can!
func LimitTime(limit time.Duration) RetryPolicy {
	return func(ri RetryInfo) (bool, time.Duration) {
		return time.Since(ri.Since) < limit, 0
	}
}

// Combine returns a RetryPolicy that combines multiple RetryPolicies.
// The function will return true if all the policies return true. (logical AND)
// If all the policies return true, the function will return the sum of the sleep durations.
//
// Sleep formula: sleep1 + sleep2 + ...
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
