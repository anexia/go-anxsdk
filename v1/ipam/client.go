package ipam

import "github.com/anexia/go-anxsdk/internal"

// Client is an anexia ipam v1 api client.
type Client struct {
	transport *internal.Transport
}

// NewClient creates a new v1 ipam api client.
func NewClient(transport *internal.Transport) *Client {
	return &Client{
		transport: transport,
	}
}

// Addresses returns an address client.
func (c *Client) Addresses() *AddressClient {
	return NewAddressClient(c.transport)
}

// Aggregates returns an aggregate client.
func (c *Client) Aggregates() *AggregateClient {
	return NewAggregateClient(c.transport)
}

// Prefixes returns a prefix client.
func (c *Client) Prefixes() *PrefixClient {
	return NewPrefixClient(c.transport)
}

// Zones returns a zone client.
func (c *Client) Zones() *ZoneClient {
	return NewZoneClient(c.transport)
}
