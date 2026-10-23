package vsphere

import "github.com/anexia/go-anxsdk/internal"

// PowerControlAPI is the interface implemented by [PowerControlClient].
// It is intentionally empty as power control is not implemented yet.
type PowerControlAPI any

// PowerControlClient is an api client for managing power control.
type PowerControlClient struct {
	transport *internal.Transport
}

// Check that PowerControlClient implements the PowerControlAPI interface.
var _ PowerControlAPI = (*PowerControlClient)(nil)

func newPowerControlClient(transport *internal.Transport) *PowerControlClient {
	return &PowerControlClient{
		transport: transport,
	}
}
