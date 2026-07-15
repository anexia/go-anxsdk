package v1

import (
	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/v1/core"
	"github.com/anexia/go-anxsdk/v1/ipam"
	"github.com/anexia/go-anxsdk/v1/lbaas"
	"github.com/anexia/go-anxsdk/v1/vsphere"
)

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

// Ipam is the entry point to ipam related clients.
func (c *Client) Ipam() *ipam.Client {
	return ipam.NewClient(c.transport)
}

// LBaaS is the entry point to lbaas related clients.
func (c *Client) LBaaS() *lbaas.Client {
	return lbaas.NewClient(c.transport)
}

// VSphere is the entry point for dynamic compute client.
func (c *Client) VSphere() *vsphere.Client {
	return vsphere.NewClient(c.transport)
}

// Core is the entry point for core engine api services.
func (c *Client) Core() *core.Client {
	return core.NewClient(c.transport)
}
