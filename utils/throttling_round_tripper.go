package utils

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/time/rate"
)

// RateLimiter is an interface allowing rate limiting of outgoing requests.
type RateLimiter interface {
	// Wait returns nil if a token is taken before the context is done.
	Wait(ctx context.Context) error
}

// Check that the golang.org/x/time rate limiter can be used as a RateLimiter out of the box.
var _ RateLimiter = (*rate.Limiter)(nil)

// ThrottlingRoundTripper is a helper http.RoundTripper implementation that throttles requests using
// a client side [RateLimiter].
type ThrottlingRoundTripper struct {
	Next    http.RoundTripper
	limiter RateLimiter
}

// NewThrottlingRoundTripper creates a new ThrottlingRoundTripper instance. Passing no limiter (or a
// nil one) disables throttling.
func NewThrottlingRoundTripper(next http.RoundTripper, limiter RateLimiter) *ThrottlingRoundTripper {
	if next == nil {
		next = http.DefaultTransport
	}

	return &ThrottlingRoundTripper{
		Next:    next,
		limiter: limiter,
	}
}

// RoundTrip implements the http.RoundTripper interface for ThrottlingRoundTripper.
func (t *ThrottlingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := t.limiter.Wait(req.Context()); err != nil {
		return nil, fmt.Errorf("waiting for rate limit token: %w", err)
	}

	resp, err := t.Next.RoundTrip(req)
	if err != nil {
		return nil, fmt.Errorf("calling next: %w", err)
	}

	return resp, nil
}
