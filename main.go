package goretry

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"
)

type RetryOptions struct {
	DelayFactor         time.Duration
	RandomizationFactor float64
	MaxDelay            time.Duration
	MaxAttempts         int
	Timeout             time.Duration
}

func (ro *RetryOptions) Delay(attempt int) time.Duration {
	delayTime := time.Duration(math.Pow(2, float64(attempt)) * float64(ro.DelayFactor) * (1 + ro.RandomizationFactor*(rand.Float64()*2-1)))
	if delayTime > ro.MaxDelay {
		delayTime = ro.MaxDelay
	}
	return delayTime
}

func Retry(ctx context.Context, fn func() error, ro RetryOptions, retryIf func(error) bool, onRetry func(int, error), onTimeout func()) error {
	var err error
	attempt := 0
	startTime := time.Now()
	for {
		if err = fn(); err == nil {
			return nil
		} else if retryIf != nil && !retryIf(err) {
			return err
		}
		if onRetry != nil {
			onRetry(attempt+1, err)
		}
		attempt++
		if attempt >= ro.MaxAttempts {
			return err
		}
		delayTime := ro.Delay(attempt)
		select {
		case <-time.After(delayTime):
			// continue with the next retry
		case <-ctx.Done():
			return ctx.Err()
		}
		if ro.Timeout > 0 && time.Since(startTime) >= ro.Timeout {
			if onTimeout != nil {
				onTimeout()
			}
			return fmt.Errorf("retry timed out after %v", time.Since(startTime))
		}
	}
}

func main() {
	ro := RetryOptions{
		DelayFactor:         200 * time.Millisecond,
		RandomizationFactor: 0.25,
		MaxDelay:            30 * time.Second,
		MaxAttempts:         8,
		Timeout:             5 * time.Minute,
	}
	fn := func() error {
		fmt.Println("Trying...")
		return fmt.Errorf("failed")
	}
	retryIf := func(err error) bool {
		return true
	}
	onRetry := func(attempt int, err error) {
		fmt.Printf("Attempt %d failed: %v\n", attempt, err)
	}
	onTimeout := func() {
		fmt.Println("Retry timed out")
	}
	if err := Retry(context.Background(), fn, ro, retryIf, onRetry, onTimeout); err != nil {
		fmt.Printf("Failed after %d attempts: %v\n", ro.MaxAttempts, err)
	}
}
