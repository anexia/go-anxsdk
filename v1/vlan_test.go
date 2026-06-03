package v1_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anexia/go-anxsdk/internal"
	v1 "github.com/anexia/go-anxsdk/v1"

	"github.com/stretchr/testify/require"
)

func newVlansClient(t *testing.T, ts *httptest.Server) *v1.VlansClient {
	t.Helper()
	tr := internal.NewTransport(ts.URL, "token", ts.Client())
	return v1.NewVlansClient(tr)
}

func TestVlansClient_Create(t *testing.T) {
	ts := RunAPITest(t, APITestCase{
		Method:  http.MethodPost,
		Path:    "/api/vlan/v1/vlan.json",
		Status:  http.StatusOK,
		Fixture: LoadFixture(t, "testdata/vlans/create_success.json"),

		ValidateRequest: ExpectJSONBody(v1.VlanCreateRequest{
			LocationIdentifier:  "loc-1",
			VMProvisioning:      true,
			DescriptionCustomer: "customer vlan",
		}),
	})

	client := newVlansClient(t, ts)

	resp, err := client.Create(context.Background(), v1.VlanCreateRequest{
		LocationIdentifier:  "loc-1",
		VMProvisioning:      true,
		DescriptionCustomer: "customer vlan",
	})

	require.NoError(t, err)
	require.Equal(t, "vlan-id", resp.Identifier)
	require.Equal(t, "vlan-name", resp.Name)
	require.Equal(t, "vlan description", resp.DescriptionCustomer)
}
