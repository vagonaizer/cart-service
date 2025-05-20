package client

import (
	"net/http"
	"time"
)

// RetryMiddleware wraps an http.RoundTripper with retry logic
type RetryMiddleware struct {
	next       http.RoundTripper
	maxRetries int
	backoff    time.Duration
}

// NewRetryMiddleware creates a new retry middleware
func NewRetryMiddleware(next http.RoundTripper, maxRetries int, backoff time.Duration) *RetryMiddleware {
	return &RetryMiddleware{
		next:       next,
		maxRetries: maxRetries,
		backoff:    backoff,
	}
}

// RoundTrip implements http.RoundTripper
func (m *RetryMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := 0; i <= m.maxRetries; i++ {
		resp, err = m.next.RoundTrip(req)
		if err != nil {
			return nil, err
		}

		// Only retry on 420 or 429 status codes
		if resp.StatusCode != 420 && resp.StatusCode != 429 {
			return resp, nil
		}

		// If this was the last retry, return the response
		if i == m.maxRetries {
			return resp, nil
		}

		// Wait before retrying
		time.Sleep(m.backoff)
	}

	return resp, nil
}
