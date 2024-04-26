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
