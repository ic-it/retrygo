package retrygo

import (
	"context"
	"time"
)

// Retry is the interface that wraps the Do method.
type Retry[T any] struct {
	policy   RetryPolicy
	recovery bool
}

type retryConfigurer[T any] func(*Retry[T])

// WithRecovery enables the recovery mode.
func WithRecovery[T any]() retryConfigurer[T] {
	return func(r *Retry[T]) {
		r.recovery = true
	}
}

// New creates a new Retry instance with the given RetryPolicy and configurers.
func New[T any](policy RetryPolicy, configurers ...retryConfigurer[T]) Retry[T] {
	r := Retry[T]{
		policy: policy,
	}
	for _, c := range configurers {
		c(&r)
	}
	return r
}

// Do calls the given function f until it returns nil error or the context is done.
func (r Retry[T]) Do(ctx context.Context, f func(context.Context) (T, error)) (T, error) {
	if !r.recovery {
		return r.do(ctx, f)
	} else {
		return r.doRecovery(ctx, f)
	}
}

func (r Retry[T]) do(ctx context.Context, f func(context.Context) (T, error)) (T, error) {
	ri := RetryInfo{
		Fails: 0,
		Since: time.Now(),
		Err:   nil,
	}
	var result T
	for {
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
		}
		result, ri.Err = f(ctx)
		if ri.Err == nil {
			return result, nil
		}
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

func (r Retry[T]) doRecovery(ctx context.Context, f func(context.Context) (T, error)) (T, error) {
	ri := RetryInfo{
		Fails: 0,
		Since: time.Now(),
		Err:   nil,
	}
	var result T
	for {
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
		}
		func() {
			defer func() {
				if r := recover(); r != nil {
					ri.Err = ErrRecovered{V: r}
				}
			}()
			result, ri.Err = f(ctx)
		}()
		if ri.Err == nil {
			return result, nil
		}
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
