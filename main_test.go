package goretry

import (
    "errors"
    "testing"
    "time"
)

func TestDelay(t *testing.T) {
    ro := RetryOptions{
        delayFactor:    200 * time.Millisecond,
        randomizationFactor: 0.25,
        maxDelay:         30 * time.Second,
        maxAttempts:      8,
    }
    
    expectedDelays := []time.Duration{
        200 * time.Millisecond,
        400 * time.Millisecond,
        800 * time.Millisecond,
        1600 * time.Millisecond,
        3200 * time.Millisecond,
        6400 * time.Millisecond,
        12000 * time.Millisecond,
        24000 * time.Millisecond,
    }
    
    for i := 0; i < len(expectedDelays); i++ {
        delay := ro.delay(i)
        if delay < expectedDelays[i] || delay > ro.maxDelay {
            t.Errorf("delay(%d) = %v, expected %v", i, delay, expectedDelays[i])
        }
    }
}

func TestRetry(t *testing.T) {
    ro := RetryOptions{
        delayFactor:      200 * time.Millisecond,
        randomizationFactor: 0.25,
        maxDelay:         30 * time.Second,
        maxAttempts:      8,
    }
    
    attempts := 0
    fn := func() error {
        attempts++
        if attempts < ro.maxAttempts {
            return errors.New("failed")
        }
        return nil
    }
    
    if err := retry(fn, ro, nil, nil); err != nil {
        t.Errorf("retry() returned error: %v", err)
    }
    
    attempts = 0
    fn = func() error {
        attempts++
        return errors.New("failed")
    }
    
    if err := retry(fn, ro, nil, nil); err == nil {
        t.Errorf("retry() did not return error")
    }
}
