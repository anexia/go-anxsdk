package lbaas

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/v2/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListByReference(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/LBaaS/v2/backendswitchingrule", r.URL.Path)
		assert.Equal(t, "be", r.URL.Query().Get("backend"))
		assert.Empty(t, r.URL.Query().Get("frontend"))
		assert.Equal(t, "identifier,name,backend,frontend", r.URL.Query().Get("attributes"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"total_pages":1,"data":[{"identifier":"bsr","name":"ours","backend":{"identifier":"be"},"frontend":{"identifier":"fe"}}]}`))
	}))
	t.Cleanup(srv.Close)

	c := NewClient(internal.NewTransport(srv.URL, "token", srv.Client()))
	resp, err := c.ListByReference(t.Context(), paging.NewParams(1, 1000), ReferenceListParams{
		Resource:   "backendswitchingrule",
		RefField:   "backend",
		RefID:      "be",
		Attributes: "identifier,name,backend,frontend",
	})
	require.NoError(t, err)
	require.Equal(t, []ReferenceListItem{{
		Identifier: "bsr",
		Name:       "ours",
		Frontend:   common.Resource{Identifier: "fe"},
		Backend:    common.Resource{Identifier: "be"},
	}}, resp.Data)
	assert.Equal(t, "be", resp.Data[0].RelationID("backend"))

	_, err = c.ListByReference(t.Context(), paging.NewParams(1, 1000), ReferenceListParams{Resource: "rule", RefField: "name", RefID: "x"})
	require.ErrorContains(t, err, "unsupported lbaas reference field")
}

func TestListByReferenceWrappedPageAndStringRelation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"total_pages":1,"data":[{"identifier":"rule-1","name":"use backend","backend":"be","frontend":{"identifier":"fe"}}]}}`))
	}))
	t.Cleanup(srv.Close)

	c := NewClient(internal.NewTransport(srv.URL, "token", srv.Client()))
	resp, err := c.ListByReference(t.Context(), paging.NewParams(1, 100), ReferenceListParams{
		Resource:   "rule",
		RefField:   "backend",
		RefID:      "be",
		Attributes: "identifier,name,backend,frontend",
	})
	require.NoError(t, err)
	require.Equal(t, []ReferenceListItem{{
		Identifier: "rule-1",
		Name:       "use backend",
		Frontend:   common.Resource{Identifier: "fe"},
		Backend:    common.Resource{Identifier: "be"},
	}}, resp.Data)
}
