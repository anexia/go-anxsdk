package vsphere

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
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
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/api/vsphere/v1/provisioning/vm.json/vm-identifier", r.URL.Path)
		assert.Equal(t, "true", r.URL.Query().Get("delayed"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"identifier":                 "vm-identifier",
			"delete_will_be_executed_at": "2026-01-01T00:00:00Z",
		})
	})
	defer ts.Close()

	client := newProvisioningClient(internal.NewTransport(ts.URL, "test-key", ts.Client()))

	resp, err := client.Deprovision(context.Background(), "vm-identifier", true)

	require.NoError(t, err)
	assert.Equal(t, "vm-identifier", resp.Identifier)
	assert.Equal(t, "2026-01-01T00:00:00Z", resp.DeleteWillBeExecutedAt)
}

func TestProvisioningClient_Deprovision_NotFound(t *testing.T) {
	ts := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("not found"))
	})
	defer ts.Close()

	client := newProvisioningClient(internal.NewTransport(ts.URL, "test-key", ts.Client()))

	_, err := client.Deprovision(context.Background(), "vm-identifier", false)

	require.Error(t, err)
	require.True(t, common.IsNotFoundError(err))
}

// validRequest returns the minimal ProvisioningRequest that passes PreValidate.
// Every optional pointer field is left nil so that callers opting into a single
// field also exercise the nil handling of all the others.
func validRequest() ProvisioningRequest {
	return ProvisioningRequest{
		Hostname: "test-vm",
		SSH:      new("ssh-ed25519 AAAA"),
	}
}

func TestProvisioningRequest_PreValidate(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*ProvisioningRequest)
		wantErr error  // nil means the request must validate
		wantMsg string // substring the error message must contain, if any
	}{
		{
			name:   "minimal valid request with ssh",
			mutate: func(*ProvisioningRequest) {},
		},
		{
			name: "password instead of ssh",
			mutate: func(r *ProvisioningRequest) {
				r.SSH = nil
				r.Password = new("hunter2")
			},
		},
		{
			name: "no credentials at all",
			mutate: func(r *ProvisioningRequest) {
				r.SSH = nil
			},
			wantErr: ErrMissingCredentials,
		},
		{
			name: "credentials present but empty",
			mutate: func(r *ProvisioningRequest) {
				r.SSH = new("")
				r.Password = new("")
			},
			wantErr: ErrMissingCredentials,
		},
		{
			name: "valid base64 script",
			mutate: func(r *ProvisioningRequest) {
				r.Script = new(base64.StdEncoding.EncodeToString([]byte("#!/bin/sh\necho hi")))
			},
		},
		{
			name: "script that is not base64",
			mutate: func(r *ProvisioningRequest) {
				r.Script = new("not base64 !!!")
			},
			wantMsg: "Script is not base64 encoded",
		},
		{
			name: "empty script is not validated",
			mutate: func(r *ProvisioningRequest) {
				r.Script = new("")
			},
		},
		{
			name: "negative MemoryMB",
			mutate: func(r *ProvisioningRequest) {
				r.MemoryMB = new(-1024)
			},
			wantErr: ErrNegativeValue,
			wantMsg: "MemoryMB must not be negative: -1024",
		},
		{
			name: "negative CPUs",
			mutate: func(r *ProvisioningRequest) {
				r.CPUs = new(-2)
			},
			wantErr: ErrNegativeValue,
			wantMsg: "CPUs",
		},
		{
			name: "negative DiskGB",
			mutate: func(r *ProvisioningRequest) {
				r.DiskGB = new(-10)
			},
			wantErr: ErrNegativeValue,
			wantMsg: "DiskGB",
		},
		{
			name: "negative Sockets",
			mutate: func(r *ProvisioningRequest) {
				r.Sockets = new(-1)
			},
			wantErr: ErrNegativeValue,
			wantMsg: "Sockets",
		},
		{
			name: "negative BootDelaySeconds",
			mutate: func(r *ProvisioningRequest) {
				r.BootDelaySeconds = new(-5)
			},
			wantErr: ErrNegativeValue,
			wantMsg: "BootDelaySeconds",
		},
		{
			name: "zero values are accepted",
			mutate: func(r *ProvisioningRequest) {
				r.MemoryMB = new(0)
				r.CPUs = new(0)
				r.DiskGB = new(0)
				r.Sockets = new(0)
				r.BootDelaySeconds = new(0)
			},
		},
		{
			name: "valid ipv4 dns",
			mutate: func(r *ProvisioningRequest) {
				r.DNS1 = new("8.8.8.8")
			},
		},
		{
			name: "valid ipv6 dns",
			mutate: func(r *ProvisioningRequest) {
				r.DNS2 = new("2001:4860:4860::8888")
			},
		},
		{
			name: "empty dns is not validated",
			mutate: func(r *ProvisioningRequest) {
				r.DNS1 = new("")
			},
		},
		{
			name: "invalid DNS1",
			mutate: func(r *ProvisioningRequest) {
				r.DNS1 = new("not-an-ip")
			},
			wantErr: ErrInvalidIP,
			wantMsg: `DNS1 is not a valid IP: "not-an-ip"`,
		},
		{
			name: "invalid DNS4",
			mutate: func(r *ProvisioningRequest) {
				r.DNS4 = new("999.999.999.999")
			},
			wantErr: ErrInvalidIP,
			wantMsg: "DNS4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validRequest()
			tt.mutate(&req)

			err := req.PreValidate()

			if tt.wantErr == nil && tt.wantMsg == "" {
				require.NoError(t, err)
				return
			}

			require.Error(t, err)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			}
			if tt.wantMsg != "" {
				assert.Contains(t, err.Error(), tt.wantMsg)
			}
		})
	}
}

// PreValidate is called on every request body, so a zero-valued request with all
// optional pointers unset must report its problems rather than panic.
func TestProvisioningRequest_PreValidate_ZeroValueDoesNotPanic(t *testing.T) {
	req := ProvisioningRequest{}

	require.NotPanics(t, func() {
		err := req.PreValidate()
		require.ErrorIs(t, err, ErrMissingCredentials)
	})
}

func TestProvisioningRequest_PreValidate_ReportsAllViolations(t *testing.T) {
	req := ProvisioningRequest{
		Hostname: "test-vm",
		Script:   new("not base64 !!!"),
		MemoryMB: new(-1),
		CPUs:     new(-1),
		DNS1:     new("not-an-ip"),
		DNS3:     new("also-not-an-ip"),
	}

	err := req.PreValidate()

	require.Error(t, err)
	require.ErrorIs(t, err, ErrMissingCredentials)
	require.ErrorIs(t, err, ErrNegativeValue)
	require.ErrorIs(t, err, ErrInvalidIP)

	msg := err.Error()
	assert.Contains(t, msg, "provisioning request validation")
	for _, want := range []string{"Script is not base64 encoded", "MemoryMB", "CPUs", "DNS1", "DNS3"} {
		assert.Contains(t, msg, want)
	}

	// Unset optional fields must not contribute errors.
	for _, unwanted := range []string{"DiskGB", "Sockets", "BootDelaySeconds", "DNS2", "DNS4"} {
		assert.NotContains(t, msg, unwanted)
	}
}

func TestProvisioningRequest_PreValidate_SatisfiesPreRequestValidator(t *testing.T) {
	var _ internal.PreRequestValidator = &ProvisioningRequest{}

	req := validRequest()
	validator, ok := any(&req).(internal.PreRequestValidator)

	require.True(t, ok)
	require.NoError(t, validator.PreValidate())
}

func TestProvisioningRequest_PreValidate_WrappedErrorsAreUnwrappable(t *testing.T) {
	req := validRequest()
	req.MemoryMB = new(-1)

	err := req.PreValidate()

	require.Error(t, err)
	require.True(t, errors.Is(err, ErrNegativeValue))
	require.False(t, errors.Is(err, ErrInvalidIP))
}
