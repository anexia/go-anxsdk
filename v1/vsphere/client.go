package vsphere

import (
	"github.com/anexia/go-anxsdk/internal"
)

//go:generate go tool counterfeiter -generate

//counterfeiter:generate -o ../../fake/v1/vsphere . API

// API is the interface implemented by [Client].
type API interface {
	// PowerControl returns a power control client.
	PowerControl() PowerControlAPI
	// Info returns an info client.
	Info() InfoAPI
	// Status returns a status client.
	Status() StatusAPI
	// Provisioning returns a provisioning client.
	Provisioning() ProvisioningAPI
	// Search returns a search client.
	Search() SearchAPI
}

// Client is an anexia vsphere v1 api client.
type Client struct {
	transport *internal.Transport
}

// Check that Client implements the API interface.
var _ API = (*Client)(nil)

// NewClient creates a new v1 vsphere api client.
func NewClient(transport *internal.Transport) *Client {
	return &Client{
		transport: transport,
	}
}

// PowerControl returns a power control client.
func (c *Client) PowerControl() PowerControlAPI {
	return newPowerControlClient(c.transport)
}

// Info returns an info client.
func (c *Client) Info() InfoAPI {
	return newInfoClient(c.transport)
}

// Status returns a status client.
func (c *Client) Status() StatusAPI {
	return newStatusClient(c.transport)
}

// Provisioning returns a provisioning client.
func (c *Client) Provisioning() ProvisioningAPI {
	return newProvisioningClient(c.transport)
}

// Search returns a search client.
func (c *Client) Search() SearchAPI {
	return newSearchClient(c.transport)
}
