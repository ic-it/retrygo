package examples

import (
	"context"
	"time"

	"github.com/ic-it/retrygo"
)

func Simple() {
	type ReturnType struct{}
	retry, _ := retrygo.New[ReturnType](
		retrygo.Combine(
			retrygo.Constant(1*time.Second),
			retrygo.LimitCount(5),
			retrygo.LimitTime(10*time.Second),
		),
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
