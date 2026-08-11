package core

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPermissionGroupClient_Create_SendsOrganizationIdentifier(t *testing.T) {
	// arrange
	require := require.New(t)
	assert := assert.New(t)

	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPost, r.Method)

		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)

		assert.Equal("organization-identifier", body["organization_identifier"])
		assert.NotContains(body, "OrganizationIdentifier")

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"identifier": "group-identifier",
		})
	})
	defer ts.Close()

	client := newPermissionGroupClient(internal.NewTransport(ts.URL, "test-key", ts.Client()))

	// act
	resp, err := client.Create(context.Background(), PermissionGroupCreateRequest{
		Name:                   "group-name",
		OrganizationIdentifier: "organization-identifier",
	})

	// assert
	require.NoError(err)
	assert.Equal("group-identifier", resp.Identifier)
}

func TestPermissionGroupClient_Get_ReadsOrganizationIdentifier(t *testing.T) {
	// arrange
	require := require.New(t)
	assert := assert.New(t)

	ts := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"identifier":              "group-identifier",
			"name":                    "group-name",
			"organization_identifier": "organization-identifier",
		})
	})
	defer ts.Close()

	client := newPermissionGroupClient(internal.NewTransport(ts.URL, "test-key", ts.Client()))

	// act
	group, err := client.Get(context.Background(), "group-identifier")

	// assert
	require.NoError(err)
	assert.Equal("organization-identifier", group.OrganizationIdentifier)
}
