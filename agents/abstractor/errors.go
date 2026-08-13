package abstractor

import (
	"errors"
	"fmt"
)

// Sentinel errors for provider failures.
var (
	ErrUnknownProvider     = errors.New("unknown provider")
	ErrConnectionRefused   = errors.New("provider connection refused")
	ErrTimeout             = errors.New("provider request timed out")
	ErrInvalidAPIKey       = errors.New("invalid or missing API key")
	ErrQuotaExceeded       = errors.New("provider quota exceeded")
	ErrProviderUnavailable = errors.New("provider temporarily unavailable")
	ErrEmptyResponse       = errors.New("provider returned empty response")
)

// ProviderError wraps a sentinel error with provider-specific context.
type ProviderError struct {
	Provider   string
	StatusCode int
	Message    string
	Err        error
}

// Error implements the error interface.
func (e *ProviderError) Error() string {
	if e.StatusCode > 0 {
		return fmt.Sprintf("%s provider error (HTTP %d): %s", e.Provider, e.StatusCode, e.Message)
	}
	return fmt.Sprintf("%s provider error: %s", e.Provider, e.Message)
}

// Unwrap supports errors.Is and errors.As.
func (e *ProviderError) Unwrap() error {
	return e.Err
}

// NewProviderError creates a ProviderError wrapping a sentinel error.
func NewProviderError(provider string, statusCode int, message string, sentinel error) *ProviderError {
	return &ProviderError{
		Provider:   provider,
		StatusCode: statusCode,
		Message:    message,
		Err:        sentinel,
	}
}

// ClassifyHTTPError maps an HTTP status code to the appropriate sentinel error.
func ClassifyHTTPError(provider string, statusCode int, body string) *ProviderError {
	switch {
	case statusCode == 401 || statusCode == 403:
		return NewProviderError(provider, statusCode, body, ErrInvalidAPIKey)
	case statusCode == 429:
		return NewProviderError(provider, statusCode, body, ErrQuotaExceeded)
	case statusCode == 503:
		return NewProviderError(provider, statusCode, body, ErrProviderUnavailable)
	case statusCode >= 500:
		return NewProviderError(provider, statusCode, body, ErrProviderUnavailable)
	default:
		return NewProviderError(provider, statusCode, body, fmt.Errorf("HTTP %d: %s", statusCode, body))
	}
}
