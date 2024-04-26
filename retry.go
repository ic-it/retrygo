package retrygo

import (
	"context"
	"time"
)

type retry[RT any] struct {
	policy RetryPolicy
}

// Retry is the interface that wraps the Do method.
type Retry[RT any] interface {
	// Do calls the given function f until it returns nil or the context is done.
	Do(ctx context.Context, f func(context.Context) (RT, error)) (RT, error)
}

// New creates a new Retry instance with the given RetryPolicy and configurers.
func New[RT any](policy RetryPolicy) Retry[RT] {
	return &retry[RT]{
		policy: policy,
	}
}

func (r retry[RT]) Do(ctx context.Context, f func(context.Context) (RT, error)) (RT, error) {
	ri := RetryInfo{
		Fails: 0,
		Since: time.Now(),
		Err:   nil,
	}
	var result RT
	for {
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
		}
		result, err := f(ctx)
		if err == nil {
			return result, nil
		}
		ri.Err = err
		ri.Fails++
		continueRetry, sleep := r.policy(ri)
		if !continueRetry {
			return result, ri.Err
		}
		if sleep > 0 {
			select {
			case <-time.After(sleep):
			case <-ctx.Done():
				return result, ctx.Err()
			}
		}
	}
}
