package lbaas

import "github.com/anexia/go-anxsdk/internal"

// API is the interface implemented by [Client].
type API interface {
	// Backends returns a backend client.
	Backends() BackendAPI
	// BackendSwitchingRules returns a backend switching rules client.
	BackendSwitchingRules() BackendSwitchingRuleAPI
	// Binds returns a binds client.
	Binds() BindAPI
	// Frontends returns a frontend client.
	Frontends() FrontendAPI
	// LoadBalancers returns a load balancers client.
	LoadBalancers() LoadBalancerAPI
	// Acls returns an acl client.
	Acls() ACLAPI
	// Maps returns a map client.
	Maps() MapAPI
	// Rules returns a rule client.
	Rules() RuleAPI
	// Servers returns a server client.
	Servers() ServerAPI
	// Ssls returns an ssl client.
	Ssls() SslAPI
}

// Client is a struct that holds access method to lbaas related resource clients.
type Client struct {
	transport *internal.Transport
}

// Check that Client implements the API interface.
var _ API = (*Client)(nil)

// NewClient returns a new lbaas.Client.
func NewClient(transport *internal.Transport) *Client {
	return &Client{
		transport: transport,
	}
}

// Backends returns a backend client.
func (c *Client) Backends() BackendAPI {
	return NewBackendClient(c.transport)
}

// BackendSwitchingRules returns a backend switching rules client.
func (c *Client) BackendSwitchingRules() BackendSwitchingRuleAPI {
	return NewBackendSwitchingRuleClient(c.transport)
}

// Binds returns a binds client.
func (c *Client) Binds() BindAPI {
	return NewBindClient(c.transport)
}

// Frontends returns a frontend client.
func (c *Client) Frontends() FrontendAPI {
	return NewFrontendClient(c.transport)
}

// LoadBalancers returns a load balancers client.
func (c *Client) LoadBalancers() LoadBalancerAPI {
	return NewLoadBalancerClient(c.transport)
}

// Acls returns an acl client.
func (c *Client) Acls() ACLAPI {
	return NewACLClient(c.transport)
}

// Maps returns a map client.
func (c *Client) Maps() MapAPI {
	return NewMapClient(c.transport)
}

// Rules returns a rule client.
func (c *Client) Rules() RuleAPI {
	return NewRuleClient(c.transport)
}

// Servers returns a server client.
func (c *Client) Servers() ServerAPI {
	return NewServerClient(c.transport)
}

// Ssls returns an ssl client.
func (c *Client) Ssls() SslAPI {
	return NewSslClient(c.transport)
}
