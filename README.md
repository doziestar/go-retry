# Go Retry Package

## Retry Function
The retry function retries a function that returns an error based on the specified RetryOptions until the maximum number of attempts is reached or the function succeeds.

## Usage
```go
func retry(fn func() error, ro RetryOptions, retryIf func(error) bool, onRetry func(int, error)) error
```

- `fn` is the function to retry. It must return an error.
- `ro` is a RetryOptions struct that contains the retry configuration.
- `retryIf` is an optional function that takes an error as input and returns a boolean. If `retryIf` is not nil and returns `false` for the error, the retry loop is terminated and the error is returned. If  `retryIf` is nil, all errors are retried.
- `onRetry` is an optional function that is called after each retry with the attempt number and error as inputs.

```go
type RetryOptions struct {
	DelayFactor         time.Duration
	RandomizationFactor float64
	MaxDelay            time.Duration
	MaxAttempts         int
	Timeout             time.Duration
}

```

- `delayFactor` is the base delay for retries.
- `randomizationFactor` is the amount of jitter to add to the delay.
- `maxDelay` is the maximum delay between retries.
- `maxAttempts` is the maximum number of attempts to make.
- `timeout` is the maximum amount of time to spend retrying.

### Functions
***`func (ro *RetryOptions) Delay(attempt int) time.Duration`*** 

This function calculates the delay time for the given retry attempt based on the retry options. It takes an integer representing the current attempt number and returns the delay time as a time.Duration value.

***`func Retry(ctx context.Context, fn func() error, ro RetryOptions, retryIf func(error) bool, onRetry func(int, error), onTimeout func()) error `***

This function performs the retries for the given function using the retry options provided. It takes the following arguments: 
- `ctx` is a context.Context that can be used to cancel the retry loop.
- `fn` is the function to retry. It must return an error.
- `ro` is a RetryOptions struct that contains the retry configuration.
- `retryIf` is an optional function that takes an error as input and returns a boolean. If `retryIf` is not nil and returns `false` for the error, the retry loop is terminated and the error is returned. If  `retryIf` is nil, all errors are retried.
- `onRetry` is an optional function that is called after each retry with the attempt number and error as inputs.
- `onTimeout` is an optional function that is called if the retry loop times out.


## Example
```go
package main

import (
    "fmt"
    "math"
    "time"
)

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
```

In this example, the `retry` function is called with a function that always returns an error. The `RetryOptions` are set to retry 8 times with a base delay of 200ms and a maximum delay of 30s. The `retryIf` function always returns true, so all errors are retried. The `onRetry` function prints the attempt number and error message after each retry. The output of the program is:
```
Trying...
Attempt 1 failed: failed
Trying...
Attempt 2 failed: failed
Trying...
Attempt 3 failed: failed
Trying...
Attempt 4 failed: failed
Trying...
Attempt 5 failed: failed
Trying...
Attempt 6 failed: failed
Trying...
Attempt 7 failed: failed
Trying...
Attempt 8 failed: failed
Failed after 8 attempts: failed
```