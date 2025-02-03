package retry

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	ctrl "sigs.k8s.io/controller-runtime"
)

var log = ctrl.Log.WithName("retry")

// RetryOptions configures the retry behavior
type RetryOptions struct {
	MaxRetries      int
	InitialInterval time.Duration
	MaxInterval     time.Duration
	BackoffFactor   float64
	RetryableErrors []error
}

// DefaultRetryOptions provides sensible defaults
var DefaultRetryOptions = RetryOptions{
	MaxRetries:      5,
	InitialInterval: 1 * time.Second,
	MaxInterval:     30 * time.Second,
	BackoffFactor:   2.0,
}

// Operation represents a retriable operation
type Operation func(context.Context) error

// IsRetryable checks if an error should be retried
func IsRetryable(err error, retryableErrors []error) bool {
	if err == nil {
		return false
	}

	// If no specific errors are defined, retry on any error
	if len(retryableErrors) == 0 {
		return true
	}

	// Check if the error matches any of the retryable errors
	for _, retryableErr := range retryableErrors {
		if err.Error() == retryableErr.Error() {
			return true
		}
	}

	return false
}

// WithRetry executes an operation with retry logic
func WithRetry(ctx context.Context, op Operation, opts RetryOptions) error {
	backoff := wait.Backoff{
		Steps:    opts.MaxRetries,
		Duration: opts.InitialInterval,
		Factor:   opts.BackoffFactor,
		Jitter:   0.1,
		Cap:      opts.MaxInterval,
	}

	var lastErr error
	var attempt int

	err := wait.ExponentialBackoffWithContext(ctx, backoff, func(ctx context.Context) (bool, error) {
		attempt++

		// Execute the operation
		if err := op(ctx); err != nil {
			lastErr = err

			// Check if error is retryable
			if !IsRetryable(err, opts.RetryableErrors) {
				return false, err // Stop retrying
			}

			log.Info("Operation failed, retrying",
				"attempt", attempt,
				"maxRetries", opts.MaxRetries,
				"error", err)

			return false, nil // Continue retrying
		}

		return true, nil // Success, stop retrying
	})

	if err != nil {
		return fmt.Errorf("operation failed after %d attempts: %v", attempt, lastErr)
	}

	return nil
}
