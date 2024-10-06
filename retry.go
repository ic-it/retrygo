package retrygo

import (
	"context"
	"time"
)

type zero struct{}

// Retry is the main type of this package.
type Retry[T any] struct {
	policy   RetryPolicy
	recovery bool
}

// type RetryOption[T any] func(*Retry[T])
type RetryOption[T any] interface {
	apply(*Retry[T]) error
}

// option[T any] is a type adapter for RetryOption[T].
type option[T any] struct {
	f func(*Retry[T]) error
}

func (o option[T]) apply(r *Retry[T]) error {
	return o.f(r)
}

// WithRecovery enables the recovery mode.
func WithRecovery[T any]() RetryOption[T] {
	return option[T]{
		f: func(r *Retry[T]) error {
			r.recovery = true
			return nil
		},
	}
}

// New creates a new Retry instance with the given RetryPolicy and RetryOptions.
func New[T any](policy RetryPolicy, options ...RetryOption[T]) (Retry[T], error) {
	r := Retry[T]{
		policy: policy,
	}
	for _, opt := range options {
		if err := opt.apply(&r); err != nil {
			return r, err
		}
	}
	return r, nil
}

// NewZero creates a new Retry instance with no return value.
func NewZero(policy RetryPolicy, options ...RetryOption[zero]) (Retry[zero], error) {
	return New(policy, options...)
}

// Do calls the given function f until it returns nil error or the context is done.
func (r Retry[T]) Do(ctx context.Context, f func(context.Context) (T, error)) (T, error) {
	if !r.recovery {
		return r.do(ctx, f)
	} else {
		return r.doRecovery(ctx, f)
	}
}

// DoZero calls the given function f until it returns nil error or the context is done.
func (r Retry[T]) DoZero(ctx context.Context, f func(context.Context) error) error {
	// fw is a function wrapper for f.
	fw := func(ctx context.Context) (T, error) {
		var zeroValue T
		return zeroValue, f(ctx)
	}
	_, err := r.Do(ctx, fw)
	return err
}

func (r Retry[T]) do(ctx context.Context, f func(context.Context) (T, error)) (T, error) {
	ri := RetryInfo{
		Fails: 0,
		Since: time.Now(),
		Err:   nil,
	}
	var result T
	timer := time.NewTimer(0)
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
		timer.Reset(sleep)
		select {
		case <-timer.C:
		case <-ctx.Done():
			return result, ctx.Err()
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
	timer := time.NewTimer(0)
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
		timer.Reset(sleep)
		select {
		case <-timer.C:
		case <-ctx.Done():
			timer.Stop()
			return result, ctx.Err()
		}
	}
}
