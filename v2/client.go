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

// ClustersByEnv returns a clusters client based on the provided env.
func (c *Client) ClustersByEnv(env KubernetesEnv) *ClustersClient {
	return newClustersClient(c.transport, env)
}

// Clusters returns a clusters client.
func (c *Client) Clusters() *ClustersClient {
	return newClustersClient(c.transport, KubernetesEnvProduction)
}

// StageClusters returns a stage clusters client.
func (c *Client) StageClusters() *ClustersClient {
	return newClustersClient(c.transport, KubernetesEnvStaging)
}

// DevClusters returns a dev clusters client.
func (c *Client) DevClusters() *ClustersClient {
	return newClustersClient(c.transport, KubernetesEnvDevelopment)
}
