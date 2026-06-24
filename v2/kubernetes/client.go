package kubernetes

import "github.com/anexia/go-anxsdk/internal"

// Client is an anexia v2 api client.
type Client struct {
	transport   *internal.Transport
	environment Env
}

// NewClient creates a new v2 api client.
func NewClient(transport *internal.Transport, environment Env) *Client {
	return &Client{
		transport:   transport,
		environment: environment,
	}
}

// Clusters returns a clusters client.
func (c *Client) Clusters() *ClustersClient {
	return newClustersClient(c.transport, c.environment)
}

// Nodepools returns a nodepools client.
func (c *Client) Nodepools() *NodepoolsClient {
	return newNodepoolsClient(c.transport, c.environment)
}

// NodepoolDisks returns a nodepool disk client.
func (c *Client) NodepoolDisks() *NodepoolDisksClient {
	return newNodepoolDisksClient(c.transport, c.environment)
}

// NodepoolNetworks returns a nodepool network client.
func (c *Client) NodepoolNetworks() *NodepoolNetworksClient {
	return newNodepoolNetworksClient(c.transport, c.environment)
}
