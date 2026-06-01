package utils

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// RateLimitRoundTripper is a helper http.RoundTripper implementation that automatically retries on rate limit api calls.
type RateLimitRoundTripper struct {
	Next             http.RoundTripper
	maxRetryDuration time.Duration
	maxRetries       int
}

// NewRateLimitRoundTripper creates a new RateLimitRoundTripper instance.
func NewRateLimitRoundTripper(
	next http.RoundTripper,
	maxRetryDuration time.Duration,
	maxRetries int,
) *RateLimitRoundTripper {
	if next == nil {
		next = http.DefaultTransport
	}

	return &RateLimitRoundTripper{
		Next:             next,
		maxRetryDuration: maxRetryDuration,
		maxRetries:       maxRetries,
	}
}

// ErrNoRewindableBody is returned during retry if the body cannot be read again.
var ErrNoRewindableBody = errors.New("cannot retry request with non-rewindable body")

// RoundTrip implements the http.RoundTripper interface for RateLimitRoundTripper.
func (t *RateLimitRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) { //nolint:revive // it's not that hard
	start := time.Now()
	retries := 0

	for {
		resp, err := t.Next.RoundTrip(req)
		if err != nil {
			return nil, fmt.Errorf("calling next: %w", err)
		}

		if resp.StatusCode != http.StatusTooManyRequests {
			return resp, nil
		}

		// Stop if we've reached the retry limit.
		if retries >= t.maxRetries {
			return resp, nil
		}

		retryAfter := time.Second

		if h := resp.Header.Get("Retry-After"); h != "" {
			if secs, err := strconv.Atoi(h); err == nil {
				retryAfter = time.Duration(secs) * time.Second
			} else if when, err := http.ParseTime(h); err == nil {
				retryAfter = time.Until(when)
				if retryAfter < 0 {
					retryAfter = 0
				}
			}
		}

		if time.Since(start)+retryAfter > t.maxRetryDuration {
			return resp, nil
		}

		_ = resp.Body.Close()

		time.Sleep(retryAfter)

		if req.Body != nil {
			if req.GetBody == nil {
				return nil, ErrNoRewindableBody
			}

			body, err := req.GetBody()
			if err != nil {
				return nil, fmt.Errorf("reading body again: %w", err)
			}

			req = req.Clone(req.Context())
			req.Body = body
		}

		retries++
	}
}
