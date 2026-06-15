package v1

import (
	"context"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/paging"
)

type LoadBalancerListParams struct {
	Search string `url:"search,omitempty"`
}

type LoadBalancerListItem struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

// LoadBalancerClient is an api client for managing load balancers.
type LoadBalancerClient struct {
	transport *internal.Transport
}

// NewLoadBalancerClient creates a new load balancer client.
func NewLoadBalancerClient(transport *internal.Transport) *LoadBalancerClient {
	return &LoadBalancerClient{
		transport: transport,
	}
}

// List returns a list of paged load balancers.
func (c *LoadBalancerClient) List(ctx context.Context, pagingParams paging.Params, params LoadBalancerListParams) (paging.PagedResponse[LoadBalancerListItem], error) {
	resp := internal.RequestWrapper[paging.PagedResponse[LoadBalancerListItem]]{}
	err := c.transport.Get(ctx, "/api/LBaaS/v1/loadbalancer.json", &resp, pagingParams, params)
	return resp.Data, mapTransportError(err)
}

// ListPageFetcher returns a paging.PageFetcher for load balancers.
func (c *LoadBalancerClient) ListPageFetcher(params LoadBalancerListParams) paging.PageFetcher[LoadBalancerListItem] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[LoadBalancerListItem], error) {
		return c.List(ctx, pageParams, params)
	}
}
