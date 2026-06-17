package lbaas

import (
	"github.com/anexia/go-anxsdk/internal"
)

// Client is an anexia lbaas v1 api client.
type Client struct {
	transport *internal.Transport
}

// NewClient creates a new v1 lbaas api client.
func NewClient(transport *internal.Transport) *Client {
	return &Client{
		transport: transport,
	}
}

// Backends returns a backend client.
func (c *Client) Backends() *BackendClient {
	return NewBackendClient(c.transport)
}

// BackendSwitchingRules returns a backend switching rules client.
func (c *Client) BackendSwitchingRules() *BackendSwitchingRuleClient {
	return NewBackendSwitchingRuleClient(c.transport)
}

// Binds returns a binds client.
func (c *Client) Binds() *BindClient {
	return NewBindClient(c.transport)
}

// Frontends returns a frontend client.
func (c *Client) Frontends() *FrontendClient {
	return NewFrontendClient(c.transport)
}

// LoadBalancers returns a load balancers client.
func (c *Client) LoadBalancers() *LoadBalancerClient {
	return NewLoadBalancerClient(c.transport)
}
