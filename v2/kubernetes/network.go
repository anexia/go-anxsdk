package kubernetes

import (
	"context"
	"fmt"
	"time"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/v2/common"
)

//revive:disable // self-explanatory constants

type NetworkBandwidthLimit string

const (
	NetworkBandwidthLimit100Mbit = NetworkBandwidthLimit("100")
	NetworkBandwidthLimit1Gbit   = NetworkBandwidthLimit("1000")
	NetworkBandwidthLimit10Gbit  = NetworkBandwidthLimit("10000")
)

//revive:enable

// NodepoolNetworkGetResponse represents the networks of a Nodepool.
type NodepoolNetworkGetResponse struct {
	Identifier string          `json:"identifier,omitempty"`
	Name       string          `json:"name,omitempty"`
	Nodepool   common.Resource `json:"nodepool,omitempty"`

	BandwidthLimit common.IDTitleTuple[NetworkBandwidthLimit] `json:"bandwidth_limit,omitempty"`
	VLAN           common.Resource                            `json:"vlan,omitempty"`

	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

// NodepoolNetworkUpdateRequest represents the networks of a Nodepool.
type NodepoolNetworkUpdateRequest struct {
	Name       string `json:"name,omitempty"`
	NodepoolID string `json:"nodepool,omitempty"`

	BandwidthLimit NetworkBandwidthLimit `json:"bandwidth_limit,omitempty"`
	VLANID         string                `json:"vlan,omitempty"`
}

// NodepoolNetworkListParams defines the available parameters for the network list endpoint.
type NodepoolNetworkListParams struct {
	Search string `url:"search,omitempty"`

	NodepoolID     string                `url:"nodepool,omitempty"`
	BandwidthLimit NetworkBandwidthLimit `url:"bandwidth_limit,omitempty"`
	VLANID         string                `url:"performance_type,omitempty"`
}

// NodepoolNetworkListItem is an item in the network list response.
type NodepoolNetworkListItem common.Resource

// NodepoolNetworksClient is an api client for managing NodepoolNetworks.
type NodepoolNetworksClient struct {
	environment Env
	transport   *internal.Transport
}

func newNodepoolNetworksClient(transport *internal.Transport, environment Env) *NodepoolNetworksClient {
	return &NodepoolNetworksClient{
		transport:   transport,
		environment: environment,
	}
}

func (c *NodepoolNetworksClient) endpointRoot() string {
	switch c.environment {
	case EnvDevelopment:
		return "/api/kubernetes-dev"
	case EnvStaging:
		return "/api/kubernetes-stage"
	case EnvProduction:
		return "/api/kubernetes"
	}

	return "/api/kubernetes"
}

// Get returns a single network by its id.
func (c *NodepoolNetworksClient) Get(ctx context.Context, identifier string) (NodepoolNetworkGetResponse, error) {
	resp := NodepoolNetworkGetResponse{}
	err := c.transport.GetSingle(ctx, fmt.Sprintf("%s/v2/node_pool_network/%s", c.endpointRoot(), identifier), &resp)
	return resp, common.MapTransportError(err)
}

// Update updates a network.
func (c *NodepoolNetworksClient) Update(ctx context.Context, identifier string, request NodepoolNetworkUpdateRequest) (NodepoolNetworkGetResponse, error) {
	resp := NodepoolNetworkGetResponse{}
	err := c.transport.Put(ctx, fmt.Sprintf("%s/v2/node_pool_network/%s", c.endpointRoot(), identifier), request, &resp)
	return resp, common.MapTransportError(err)
}

// List returns a list of paged Networks.
func (c *NodepoolNetworksClient) List(ctx context.Context, pageParams paging.Params, params NodepoolNetworkListParams) (paging.PagedResponse[NodepoolNetworkListItem], error) {
	resp := paging.PagedResponse[NodepoolNetworkListItem]{}
	err := c.transport.Get(ctx, fmt.Sprintf("%s/v2/node_pool_network", c.endpointRoot()), &resp, pageParams, params)
	return resp, common.MapTransportError(err)
}

// ListPageFetcher returns a paging.PageFetcher for Networks.
func (c *NodepoolNetworksClient) ListPageFetcher(params NodepoolNetworkListParams) paging.PageFetcher[NodepoolNetworkListItem] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[NodepoolNetworkListItem], error) {
		return c.List(ctx, pageParams, params)
	}
}

// ListFull returns a list of paged Networks with all attributes included.
func (c *NodepoolNetworksClient) ListFull(ctx context.Context, pageParams paging.Params, params NodepoolNetworkListParams,
) (paging.PagedResponse[NodepoolNetworkGetResponse], error) {
	resp := paging.PagedResponse[NodepoolNetworkGetResponse]{}
	wrap := internal.NewAllAttributesWrapper(params)
	err := c.transport.Get(ctx, fmt.Sprintf("%s/v2/node_pool_network", c.endpointRoot()), &resp, pageParams, wrap)
	return resp, common.MapTransportError(err)
}

// ListFullPageFetcher returns a paging.PageFetcher for Networks with all attributes included.
func (c *NodepoolNetworksClient) ListFullPageFetcher(params NodepoolNetworkListParams) paging.PageFetcher[NodepoolNetworkGetResponse] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[NodepoolNetworkGetResponse], error) {
		return c.ListFull(ctx, pageParams, params)
	}
}
