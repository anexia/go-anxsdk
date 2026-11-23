package kubernetes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClustersClient_Delete(t *testing.T) {
	tests := []struct {
		env  Env
		path string
	}{
		{EnvDevelopment, "/api/kubernetes-dev/v2/cluster/cluster-id"},
		{EnvStaging, "/api/kubernetes-stage/v2/cluster/cluster-id"},
		{EnvProduction, "/api/kubernetes/v2/cluster/cluster-id"},
	}

	for _, tt := range tests {
		t.Run(string(tt.env), func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodDelete, r.Method)
				assert.Equal(t, tt.path, r.URL.Path)
				assert.Equal(t, "true", r.URL.Query().Get("skip_queued_automation_rules"))
				w.WriteHeader(http.StatusNoContent)
			}))
			t.Cleanup(srv.Close)

			c := newClustersClient(internal.NewTransport(srv.URL, "token", srv.Client()), tt.env)
			require.NoError(t, c.Delete(t.Context(), "cluster-id"))
		})
	}
}
