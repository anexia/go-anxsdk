package kubernetes

import "github.com/anexia/go-anxsdk/internal"

// Client is an anexia v2 api client.
type Client struct {
	transport   *internal.Transport
	environment KubernetesEnv
}

// NewClient creates a new v2 api client.
func NewClient(transport *internal.Transport, environment KubernetesEnv) *Client {
	return &Client{
		transport:   transport,
		environment: environment,
	}
}

// Clusters returns a clusters client.
func (c *Client) Clusters() *ClustersClient {
	return NewClustersClient(c.transport, c.environment)
}

// Nodepools returns a nodepools client.
func (c *Client) Nodepools() *NodepoolsClient {
	return NewNodepoolsClient(c.transport, c.environment)
}
