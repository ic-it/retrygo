package benchmarks

import (
	"fmt"
	"testing"

	"github.com/avast/retry-go/v4"

	benchm "github.com/ic-it/retrygo/benchmarks"
)

const LIB_NAME = "avast-retry-go"

var (
	failsRange = []uint{1, 10, 30}
)

func BenchmarkDo(b *testing.B) {
	for _, fails := range failsRange {
		b.Run(fmt.Sprintf("lib=%s/fails=%d", LIB_NAME, fails), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = retry.Do(
					func() error {
						return benchm.ErrBench
					},
					retry.Attempts(fails),
					retry.Delay(0),
				)
			}
		})
	}
}

func BenchmarkDoWithData(b *testing.B) {
	for _, fails := range failsRange {
		b.Run(fmt.Sprintf("lib=%s/fails=%d", LIB_NAME, fails), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = retry.DoWithData(
					func() (benchm.Zero, error) {
						return benchm.Zero{}, benchm.ErrBench
					},
					retry.Attempts(fails),
					retry.Delay(0),
				)
			}
		})
	}
}

func BenchmarkDoSuccess(b *testing.B) {
	for _, fails := range failsRange {
		b.Run(fmt.Sprintf("lib=%s/fails=%d", LIB_NAME, fails), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = retry.Do(
					func() error {
						return nil
					},
					retry.Attempts(fails),
					retry.Delay(0),
				)
			}
		})
	}
}

func BenchmarkDoWithDataSuccess(b *testing.B) {
	for _, fails := range failsRange {
		b.Run(fmt.Sprintf("lib=%s/fails=%d", LIB_NAME, fails), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = retry.DoWithData(
					func() (benchm.Zero, error) {
						return benchm.Zero{}, nil
					},
					retry.Attempts(fails),
					retry.Delay(0),
				)
			}
		})
	}
}
