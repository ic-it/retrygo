package retrygo_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ic-it/retrygo"
)

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		retrygo.New[int](
			func(ri retrygo.RetryInfo) (bool, time.Duration) {
				return ri.Fails < 3, 0
			},
		)
	}
}

func TestDo(t *testing.T) {
	type zero struct{}
	// Create a retry instance with a mock RetryPolicy.
	retry, _ := retrygo.New[zero](
		func(ri retrygo.RetryInfo) (bool, time.Duration) {
			t.Log("retrying")
			return ri.Fails < 3, time.Duration(ri.Fails) * time.Second
		})

	// Call the Do method with a mock function that returns an error after 3 calls.
	_, err := retry.Do(context.Background(), func(context.Context) (zero, error) {
		return zero{}, fmt.Errorf("error")
	})

	// Check if the error is not nil.
	if err == nil {
		t.Error("expected error")
	}
}

func BenchmarkDo(b *testing.B) {
	err := fmt.Errorf("error")
	ctx := context.Background()
	for _, maxFails := range []int{1, 10, 100, 1000} {
		retry, _ := retrygo.New[int](
			func(ri retrygo.RetryInfo) (bool, time.Duration) {
				return ri.Fails < maxFails, 0
			})
		b.Run(fmt.Sprintf("maxFails=%d", maxFails), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				retry.Do(ctx, func(context.Context) (int, error) {
					return 0, err
				})
			}
		})
	}
}

func TestDoZero(t *testing.T) {
	// Create a retry instance with a mock RetryPolicy.
	retry, _ := retrygo.NewZero(
		func(ri retrygo.RetryInfo) (bool, time.Duration) {
			t.Log("retrying")
			return ri.Fails < 3, time.Duration(ri.Fails) * time.Second
		})

	// Call the DoZero method with a mock function that returns an error after 3 calls.
	err := retry.DoZero(context.Background(), func(context.Context) error {
		return fmt.Errorf("error")
	})

	// Check if the error is not nil.
	if err == nil {
		t.Error("expected error")
	}
}

func BenchmarkDoZero(b *testing.B) {
	err := fmt.Errorf("error")
	ctx := context.Background()
	for _, maxFails := range []int{1, 10, 100, 1000} {
		retry, _ := retrygo.NewZero(
			func(ri retrygo.RetryInfo) (bool, time.Duration) {
				return ri.Fails < maxFails, 0
			})
		b.Run(fmt.Sprintf("maxFails=%d", maxFails), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				retry.DoZero(ctx, func(context.Context) error {
					return err
				})
			}
		})
	}
}

func TestDoSuccess(t *testing.T) {
	// Create a retry instance with a mock RetryPolicy.
	retry, _ := retrygo.New[string](
		func(ri retrygo.RetryInfo) (bool, time.Duration) {
			t.Log("retrying")
			return ri.Fails < 3, time.Duration(ri.Fails) * time.Second
		})

	errcount := 0
	// Call the Do method with a mock function that returns nil.
	val, err := retry.Do(context.Background(), func(ctx context.Context) (string, error) {
		if errcount < 2 {
			errcount++
			return "", fmt.Errorf("error")
		}
		return "success", nil
	})

	// Check if the error is nil.
	if err != nil {
		t.Error("unexpected error")
	}

	// Check if the value is "success".
	if val != "success" {
		t.Error("unexpected value")
	}
}

func BenchmarkDoSuccess(b *testing.B) {
	ctx := context.Background()
	for _, maxFails := range []int{1, 10, 100, 1000} {
		retry, _ := retrygo.New[int](
			func(ri retrygo.RetryInfo) (bool, time.Duration) {
				return ri.Fails < maxFails, 0
			})
		b.Run(fmt.Sprintf("maxFails=%d", maxFails), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				retry.Do(ctx, func(context.Context) (int, error) {
					return 0, nil
				})
			}
		})
	}
}

func TestDoContextCancel(t *testing.T) {
	type zero struct{}
	// Create a retry instance with a mock RetryPolicy.
	retry, _ := retrygo.New[zero](
		func(ri retrygo.RetryInfo) (bool, time.Duration) {
			t.Log("retrying")
			return ri.Fails < 3, time.Duration(ri.Fails) * time.Second
		})

	// Create a context with a timeout of 1 second.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := retry.Do(ctx, func(ctx context.Context) (zero, error) {
		return zero{}, fmt.Errorf("error")
	})

	// Check if the error is not nil.
	if err == nil {
		t.Error("expected error")
	}

	// Check if the error is a context deadline exceeded error.
	if err != context.DeadlineExceeded {
		t.Error("expected context deadline exceeded error")
	}
}

func TestDoContextCancelBeforeDo(t *testing.T) {
	type zero struct{}
	// Create a retry instance with a mock RetryPolicy.
	retry, _ := retrygo.New[zero](
		func(ri retrygo.RetryInfo) (bool, time.Duration) {
			t.Log("retrying")
			return ri.Fails < 3, time.Duration(ri.Fails) * time.Second
		})

	// Create a context with a timeout of 1 second.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	cancel()

	_, err := retry.Do(ctx, func(ctx context.Context) (zero, error) {
		return zero{}, fmt.Errorf("error")
	})

	// Check if the error is not nil.
	if err == nil {
		t.Error("expected error")
	}

	// Check if the error is a context canceled error.
	if err != context.Canceled {
		t.Error("expected context canceled error")
	}
}

func TestDoContextCancelSuccess(t *testing.T) {
	// Create a retry instance with a mock RetryPolicy.
	retry, _ := retrygo.New[string](
		func(ri retrygo.RetryInfo) (bool, time.Duration) {
			t.Log("retrying")
			return ri.Fails < 3, time.Duration(ri.Fails) * time.Second
		})

	// Create a context with a timeout of 1 second.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Call the Do method with a mock function that returns nil.
	_, err := retry.Do(ctx, func(ctx context.Context) (string, error) {
		return "success", nil
	})

	// Check if the error is nil.
	if err != nil {
		t.Error("unexpected error")
	}
}

func TestDoMultipleTimes(t *testing.T) {
	ctx := context.Background()
	// Create a retry instance with a mock RetryPolicy.
	retry, _ := retrygo.New[string](
		func(ri retrygo.RetryInfo) (bool, time.Duration) {
			t.Log("retrying")
			return ri.Fails < 3, time.Duration(ri.Fails) * time.Second
		})

	// Call the Do method with a mock function that returns an error after 3 calls.
	_, err := retry.Do(ctx, func(ctx context.Context) (string, error) {
		return "", fmt.Errorf("error")
	})

	// Check if the error is not nil.
	if err == nil {
		t.Error("expected error")
	}

	// Call the Do method with a mock function that returns nil.
	val, err := retry.Do(ctx, func(ctx context.Context) (string, error) {
		return "success", nil
	})

	// Check if the error is nil.
	if err != nil {
		t.Error("unexpected error")
	}

	// Check if the value is "success".
	if val != "success" {
		t.Error("unexpected value")
	}
}

func BenchmarkNewDo(b *testing.B) {
	err := fmt.Errorf("error")
	ctx := context.Background()
	for _, maxFails := range []int{1, 10, 100, 1000} {
		b.Run(fmt.Sprintf("maxFails=%d", maxFails), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				retry, _ := retrygo.New[int](
					func(ri retrygo.RetryInfo) (bool, time.Duration) {
						return ri.Fails < maxFails, 0
					})
				retry.Do(ctx, func(context.Context) (int, error) {
					return 0, err
				})
			}
		})
	}
}

// Test recovery mode.
func TestDoRecovery(t *testing.T) {
	type zero struct{}
	// Create a retry instance with a mock RetryPolicy and recovery mode enabled.
	retry, _ := retrygo.New[zero](
		func(ri retrygo.RetryInfo) (bool, time.Duration) {
			t.Log("retrying", ri.Err)
			return ri.Fails < 3, time.Duration(ri.Fails) * time.Second
		},
		retrygo.WithRecovery[zero](),
	)

	// Call the Do method with a mock function that returns an error after 3 calls.
	_, err := retry.Do(context.Background(), func(context.Context) (zero, error) {
		panic("error")
	})

	// Check if the error is not nil.
	if err == nil {
		t.Error("expected error")
	}
}

func BenchmarkDoRecovery(b *testing.B) {
	for _, maxFails := range []int{1, 10, 100, 1000} {
		b.Run(fmt.Sprintf("maxFails=%d", maxFails), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				retry, _ := retrygo.New[int](
					func(ri retrygo.RetryInfo) (bool, time.Duration) {
						return ri.Fails < maxFails, 0
					},
					retrygo.WithRecovery[int](),
				)

				retry.Do(context.Background(), func(context.Context) (int, error) {
					panic("error")
				})
			}
		})
	}
}

func TestDoRecoverySuccess(t *testing.T) {
	// Create a retry instance with a mock RetryPolicy and recovery mode enabled.
	retry, _ := retrygo.New[string](
		func(ri retrygo.RetryInfo) (bool, time.Duration) {
			t.Log("retrying", ri.Err)
			return ri.Fails < 3, time.Duration(ri.Fails) * time.Second
		},
		retrygo.WithRecovery[string](),
	)

	errcount := 0
	// Call the Do method with a mock function that returns nil.
	val, err := retry.Do(context.Background(), func(context.Context) (string, error) {
		if errcount < 2 {
			errcount++
			panic("error")
		}
		return "success", nil
	})

	// Check if the error is nil.
	if err != nil {
		t.Error("unexpected error")
		return
	}

	// Check if the value is "success".
	if val != "success" {
		t.Error("unexpected value")
		return
	}
}

func TestDoRecoveryContextCancel(t *testing.T) {
	type zero struct{}
	// Create a retry instance with a mock RetryPolicy and recovery mode enabled.
	retry, _ := retrygo.New[zero](
		func(ri retrygo.RetryInfo) (bool, time.Duration) {
			t.Log("retrying", ri.Err)
			return ri.Fails < 3, time.Duration(ri.Fails) * time.Second
		},
		retrygo.WithRecovery[zero](),
	)

	// Create a context with a timeout of 1 second.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := retry.Do(ctx, func(context.Context) (zero, error) {
		panic("error")
	})

	// Check if the error is not nil.
	if err == nil {
		t.Error("expected error")
	}

	// Check if the error is a context deadline exceeded error.
	if err != context.DeadlineExceeded {
		t.Error("expected context deadline exceeded error")
	}
}

func TestDoRecoveryContextCancelBeforeDo(t *testing.T) {
	type zero struct{}
	// Create a retry instance with a mock RetryPolicy and recovery mode enabled.
	retry, _ := retrygo.New[zero](
		func(ri retrygo.RetryInfo) (bool, time.Duration) {
			t.Log("retrying", ri.Err)
			return ri.Fails < 3, time.Duration(ri.Fails) * time.Second
		},
		retrygo.WithRecovery[zero](),
	)

	// Create a context with a timeout of 1 second.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	cancel()

	_, err := retry.Do(ctx, func(context.Context) (zero, error) {
		panic("error")
	})

	// Check if the error is not nil.
	if err == nil {
		t.Error("expected error")
	}

	// Check if the error is a context canceled error.
	if err != context.Canceled {
		t.Error("expected context canceled error")
	}
}
