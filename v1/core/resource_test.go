package core

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anexia/go-anxsdk/internal"
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

func TestResourceClient_GetTags_UsesIdentifier(t *testing.T) {
	// arrange
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/core/v1/resource.json/resource-identifier/tags", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode([]map[string]any{
			{
				"identifier": "tag-identifier",
				"name":       "tag-name",
			},
		})
	})
	defer ts.Close()

	client := newResourceClient(internal.NewTransport(ts.URL, "test-key", ts.Client()))

	// act
	tags, err := client.GetTags(context.Background(), "resource-identifier")

	// assert
	require.NoError(t, err)
	require.Len(t, tags, 1)
	assert.Equal(t, "tag-identifier", tags[0].Identifier)
}
