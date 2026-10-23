package v1

import (
	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/v1/core"
	"github.com/anexia/go-anxsdk/v1/ipam"
	"github.com/anexia/go-anxsdk/v1/lbaas"
	"github.com/anexia/go-anxsdk/v1/vsphere"
)

// API is the interface implemented by [Client].
type API interface {
	// Vlans returns a vlans client.
	Vlans() VlansAPI
	// Clusters returns a clusters client.
	Clusters() ClustersAPI
	// StageClusters returns a stage clusters client.
	StageClusters() ClustersAPI
	// DevClusters returns a dev clusters client.
	DevClusters() ClustersAPI
	// Ipam is the entry point to ipam related clients.
	Ipam() ipam.API
	// LBaaS is the entry point to lbaas related clients.
	LBaaS() lbaas.API
	// VSphere is the entry point for dynamic compute client.
	VSphere() vsphere.API
	// Core is the entry point for core engine api services.
	Core() core.API
}

// Client is an anexia v1 api client.
type Client struct {
	transport *internal.Transport
}

// Check that Client implements the API interface.
var _ API = (*Client)(nil)

// NewClient creates a new v1 api client.
func NewClient(transport *internal.Transport) *Client {
	return &Client{
		transport: transport,
	}
}

// Vlans returns a vlans client.
func (c *Client) Vlans() VlansAPI {
	return NewVlansClient(c.transport)
}

// Clusters returns a clusters client.
func (c *Client) Clusters() ClustersAPI {
	return newClustersClient(c.transport, kubernetesEnvironmentProduction)
}

// StageClusters returns a stage clusters client.
func (c *Client) StageClusters() ClustersAPI {
	return newClustersClient(c.transport, kubernetesEnvironmentStaging)
}

// DevClusters returns a dev clusters client.
func (c *Client) DevClusters() ClustersAPI {
	return newClustersClient(c.transport, kubernetesEnvironmentDevelopment)
}

// Ipam is the entry point to ipam related clients.
func (c *Client) Ipam() ipam.API {
	return ipam.NewClient(c.transport)
}

// LBaaS is the entry point to lbaas related clients.
func (c *Client) LBaaS() lbaas.API {
	return lbaas.NewClient(c.transport)
}

// VSphere is the entry point for dynamic compute client.
func (c *Client) VSphere() vsphere.API {
	return vsphere.NewClient(c.transport)
}

// Core is the entry point for core engine api services.
func (c *Client) Core() core.API {
	return core.NewClient(c.transport)
}
