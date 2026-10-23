package ipam

import "github.com/anexia/go-anxsdk/internal"

//go:generate go tool counterfeiter -generate

//counterfeiter:generate -o ../../fake/v1/ipam . API

// API is the interface implemented by [Client].
type API interface {
	// Addresses returns an address client.
	Addresses() AddressAPI
	// Aggregates returns an aggregate client.
	Aggregates() AggregateAPI
	// Prefixes returns a prefix client.
	Prefixes() PrefixAPI
	// Zones returns a zone client.
	Zones() ZoneAPI
}

// Client is an anexia ipam v1 api client.
type Client struct {
	transport *internal.Transport
}

// Check that Client implements the API interface.
var _ API = (*Client)(nil)

// NewClient creates a new v1 ipam api client.
func NewClient(transport *internal.Transport) *Client {
	return &Client{
		transport: transport,
	}
}

// Addresses returns an address client.
func (c *Client) Addresses() AddressAPI {
	return NewAddressClient(c.transport)
}

// Aggregates returns an aggregate client.
func (c *Client) Aggregates() AggregateAPI {
	return NewAggregateClient(c.transport)
}

// Prefixes returns a prefix client.
func (c *Client) Prefixes() PrefixAPI {
	return NewPrefixClient(c.transport)
}

// Zones returns a zone client.
func (c *Client) Zones() ZoneAPI {
	return NewZoneClient(c.transport)
}
