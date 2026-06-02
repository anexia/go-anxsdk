package v2

import "code.anexia.com/se/ks/go-anxsdk/internal"

// Client is an anexia v2 api client.
type Client struct {
	transport *internal.Transport
}

// NewClient creates a new v2 api client.
func NewClient(transport *internal.Transport) *Client {
	return &Client{
		transport: transport,
	}
}

// Clusters returns a clusters client.
func (c *Client) Clusters() *ClustersClient {
	return newClustersClient(c.transport, kubernetesEnvironmentProduction)
}

// StageClusters returns a stage clusters client.
func (c *Client) StageClusters() *ClustersClient {
	return newClustersClient(c.transport, kubernetesEnvironmentStaging)
}

// DevClusters returns a dev clusters client.
func (c *Client) DevClusters() *ClustersClient {
	return newClustersClient(c.transport, kubernetesEnvironmentDevelopment)
}
