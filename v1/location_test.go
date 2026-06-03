package v1_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/paging"
	v1 "github.com/anexia/go-anxsdk/v1"

	"github.com/stretchr/testify/require"
)

func newLocationsClient(t *testing.T, ts *httptest.Server) *v1.LocationsClient {
	t.Helper()
	tr := internal.NewTransport(ts.URL, "token", ts.Client())
	return v1.NewLocationsClient(tr)
}

func TestLocationsClient_List(t *testing.T) {
	// arrange
	ts := RunAPITest(t, APITestCase{
		Method:  http.MethodGet,
		Path:    "/api/core/v1/location.json",
		Status:  http.StatusOK,
		Fixture: LoadFixture(t, "testdata/locations/list_success.json"),

		ValidateRequest: func(t *testing.T, r *http.Request) {
			ExpectQuery(t, map[string]string{
				"page":   "2",
				"limit":  "50",
				"search": "vienna",
			})(r)
		},
	})

	client := newLocationsClient(t, ts)

	// act
	resp, err := client.List(context.Background(), paging.Params{Page: 2, Limit: 50}, v1.LocationListParams{
		Search: "vienna",
	})

	// assert
	require.NoError(t, err)

	require.Equal(t, 2, resp.Page)
	require.Equal(t, 50, resp.Limit)
	require.NotEmpty(t, resp.Data)

	first := resp.Data[0]
	require.Equal(t, "loc-1", first.Identifier)
	require.Equal(t, "VIE", first.Code)
	require.Equal(t, "Vienna", first.Name)
}
