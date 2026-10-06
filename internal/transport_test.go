package internal

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
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
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		// request assertions (contract validation)
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/v1/test", r.URL.Path)
		assert.Equal(t, "Token test-key", r.Header.Get("Authorization"))
		assert.Equal(t, "2", r.URL.Query().Get("foo"))

		w.Header().Set("Content-Type", "application/json")
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
	require.NoError(t, err)
	assert.Equal(t, "ok", out.Data.Message)
}

func TestTransport_Get_SuccessWithAllAttributes(t *testing.T) {
	// arrange
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		// request assertions (contract validation)
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/v1/test", r.URL.Path)
		assert.Equal(t, "Token test-key", r.Header.Get("Authorization"))
		assert.Equal(t, "2", r.URL.Query().Get("foo"))
		assert.Equal(t, "all", r.URL.Query().Get("attributes"))

		w.Header().Set("Content-Type", "application/json")
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
	err := tr.Get(context.Background(), "/v1/test", &out, paging.DefaultParams(), NewAllAttributesWrapper(params))

	// assert
	require.NoError(t, err)
	assert.Equal(t, "ok", out.Data.Message)
}

func TestTransport_Post_Success(t *testing.T) {
	// arrange
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)

		assert.Equal(t, "value", body["key"])

		w.Header().Set("Content-Type", "application/json")
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
	require.NoError(t, err)
	assert.Equal(t, "123", out.Data.ID)
}

func TestTransport_DeleteWithResponse_Success(t *testing.T) {
	// arrange
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/v1/test/123", r.URL.Path)
		assert.Equal(t, "true", r.URL.Query().Get("delayed"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "123",
		})
	})
	defer ts.Close()

	tr := NewTransport(ts.URL, "test-key", ts.Client())

	var out struct {
		ID string `json:"id"`
	}

	params := struct {
		Delayed bool `url:"delayed"`
	}{
		Delayed: true,
	}

	// act
	err := tr.DeleteWithResponse(context.Background(), "/v1/test/123", params, &out)

	// assert
	require.NoError(t, err)
	assert.Equal(t, "123", out.ID)
}

func TestTransport_Do_APIError(t *testing.T) {
	// arrange
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
	require.Error(t, err)

	var apiErr *TransportError
	require.ErrorAs(t, err, &apiErr)

	assert.Equal(t, 400, apiErr.StatusCode)
	assert.Equal(t, "bad request", apiErr.Body)
}

func TestTransport_Do_WrongContentType(t *testing.T) {
	// arrange
	ts := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte(`{"text": "json but still wrong content type"}`))
	})
	defer ts.Close()

	tr := NewTransport(ts.URL, "", ts.Client())

	var out map[string]any

	// act
	err := tr.Get(context.Background(), "/v1/error", &out, paging.DefaultParams(), nil)

	// assert
	require.Error(t, err)
	assert.ErrorContains(t, err, "sending request: unsupported response Content-Type")
}

func TestTransport_BuildRequestUrlWithPaging(t *testing.T) {
	// arrange
	tr := NewTransport("https://api.example.com", "", nil)

	// act
	url, err := tr.buildRequestURL("/v1/test", &paging.Params{
		Page:  3,
		Limit: 77,
	}, nil)

	// assert
	require.NoError(t, err)
	assert.Equal(t, "https://api.example.com/v1/test?limit=77&page=3", url)
}

func TestTransport_BuildRequestUrl_PageError(t *testing.T) {
	// arrange
	tr := NewTransport("https://api.example.com", "", nil)

	// act
	_, err := tr.buildRequestURL("/v1/test", &paging.Params{
		Page:  0,
		Limit: 77,
	}, nil)

	// assert
	require.Error(t, err)
}

func TestTransport_JSONMarshallingFailed(t *testing.T) {
	// arrange
	ts := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	defer ts.Close()

	tr := NewTransport(ts.URL, "test-key", ts.Client())

	type doesNotSerializeJSON struct {
		Fail func()
	}

	reqBody := doesNotSerializeJSON{
		Fail: func() {},
	}

	// act
	err := tr.Post(context.Background(), "/v1/create", reqBody, nil)

	// assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "marshalling request")
}

//
// PreRequestValidator
//

// validatedRequest implements PreRequestValidator with a value receiver, so that both
// the value and a pointer to it satisfy the interface.
type validatedRequest struct {
	Key string `json:"key"`

	err   error
	calls *atomic.Int32
}

func (r validatedRequest) PreValidate() error {
	r.calls.Add(1)
	return r.err
}

// plainRequest deliberately does not implement PreRequestValidator.
type plainRequest struct {
	Key string `json:"key"`
}

// pointerValidatedRequest implements PreRequestValidator with a pointer receiver, mirroring
// v1/vsphere.ProvisioningRequest.
type pointerValidatedRequest struct {
	Key string `json:"key"`

	calls *atomic.Int32
}

func (r *pointerValidatedRequest) PreValidate() error {
	r.calls.Add(1)
	return nil
}

var errValidationFailed = errors.New("request is not valid")

func TestTransport_PreValidate_BlocksInvalidRequest(t *testing.T) {
	// arrange
	var serverHits atomic.Int32
	ts := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		serverHits.Add(1)
		w.WriteHeader(http.StatusOK)
	})
	defer ts.Close()

	tr := NewTransport(ts.URL, "test-key", ts.Client())

	var calls atomic.Int32
	reqBody := validatedRequest{Key: "value", err: errValidationFailed, calls: &calls}

	// act
	err := tr.Post(context.Background(), "/v1/create", reqBody, nil)

	// assert
	require.Error(t, err)
	require.ErrorIs(t, err, errValidationFailed)
	assert.Contains(t, err.Error(), "pre-request validation failed")
	assert.Equal(t, int32(1), calls.Load(), "PreValidate must be called exactly once")
	assert.Equal(t, int32(0), serverHits.Load(), "invalid request must not reach the server")
}

func TestTransport_PreValidate_SendsValidRequest(t *testing.T) {
	// arrange
	var serverHits atomic.Int32
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		serverHits.Add(1)

		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		assert.Equal(t, "value", body["key"])

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "123"})
	})
	defer ts.Close()

	tr := NewTransport(ts.URL, "test-key", ts.Client())

	var calls atomic.Int32
	reqBody := validatedRequest{Key: "value", calls: &calls}

	var out struct {
		ID string `json:"id"`
	}

	// act
	err := tr.Post(context.Background(), "/v1/create", reqBody, &out)

	// assert
	require.NoError(t, err)
	assert.Equal(t, "123", out.ID)
	assert.Equal(t, int32(1), calls.Load(), "PreValidate must be called exactly once")
	assert.Equal(t, int32(1), serverHits.Load())
}

// A request type that does not implement PreRequestValidator must be sent unchanged.
// This guards the type assertion in doRequest: an inverted check would call PreValidate
// on a nil interface here and panic.
func TestTransport_PreValidate_SkippedForNonValidator(t *testing.T) {
	// arrange
	var serverHits atomic.Int32
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		serverHits.Add(1)

		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		assert.Equal(t, "value", body["key"])

		w.WriteHeader(http.StatusOK)
	})
	defer ts.Close()

	tr := NewTransport(ts.URL, "test-key", ts.Client())

	// act / assert
	require.NotPanics(t, func() {
		err := tr.Post(context.Background(), "/v1/create", plainRequest{Key: "value"}, nil)
		require.NoError(t, err)
	})
	assert.Equal(t, int32(1), serverHits.Load())
}

func TestTransport_PreValidate_RunsForPut(t *testing.T) {
	// arrange
	ts := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	defer ts.Close()

	tr := NewTransport(ts.URL, "test-key", ts.Client())

	var calls atomic.Int32
	reqBody := validatedRequest{Key: "value", err: errValidationFailed, calls: &calls}

	// act
	err := tr.Put(context.Background(), "/v1/update", reqBody, nil)

	// assert
	require.ErrorIs(t, err, errValidationFailed)
	assert.Equal(t, int32(1), calls.Load())
}

// Requests without a body (GET, plain DELETE) must not attempt validation.
func TestTransport_PreValidate_SkippedForBodylessRequests(t *testing.T) {
	// arrange
	ts := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{})
	})
	defer ts.Close()

	tr := NewTransport(ts.URL, "test-key", ts.Client())

	// act / assert
	require.NotPanics(t, func() {
		require.NoError(t, tr.GetSingle(context.Background(), "/v1/test", nil))
		require.NoError(t, tr.Delete(context.Background(), "/v1/test/123"))
	})
}

// doRequest asserts against the dynamic type it is handed. A type whose PreValidate has a
// pointer receiver is only validated when a pointer is passed -- passing it by value
// silently skips validation. Callers must therefore pass such request bodies by pointer.
func TestTransport_PreValidate_PointerReceiverRequiresPointer(t *testing.T) {
	// arrange
	ts := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	defer ts.Close()

	tr := NewTransport(ts.URL, "test-key", ts.Client())

	// act: passed by value -- not part of the method set, so validation is skipped
	var valueCalls atomic.Int32
	err := tr.Post(context.Background(), "/v1/create", pointerValidatedRequest{Key: "value", calls: &valueCalls}, nil)

	// assert
	require.NoError(t, err)
	assert.Equal(t, int32(0), valueCalls.Load(), "value of a pointer-receiver type is not a PreRequestValidator")

	// act: passed by pointer -- validation runs
	var pointerCalls atomic.Int32
	err = tr.Post(context.Background(), "/v1/create", &pointerValidatedRequest{Key: "value", calls: &pointerCalls}, nil)

	// assert
	require.NoError(t, err)
	assert.Equal(t, int32(1), pointerCalls.Load(), "pointer to a pointer-receiver type is a PreRequestValidator")
}
