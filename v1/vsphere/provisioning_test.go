package vsphere

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/v1/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	return httptest.NewServer(handler)
}

func TestProvisioningClient_Deprovision(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodDelete, r.Method)
		assert.Equal("/api/vsphere/v1/provisioning/vm.json/vm-identifier", r.URL.Path)
		assert.Equal("true", r.URL.Query().Get("delayed"))

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"identifier":                 "vm-identifier",
			"delete_will_be_executed_at": "2026-01-01T00:00:00Z",
		})
	})
	defer ts.Close()

	client := newProvisioningClient(internal.NewTransport(ts.URL, "test-key", ts.Client()))

	resp, err := client.Deprovision(context.Background(), "vm-identifier", true)

	require.NoError(err)
	assert.Equal("vm-identifier", resp.Identifier)
	assert.Equal("2026-01-01T00:00:00Z", resp.DeleteWillBeExecutedAt)
}

func TestProvisioningClient_Deprovision_NotFound(t *testing.T) {
	require := require.New(t)

	ts := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("not found"))
	})
	defer ts.Close()

	client := newProvisioningClient(internal.NewTransport(ts.URL, "test-key", ts.Client()))

	_, err := client.Deprovision(context.Background(), "vm-identifier", false)

	require.Error(err)
	require.True(common.IsNotFoundError(err))
}
