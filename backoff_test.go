package retrygo_test

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/ic-it/retrygo"
)

// Test LimitCount
func TestLimitCount(t *testing.T) {
	type zero struct{}
	const countLimit = 3

	retry, _ := retrygo.New[zero](
		retrygo.LimitCount(countLimit),
	)
	retryCount := 0
	_, err := retry.Do(context.TODO(), func(context.Context) (zero, error) {
		retryCount++
		return zero{}, fmt.Errorf("error")
	})
	if err == nil {
		t.Error("expected error")
	}
	if retryCount != countLimit {
		t.Errorf("expected %d retries, got %d", countLimit, retryCount)
	}
}

// Test LimitTime
func TestLimitTime(t *testing.T) {
	type zero struct{}
	const testTolerance = 100 * time.Millisecond
	const timeLimit = 3 * time.Second

	retry, _ := retrygo.New[zero](
		retrygo.LimitTime(timeLimit),
	)
	start := time.Now()
	_, err := retry.Do(context.TODO(), func(context.Context) (zero, error) {
		time.Sleep(time.Second)
		return zero{}, fmt.Errorf("error")
	})
	end := time.Now()
	if err == nil {
		t.Error("expected error")
	}
	duration := end.Sub(start).Milliseconds()
	if math.Abs(float64(duration-timeLimit.Milliseconds())) > float64(testTolerance.Milliseconds()) {
		t.Errorf("expected %s(+/-%s), got %s", timeLimit, testTolerance, end.Sub(start))
	}
}

// Test Combine
func TestCombine(t *testing.T) {
	type zero struct{}
	const countLimit = 3
	const timeLimit = 3 * time.Second

	retry, _ := retrygo.New[zero](
		retrygo.Combine(
			retrygo.LimitCount(countLimit),
			retrygo.LimitCount(5),
			retrygo.LimitTime(timeLimit),
		),
	)
	retryCount := 0
	_, err := retry.Do(context.TODO(), func(context.Context) (zero, error) {
		retryCount++
		time.Sleep(time.Second)
		return zero{}, fmt.Errorf("error")
	})
	if err == nil {
		t.Error("expected error")
	}
	if retryCount != countLimit {
		t.Errorf("expected %d retries, got %d", countLimit, retryCount)
	}
}

// Test Constant
func TestConstant(t *testing.T) {
	type zero struct{}
	const testTolerance = 100 * time.Millisecond
	const interval = 2 * time.Second // 2 + 2 = 4 seconds
	const expectedDuration = 4 * time.Second
	const countLimit = 3 // Due to the fact that the last retry doesn't do sleep

	retry, _ := retrygo.New[zero](
		retrygo.Combine(
			retrygo.Constant(interval),
			retrygo.LimitCount(countLimit),
		),
	)
	retryCount := 0
	start := time.Now()
	_, err := retry.Do(context.TODO(), func(context.Context) (zero, error) {
		retryCount++
		return zero{}, fmt.Errorf("error")
	})
	end := time.Now()
	if err == nil {
		t.Error("expected error")
	}
	if retryCount != countLimit {
		t.Errorf("expected %d retries, got %d", countLimit, retryCount)
	}
	duration := end.Sub(start).Milliseconds()
	if math.Abs(float64(duration-expectedDuration.Milliseconds())) > float64(testTolerance.Milliseconds()) {
		t.Errorf("expected %s(+/-%d), got %s", expectedDuration, testTolerance, end.Sub(start))
	}
}

// Test Linear
func TestLinear(t *testing.T) {
	type zero struct{}
	const testTolerance = 100 * time.Millisecond
	const interval = 2 * time.Second // 2 + 4 = 6 seconds
	const expectedDuration = 6 * time.Second
	const countLimit = 3 // Due to the fact that the last retry doesn't do sleep

	retry, _ := retrygo.New[zero](
		retrygo.Combine(
			retrygo.Linear(interval),
			retrygo.LimitCount(countLimit),
		),
	)
	retryCount := 0
	start := time.Now()
	_, err := retry.Do(context.TODO(), func(context.Context) (zero, error) {
		retryCount++
		return zero{}, fmt.Errorf("error")
	})
	end := time.Now()
	if err == nil {
		t.Error("expected error")
	}
	if retryCount != countLimit {
		t.Errorf("expected %d retries, got %d", countLimit, retryCount)
	}
	duration := end.Sub(start).Milliseconds()
	if math.Abs(float64(duration-expectedDuration.Milliseconds())) > float64(testTolerance.Milliseconds()) {
		t.Errorf("expected %s(+/-%s), got %s", expectedDuration, testTolerance, end.Sub(start))
	}
}

// Test Exponential
func TestExponential(t *testing.T) {
	type zero struct{}
	const testTolerance = 100 * time.Millisecond
	const interval = 1 * time.Second // 2 + 4 + 8 = 14 seconds
	const expectedDuration = 14 * time.Second
	const countLimit = 4 // Due to the fact that the last retry doesn't do sleep

	retry, _ := retrygo.New[zero](
		retrygo.Combine(
			retrygo.Exponential(interval),
			retrygo.LimitCount(countLimit),
		),
	)
	retryCount := 0
	start := time.Now()
	_, err := retry.Do(context.TODO(), func(context.Context) (zero, error) {
		retryCount++
		return zero{}, fmt.Errorf("error")
	})
	end := time.Now()
	if err == nil {
		t.Error("expected error")
	}
	if retryCount != countLimit {
		t.Errorf("expected %d retries, got %d", countLimit, retryCount)
	}
	duration := end.Sub(start).Milliseconds()
	if math.Abs(float64(duration-expectedDuration.Milliseconds())) > float64(testTolerance.Milliseconds()) {
		t.Errorf("expected %s(+/-%s), got %s", expectedDuration, testTolerance, end.Sub(start))
	}
}

// Test Jitter
func TestJitter(t *testing.T) {
	type zero struct{}
	const testTolerance = 4 * time.Second
	const interval = 2 * time.Second // 2(+[0,2)) + 2(+[0,2)) = 4(+[0,4)) seconds
	const expectedDuration = 4 * time.Second
	const countLimit = 3 // Due to the fact that the last retry doesn't do sleep

	retry, _ := retrygo.New[zero](
		retrygo.Combine(
			retrygo.Jitter(interval),
			retrygo.LimitCount(countLimit),
		),
	)
	retryCount := 0
	start := time.Now()
	_, err := retry.Do(context.TODO(), func(context.Context) (zero, error) {
		retryCount++
		return zero{}, fmt.Errorf("error")
	})
	end := time.Now()
	if err == nil {
		t.Error("expected error")
	}
	if retryCount != countLimit {
		t.Errorf("expected %d retries, got %d", countLimit, retryCount)
	}
	duration := end.Sub(start).Milliseconds()
	if math.Abs(float64(duration-expectedDuration.Milliseconds())) > float64(testTolerance.Milliseconds()) {
		t.Errorf("expected %s(+/-%s), got %s", expectedDuration, testTolerance, end.Sub(start))
	}
}
