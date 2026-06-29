package vsphere

import "github.com/anexia/go-anxsdk/internal"

// PowerControlClient is an api client for managing power control.
type PowerControlClient struct {
	transport *internal.Transport
}

func newPowerControlClient(transport *internal.Transport) *PowerControlClient {
	return &PowerControlClient{
		transport: transport,
	}
}
