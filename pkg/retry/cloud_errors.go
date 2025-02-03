package retry

import (
	"errors"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/aws/smithy-go"
	"google.golang.org/api/googleapi"
)

// Common retryable errors
var (
	ErrTimeout     = errors.New("timeout")
	ErrRateLimit   = errors.New("rate limit exceeded")
	ErrThrottling  = errors.New("throttling")
	ErrUnavailable = errors.New("service unavailable")
)

// IsAzureRetryable checks if an Azure error is retryable
func IsAzureRetryable(err error) bool {
	var azErr *azcore.ResponseError
	if errors.As(err, &azErr) {
		// Azure status codes that are typically retryable
		switch azErr.StatusCode {
		case 408, // Request Timeout
			429, // Too Many Requests
			500, // Internal Server Error
			502, // Bad Gateway
			503, // Service Unavailable
			504: // Gateway Timeout
			return true
		}
	}
	return false
}

// IsAWSRetryable checks if an AWS error is retryable
func IsAWSRetryable(err error) bool {
	var awsErr smithy.APIError
	if errors.As(err, &awsErr) {
		if strings.Contains(awsErr.ErrorCode(), "Throttling") ||
			strings.Contains(awsErr.ErrorCode(), "RequestLimitExceeded") ||
			strings.Contains(awsErr.ErrorCode(), "ServiceUnavailable") {
			return true
		}
	}
	return false
}

// IsGCPRetryable checks if a GCP error is retryable
func IsGCPRetryable(err error) bool {
	var gcpErr *googleapi.Error
	if errors.As(err, &gcpErr) {
		switch gcpErr.Code {
		case 408, // Request Timeout
			429, // Too Many Requests
			500, // Internal Server Error
			502, // Bad Gateway
			503, // Service Unavailable
			504: // Gateway Timeout
			return true
		}
	}
	return false
}
