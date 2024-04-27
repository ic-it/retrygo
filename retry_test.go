package retrygo_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ic-it/retrygo"
)

func TestDo(t *testing.T) {
	type zero struct{}
	// Create a retry instance with a mock RetryPolicy.
	retry := retrygo.New[zero](
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
	testerr := fmt.Errorf("error")
	for i := 0; i < b.N; i++ {
		retry := retrygo.New[int](
			func(ri retrygo.RetryInfo) (bool, time.Duration) {
				return ri.Fails < i, 0
			})

		retry.Do(context.Background(), func(context.Context) (int, error) {
			return 0, testerr
		})
	}
}

func TestDoSuccess(t *testing.T) {
	// Create a retry instance with a mock RetryPolicy.
	retry := retrygo.New[string](
		func(ri retrygo.RetryInfo) (bool, time.Duration) {
			t.Log("retrying")
			return ri.Fails < 3, time.Duration(ri.Fails) * time.Second
		})

	// Call the Do method with a mock function that returns nil.
	val, err := retry.Do(context.Background(), func(ctx context.Context) (string, error) {
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
	for i := 0; i < b.N; i++ {
		retry := retrygo.New[int](
			func(ri retrygo.RetryInfo) (bool, time.Duration) {
				return ri.Fails < i, 0
			})

		retry.Do(context.Background(), func(context.Context) (int, error) {
			return 0, nil
		})
	}
}

func TestDoContextCancel(t *testing.T) {
	type zero struct{}
	// Create a retry instance with a mock RetryPolicy.
	retry := retrygo.New[zero](
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

func BenchmarkDoContextCancel(b *testing.B) {
	testerr := fmt.Errorf("error")
	for i := 0; i < b.N; i++ {
		retry := retrygo.New[int](
			func(ri retrygo.RetryInfo) (bool, time.Duration) {
				return ri.Fails < i, 0
			})

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		retry.Do(ctx, func(context.Context) (int, error) {
			return 0, testerr
		})
	}
}

func TestDoContextCancelSuccess(t *testing.T) {
	// Create a retry instance with a mock RetryPolicy.
	retry := retrygo.New[string](
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

func BenchmarkDoContextCancelSuccess(b *testing.B) {
	for i := 0; i < b.N; i++ {
		retry := retrygo.New[string](
			func(ri retrygo.RetryInfo) (bool, time.Duration) {
				return ri.Fails < i, 0
			})

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		retry.Do(ctx, func(context.Context) (string, error) {
			return "success", nil
		})
	}
}

func TestDoMultipleTimes(t *testing.T) {
	ctx := context.Background()
	// Create a retry instance with a mock RetryPolicy.
	retry := retrygo.New[string](
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

func BenchmarkDoReuse(b *testing.B) {
	testerr := fmt.Errorf("error")
	var j *int = nil
	retry := retrygo.New[int](
		func(ri retrygo.RetryInfo) (bool, time.Duration) {
			return ri.Fails < *j, 0
		})

	for i := 0; i < b.N; i++ {
		j = &i
		retry.Do(context.Background(), func(context.Context) (int, error) {
			return 0, testerr
		})
	}
}
