package goretry

import (
    "fmt"
    "math"
    "time"
)

type RetryOptions struct {
    delayFactor      time.Duration
    randomizationFactor float64
    maxDelay         time.Duration
    maxAttempts      int
}

func (ro *RetryOptions) delay(attempt int) time.Duration {
    delayTime := time.Duration(math.Pow(2, float64(attempt)) * float64(ro.delayFactor) * (1 + ro.randomizationFactor*(math.Pow(2, float64(attempt))-1)))
    if delayTime > ro.maxDelay {
        delayTime = ro.maxDelay
    }
    return delayTime
}

func retry(fn func() error, ro RetryOptions, retryIf func(error) bool, onRetry func(int, error)) error {
    var err error
    for attempt := 1; attempt <= ro.maxAttempts; attempt++ {
        if err = fn(); err == nil {
            return nil
        } else if retryIf != nil && !retryIf(err) {
            return err
        }
        if onRetry != nil {
            onRetry(attempt, err)
        }
        time.Sleep(ro.delay(attempt))
    }
    return err
}

func main() {
    ro := RetryOptions{
        delayFactor:      200 * time.Millisecond,
        randomizationFactor: 0.25,
        maxDelay:         30 * time.Second,
        maxAttempts:      8,
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
    if err := retry(fn, ro, retryIf, onRetry); err != nil {
        fmt.Printf("Failed after %d attempts: %v\n", ro.maxAttempts, err)
    }
}
