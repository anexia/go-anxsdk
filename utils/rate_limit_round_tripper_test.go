package utils

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//
// Helpers
//

func newTestServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	return httptest.NewServer(handler)
}

//
// Tests
//

func TestRateLimitRoundTripper_RetriesOn429(t *testing.T) {
	// arrange
	require := require.New(t)
	assert := assert.New(t)

	var requests int

	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		requests++

		if requests == 1 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	defer ts.Close()

	client := &http.Client{
		Transport: NewRateLimitRoundTripper(nil, time.Second, 3),
	}

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	assert.NoError(err)

	// act
	resp, err := client.Do(req)

	// assert
	require.NoError(err)
	require.NotNil(resp)
	defer resp.Body.Close()

	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(2, requests)
}

func TestRateLimitRoundTripper_MaxRetries(t *testing.T) {
	// arrange
	require := require.New(t)
	assert := assert.New(t)

	var requests int

	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		requests++

		w.Header().Set("Retry-After", "0")
		w.WriteHeader(http.StatusTooManyRequests)
	})
	defer ts.Close()

	client := &http.Client{
		Transport: NewRateLimitRoundTripper(nil, time.Second, 2),
	}

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	require.NoError(err)

	// act
	resp, err := client.Do(req)

	// assert
	require.NoError(err)
	require.NotNil(resp)
	defer resp.Body.Close()

	assert.Equal(http.StatusTooManyRequests, resp.StatusCode)
	assert.Equal(3, requests) // initial request + 2 retries
}

func TestRateLimitRoundTripper_DoesNotRetryNon429(t *testing.T) {
	// arrange
	require := require.New(t)
	assert := assert.New(t)

	var requests int

	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		requests++
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer ts.Close()

	client := &http.Client{
		Transport: NewRateLimitRoundTripper(nil, time.Second, 5),
	}

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	require.NoError(err)

	// act
	resp, err := client.Do(req)

	// assert
	require.NoError(err)
	require.NotNil(resp)
	_ = resp.Body.Close()

	assert.Equal(http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(1, requests)
}

func TestRateLimitRoundTripper_MaxRetryDuration(t *testing.T) {
	// arrange
	require := require.New(t)
	assert := assert.New(t)

	var requests int

	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		requests++

		w.Header().Set("Retry-After", "5")
		w.WriteHeader(http.StatusTooManyRequests)
	})
	defer ts.Close()

	client := &http.Client{
		Transport: NewRateLimitRoundTripper(nil, time.Second, 10),
	}

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	require.NoError(err)

	// act
	resp, err := client.Do(req)

	// assert
	require.NoError(err)
	require.NotNil(resp)
	_ = resp.Body.Close()

	assert.Equal(http.StatusTooManyRequests, resp.StatusCode)
	assert.Equal(1, requests)
}

func TestRateLimitRoundTripper_RetriesRequestBody(t *testing.T) {
	// arrange
	require := require.New(t)
	assert := assert.New(t)

	var requests int

	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		requests++

		body, err := io.ReadAll(r.Body)
		require.NoError(err)

		assert.Equal("hello", string(body))

		if requests == 1 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
	defer ts.Close()

	client := &http.Client{
		Transport: NewRateLimitRoundTripper(nil, time.Second, 3),
	}

	req, err := http.NewRequest(
		http.MethodPost,
		ts.URL,
		strings.NewReader("hello"),
	)
	require.NoError(err)

	// act
	resp, err := client.Do(req)

	// assert
	require.NoError(err)
	require.NotNil(resp)
	_ = resp.Body.Close()

	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(2, requests)
}
