package kubernetes

import "github.com/anexia/go-anxsdk/internal"

// API is the interface implemented by [Client].
type API interface {
	// Clusters returns a clusters client.
	Clusters() ClustersAPI
	// Nodepools returns a nodepools client.
	Nodepools() NodepoolsAPI
	// NodepoolDisks returns a nodepool disk client.
	NodepoolDisks() NodepoolDisksAPI
	// NodepoolNetworks returns a nodepool network client.
	NodepoolNetworks() NodepoolNetworksAPI
}

// Client is an anexia v2 api client.
type Client struct {
	transport   *internal.Transport
	environment Env
}

// Check that Client implements the API interface.
var _ API = (*Client)(nil)

// NewClient creates a new v2 api client.
func NewClient(transport *internal.Transport, environment Env) *Client {
	return &Client{
		transport:   transport,
		environment: environment,
	}
}

// Clusters returns a clusters client.
func (c *Client) Clusters() ClustersAPI {
	return newClustersClient(c.transport, c.environment)
}

// Nodepools returns a nodepools client.
func (c *Client) Nodepools() NodepoolsAPI {
	return newNodepoolsClient(c.transport, c.environment)
}

// NodepoolDisks returns a nodepool disk client.
func (c *Client) NodepoolDisks() NodepoolDisksAPI {
	return newNodepoolDisksClient(c.transport, c.environment)
}

// NodepoolNetworks returns a nodepool network client.
func (c *Client) NodepoolNetworks() NodepoolNetworksAPI {
	return newNodepoolNetworksClient(c.transport, c.environment)
}
