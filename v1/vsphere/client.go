package vsphere

import (
	"github.com/anexia/go-anxsdk/internal"
)

// Client is an anexia vsphere v1 api client.
type Client struct {
	transport *internal.Transport
}

// NewClient creates a new v1 vsphere api client.
func NewClient(transport *internal.Transport) *Client {
	return &Client{
		transport: transport,
	}
}

// PowerControl returns a power control client.
func (c *Client) PowerControl() *PowerControlClient {
	return newPowerControlClient(c.transport)
}

// Info returns an info client.
func (c *Client) Info() *InfoClient {
	return newInfoClient(c.transport)
}

// Status returns a status client.
func (c *Client) Status() *StatusClient {
	return newStatusClient(c.transport)
}

// Provisioning returns a provisioning client.
func (c *Client) Provisioning() *ProvisioningClient {
	return newProvisioningClient(c.transport)
}
