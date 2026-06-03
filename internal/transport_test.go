package internal

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anexia/go-anxsdk/paging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//
// Helpers (kept minimal and explicit)
//

func newTestServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	return httptest.NewServer(handler)
}

//
// Tests
//

func TestTransport_Get_Success(t *testing.T) {
	// arrange
	require := require.New(t)
	assert := assert.New(t)

	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		// request assertions (contract validation)
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/v1/test", r.URL.Path)
		assert.Equal("Token test-key", r.Header.Get("Authorization"))
		assert.Equal("2", r.URL.Query().Get("foo"))

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"message": "ok",
			},
		})
	})
	defer ts.Close()

	tr := NewTransport(ts.URL, "test-key", ts.Client())

	var out struct {
		Data struct {
			Message string `json:"message"`
		} `json:"data"`
	}

	params := struct {
		Foo int `url:"foo"`
	}{
		Foo: 2,
	}

	// act
	err := tr.Get(context.Background(), "/v1/test", &out, paging.DefaultParams(), params)

	// assert
	require.NoError(err)
	assert.Equal("ok", out.Data.Message)
}

func TestTransport_Post_Success(t *testing.T) {
	// arrange
	require := require.New(t)
	assert := assert.New(t)

	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPost, r.Method)
		assert.Equal("application/json", r.Header.Get("Content-Type"))

		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)

		assert.Equal("value", body["key"])

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id": "123",
			},
		})
	})
	defer ts.Close()

	tr := NewTransport(ts.URL, "token", ts.Client())

	reqBody := struct {
		Key string `json:"key"`
	}{
		Key: "value",
	}

	var out struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	// act
	err := tr.Post(context.Background(), "/v1/create", reqBody, &out)

	// assert
	require.NoError(err)
	assert.Equal("123", out.Data.ID)
}

func TestTransport_Do_APIError(t *testing.T) {
	// arrange
	require := require.New(t)
	assert := assert.New(t)

	ts := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("bad request"))
	})
	defer ts.Close()

	tr := NewTransport(ts.URL, "", ts.Client())

	var out map[string]any

	// act
	err := tr.Get(context.Background(), "/v1/error", &out, paging.DefaultParams(), nil)

	// assert
	require.Error(err)

	var apiErr *TransportError
	require.ErrorAs(err, &apiErr)

	assert.Equal(400, apiErr.StatusCode)
	assert.Equal("bad request", apiErr.Body)
}

func TestTransport_BuildRequestUrlWithPaging(t *testing.T) {
	// arrange
	require := require.New(t)

	tr := NewTransport("https://api.example.com", "", nil)

	// act
	url, err := tr.buildRequestURL("/v1/test", &paging.Params{
		Page:  3,
		Limit: 77,
	}, nil)

	// assert
	require.NoError(err)
	assert.Equal(t, "https://api.example.com/v1/test?limit=77&page=3", url)
}

func TestTransport_BuildRequestUrl_PageError(t *testing.T) {
	// arrange
	require := require.New(t)

	tr := NewTransport("https://api.example.com", "", nil)

	// act
	_, err := tr.buildRequestURL("/v1/test", &paging.Params{
		Page:  0,
		Limit: 77,
	}, nil)

	// assert
	require.Error(err)
}
