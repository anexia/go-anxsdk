package lbaas

import "github.com/anexia/go-anxsdk/internal"

// A struct that hols access method to lbaas related resource clients.
type Client struct {
	transport *internal.Transport
}

// NewClient returns a new lbaas.Client.
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

// Acls returns an acl client.
func (c *Client) Acls() *ACLClient {
	return NewACLClient(c.transport)
}

// Maps returns a map client.
func (c *Client) Maps() *MapClient {
	return NewMapClient(c.transport)
}

// Rules returns a rule client.
func (c *Client) Rules() *RuleClient {
	return NewRuleClient(c.transport)
}

// Servers returns a server client.
func (c *Client) Servers() *ServerClient {
	return NewServerClient(c.transport)
}

// Ssls returns an ssl client.
func (c *Client) Ssls() *SslClient {
	return NewSslClient(c.transport)
}
