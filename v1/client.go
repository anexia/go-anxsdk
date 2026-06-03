package v1

import "github.com/anexia/go-anxsdk/internal"

// Client is an anexia v1 api client.
type Client struct {
	transport *internal.Transport
}

// NewClient creates a new v1 api client.
func NewClient(transport *internal.Transport) *Client {
	return &Client{
		transport: transport,
	}
}

// Vlans returns a vlans client.
func (c *Client) Vlans() *VlansClient {
	return NewVlansClient(c.transport)
}

// Locations returns a locations client.
func (c *Client) Locations() *LocationsClient {
	return NewLocationsClient(c.transport)
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
