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
	const interval = 2 * time.Second

	requiredValues := []time.Duration{2 * time.Second, 2 * time.Second, 2 * time.Second, 2 * time.Second}

	backoff := retrygo.Constant(interval)
	info := retrygo.RetryInfo{Fails: 1}
	for _, expectedValue := range requiredValues {
		_, sleep := backoff(info)
		if sleep != expectedValue {
			t.Errorf("expected %s, got %s", expectedValue, sleep)
		}
		info.Fails++
	}
}

// Test Linear
func TestLinear(t *testing.T) {
	const interval = 2 * time.Second

	requiredValues := []time.Duration{2 * time.Second, 4 * time.Second, 6 * time.Second, 8 * time.Second}

	backoff := retrygo.Linear(interval)
	info := retrygo.RetryInfo{Fails: 1}
	for _, expectedValue := range requiredValues {
		_, sleep := backoff(info)
		if sleep != expectedValue {
			t.Errorf("expected %s, got %s", expectedValue, sleep)
		}
		info.Fails++
	}
}

// Test Exponential
func TestExponential(t *testing.T) {
	const interval = 2 * time.Second

	requiredValues := []time.Duration{2 * time.Second, 4 * time.Second, 8 * time.Second, 16 * time.Second}

	backoff := retrygo.Exponential(interval)
	info := retrygo.RetryInfo{Fails: 1}
	for _, expectedValue := range requiredValues {
		_, sleep := backoff(info)
		if sleep != expectedValue {
			t.Errorf("expected %s, got %s", expectedValue, sleep)
		}
		info.Fails++
	}
}

// Test Jitter
func TestJitter(t *testing.T) {
	const interval = 2 * time.Second

	backoff := retrygo.Jitter(interval)
	info := retrygo.RetryInfo{Fails: 1}
	for i := 0; i < 100; i++ {
		_, sleep := backoff(info)
		if sleep < interval || sleep > interval*2 {
			t.Errorf("expected %s to %s, got %s", interval, interval*2, sleep)
		}
		info.Fails++
	}
}
