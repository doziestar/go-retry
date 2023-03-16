package goretry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	// Test with a function that always succeeds.
	t.Run("AlwaysSucceeds", func(t *testing.T) {
		fn := func() error {
			return nil
		}
		ro := RetryOptions{
			DelayFactor:         1 * time.Millisecond,
			RandomizationFactor: 0.0,
			MaxDelay:            1 * time.Millisecond,
			MaxAttempts:         5,
			Timeout:             100 * time.Millisecond,
		}
		err := Retry(context.Background(), fn, ro, nil, nil, nil)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	// Test with a function that always fails.
	t.Run("AlwaysFails", func(t *testing.T) {
		fn := func() error {
			return errors.New("always fails")
		}
		ro := RetryOptions{
			DelayFactor:         1 * time.Millisecond,
			RandomizationFactor: 0.0,
			MaxDelay:            1 * time.Millisecond,
			MaxAttempts:         5,
			Timeout:             100 * time.Millisecond,
		}
		err := Retry(context.Background(), fn, ro, nil, nil, nil)
		if err == nil {
			t.Error("expected an error, got nil")
		}
	})

	// Test with a function that fails until the last attempt.
	t.Run("FailsUntilLastAttempt", func(t *testing.T) {
		count := 0
		fn := func() error {
			if count < 4 {
				count++
				return errors.New("failed")
			}
			return nil
		}
		var lastAttempt int
		var lastError error
		onRetry := func(attempt int, err error) {
			lastAttempt = attempt
			lastError = err
		}
		ro := RetryOptions{
			DelayFactor:         1 * time.Millisecond,
			RandomizationFactor: 0.0,
			MaxDelay:            1 * time.Millisecond,
			MaxAttempts:         5,
			Timeout:             100 * time.Millisecond,
		}
		err := Retry(context.Background(), fn, ro, nil, onRetry, nil)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if lastAttempt != 5 {
			t.Errorf("expected last attempt to be 5, got %d", lastAttempt)
		}
		if lastError != nil {
			t.Errorf("expected last error to be nil, got %v", lastError)
		}
	})

	// Test with a custom retryIf function that retries on specific errors.
	t.Run("RetryOnSpecificErrors", func(t *testing.T) {
		count := 0
		fn := func() error {
			if count < 4 {
				count++
				return errors.New("retry on this error")
			}
			return errors.New("do not retry on this error")
		}
		retryIf := func(err error) bool {
			return err.Error() == "retry on this error"
		}
		ro := RetryOptions{
			DelayFactor:         1 * time.Millisecond,
			RandomizationFactor: 0.0,
			MaxDelay:            1 * time.Millisecond,
			MaxAttempts:         5,
			Timeout:             100 * time.Millisecond,
		}
		err := Retry(context.Background(), fn, ro, retryIf, nil, nil)
		if err == nil {
			t.Error("expected an error, got nil")
		}
		if err.Error() != "do not retry on this error" {
			t.Errorf("expected error to be 'do not retry on this error', got %v", err)
		}
	})

}

func TestRetrySuccess(t *testing.T) {
	ro := RetryOptions{
		DelayFactor:         100 * time.Millisecond,
		RandomizationFactor: 0.0,
		MaxDelay:            1 * time.Second,
		MaxAttempts:         5,
		Timeout:             5 * time.Second,
	}

	fn := func() error {
		return nil
	}

	err := Retry(context.Background(), fn, ro, nil, nil, nil)
	if err != nil {
		t.Errorf("Expected success, but got error: %v", err)
	}
}

func TestRetryMaxAttempts(t *testing.T) {
	ro := RetryOptions{
		DelayFactor:         100 * time.Millisecond,
		RandomizationFactor: 0.0,
		MaxDelay:            1 * time.Second,
		MaxAttempts:         5,
		Timeout:             5 * time.Second,
	}

	fn := func() error {
		return errors.New("failed")
	}

	err := Retry(context.Background(), fn, ro, nil, nil, nil)
	if err == nil {
		t.Error("Expected error, but got success")
	}

	if !errors.Is(err, fn()) {
		t.Errorf("Expected error %v, but got %v", fn(), err)
	}
}
