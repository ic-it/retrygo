package benchmarks

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ic-it/retrygo"

	benchm "github.com/ic-it/retrygo/benchmarks"
)

const LIB_NAME = "icit-retrygo"

var (
	failsRange = []int{1, 10, 100, 1000}
)

// BenchmarkNew benchmarks the performance of creating a new retry.
func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		retrygo.New[benchm.Zero](
			func(ri retrygo.RetryInfo) (bool, time.Duration) {
				return ri.Fails < 3, 0
			},
		)
	}
}

// BenchmarkDo benchmarks the performance of the Do method with error.
func BenchmarkDo(b *testing.B) {
	ctx := context.Background()
	for _, fails := range failsRange {
		retry := retrygo.New[benchm.Zero](
			func(ri retrygo.RetryInfo) (bool, time.Duration) {
				return ri.Fails < fails, 0
			})
		b.Run(fmt.Sprintf("lib=%s/fails=%d", LIB_NAME, fails), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				retry.Do(ctx, func(context.Context) (benchm.Zero, error) {
					return benchm.Zero{}, benchm.ErrBench
				})
			}
		})
	}
}

// BenchmarkDoSuccess benchmarks the performance of the Do method with success.
func BenchmarkDoSuccess(b *testing.B) {
	ctx := context.Background()
	for _, fails := range failsRange {
		retry := retrygo.New[benchm.Zero](
			func(ri retrygo.RetryInfo) (bool, time.Duration) {
				return ri.Fails < fails, 0
			})

		b.Run(fmt.Sprintf("lib=%s/fails=%d", LIB_NAME, fails), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				retry.Do(ctx, func(context.Context) (benchm.Zero, error) {
					return benchm.Zero{}, nil
				})
			}
		})
	}
}

// BenchmarkNewDo benchmarks the performance of creating a new retry and calling Do.
// This is worse scenario because it creates a new retry for each iteration.
func BenchmarkNewDo(b *testing.B) {
	ctx := context.Background()
	for _, fails := range failsRange {
		b.Run(fmt.Sprintf("lib=%s/fails=%d", LIB_NAME, fails), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				retry := retrygo.New[benchm.Zero](
					func(ri retrygo.RetryInfo) (bool, time.Duration) {
						return ri.Fails < fails, 0
					})
				retry.Do(ctx, func(context.Context) (benchm.Zero, error) {
					return benchm.Zero{}, benchm.ErrBench
				})
			}
		})
	}
}

// BenchmarkDoRecovery benchmarks the performance of the Do method with recovery.
func BenchmarkDoRecovery(b *testing.B) {
	for _, fails := range failsRange {
		retry := retrygo.New(
			func(ri retrygo.RetryInfo) (bool, time.Duration) {
				return ri.Fails < fails, 0
			},
			retrygo.WithRecovery[benchm.Zero](),
		)
		b.Run(fmt.Sprintf("lib=%s/fails=%d", LIB_NAME, fails), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				retry.Do(context.Background(), func(context.Context) (benchm.Zero, error) {
					panic("error")
				})
			}
		})
	}
}
