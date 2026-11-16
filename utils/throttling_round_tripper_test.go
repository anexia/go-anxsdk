package utils

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"
)

// stubRateLimiter is a controllable RateLimiter implementation for testing.
type stubRateLimiter struct {
	waitError error
	waits     int
}

func (s *stubRateLimiter) Wait(_ context.Context) error {
	s.waits++
	return s.waitError
}

// countingRoundTripper counts how often it is called.
type countingRoundTripper struct {
	count int
}

func (c *countingRoundTripper) RoundTrip(_ *http.Request) (*http.Response, error) {
	c.count++
	return &http.Response{StatusCode: http.StatusOK}, nil
}

// errThrottled is a sentinel error used to simulate a throttled request.
var errThrottled = errors.New("throttled")

//
// Tests
//

func TestThrottlingRoundTripper_DelegatesToNext(t *testing.T) {
	// arrange
	limiter := &stubRateLimiter{}

	client := &http.Client{
		Transport: NewThrottlingRoundTripper(nil, limiter),
	}

	// act
	resp, err := client.Get("https://example.com")

	// assert
	require.NoError(t, err)
	require.NotNil(t, resp)
	_ = resp.Body.Close()

	assert.Equal(t, 1, limiter.waits)
}

func TestThrottlingRoundTripper_ReturnsErrorWhenThrottled(t *testing.T) {
	// arrange
	next := &countingRoundTripper{}
	limiter := &stubRateLimiter{waitError: errThrottled}

	client := &http.Client{
		Transport: NewThrottlingRoundTripper(next, limiter),
	}

	// act
	resp, err := client.Get("https://example.com")

	// assert
	require.Error(t, err)
	if resp != nil {
		_ = resp.Body.Close()
	}
	assert.ErrorContains(t, err, "throttled")
	assert.Equal(t, 1, limiter.waits)
	assert.Equal(t, 0, next.count, "expected next round tripper not to be called when throttled")
}

func TestThrottlingRoundTripper_WithGolangRateLimiter(t *testing.T) {
	// arrange
	next := &countingRoundTripper{}
	limiter := rate.NewLimiter(rate.Limit(5), 10)

	client := &http.Client{
		Transport: NewThrottlingRoundTripper(next, limiter),
	}

	// act
	resp, err := client.Get("https://example.com")

	// assert
	require.NoError(t, err)
	require.NotNil(t, resp)
	_ = resp.Body.Close()

	assert.Equal(t, 1, next.count)
}
