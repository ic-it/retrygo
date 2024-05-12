package benchmarks

import (
	"context"
	"fmt"
	"testing"
	"time"

	benchm "github.com/ic-it/retrygo/benchmarks"
	"github.com/sethvargo/go-retry"
)

const LIB_NAME = "sethvargo-go-retry"

var (
	failsRange = []uint64{1, 10, 100, 1000}
)

var (
	errRetryable = retry.RetryableError(benchm.ErrBench)
)

func BenchmarkDo(b *testing.B) {
	ctx := context.Background()
	for _, fails := range failsRange {
		b.Run(fmt.Sprintf("lib=%s/fails=%d", LIB_NAME, fails), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bf := retry.WithMaxRetries(
					fails,
					retry.BackoffFunc(func() (time.Duration, bool) {
						return 0, false
					}),
				)
				_ = retry.Do(ctx, bf, func(ctx context.Context) error {
					return errRetryable
				})
			}
		})
	}
}

func BenchmarkDoSuccess(b *testing.B) {
	ctx := context.Background()
	for _, fails := range failsRange {
		b.Run(fmt.Sprintf("lib=%s/fails=%d", LIB_NAME, fails), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bf := retry.WithMaxRetries(
					fails,
					retry.BackoffFunc(func() (time.Duration, bool) {
						return 0, false
					}),
				)
				_ = retry.Do(ctx, bf, func(ctx context.Context) error {
					return nil
				})
			}
		})
	}
}
