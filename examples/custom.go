package examples

import (
	"context"
	"time"

	"github.com/ic-it/retrygo"
)

func Custom() {
	type ReturnType struct{}
	retry, _ := retrygo.New[ReturnType](
		func(ri retrygo.RetryInfo) (continueRetry bool, sleep time.Duration) {
			// Custom logic
			return false, 0
		},
	)

	val, err := retry.Do(context.TODO(), func(context.Context) (ReturnType, error) {
		// Do something
		return ReturnType{}, nil
	})

	if err != nil {
		// Handle error
		return
	}

	// Continue with val
	_ = val
}
