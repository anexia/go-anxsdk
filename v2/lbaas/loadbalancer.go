package lbaas

import (
	"context"
	"fmt"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/v2/common"
)

// LoadBalancerListParams defines the available parameters for the load balancer list endpoint.
type LoadBalancerListParams struct {
	Search string `url:"search,omitempty"`
}

// LoadBalancerListItem is an item in the load balancer list response.
type LoadBalancerListItem struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

// LoadBalancerGetResponse represents the response of the load balancer get endpoint.
type LoadBalancerGetResponse struct {
	CustomerIdentifier         *string              `json:"customer_identifier"`
	ResellerIdentifier         string               `json:"reseller_identifier"`
	CriticalOperationPassword  *string              `json:"critical_operation_password"`
	CriticalOperationConfirmed bool                 `json:"critical_operation_confirmed"`
	Identifier                 string               `json:"identifier"`
	Name                       string               `json:"name"`
	State                      common.State[string] `json:"state"`
	IPAddress                  string               `json:"ip_address"`
	AutomationRules            []common.Resource    `json:"automation_rules"`
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

// Get returns a single load balancer by its id.
func (c *LoadBalancerClient) Get(ctx context.Context, identifier string) (LoadBalancerGetResponse, error) {
	resp := LoadBalancerGetResponse{}
	err := c.transport.GetSingle(ctx, fmt.Sprintf("/api/LBaaS/v2/loadbalancer/%s", identifier), &resp)
	return resp, common.MapTransportError(err)
}

// List returns a list of paged load balancers.
func (c *LoadBalancerClient) List(ctx context.Context, pagingParams paging.Params, params LoadBalancerListParams) (paging.PagedResponse[LoadBalancerListItem], error) {
	resp := paging.PagedResponse[LoadBalancerListItem]{}
	err := c.transport.Get(ctx, "/api/LBaaS/v2/loadbalancer", &resp, pagingParams, params)
	return resp, common.MapTransportError(err)
}

// ListPageFetcher returns a paging.PageFetcher for clusters.
func (c *LoadBalancerClient) ListPageFetcher(params LoadBalancerListParams) paging.PageFetcher[LoadBalancerListItem] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[LoadBalancerListItem], error) {
		return c.List(ctx, pageParams, params)
	}
}

// ListFull returns a list of paged clusters with all attributes included.
func (c *LoadBalancerClient) ListFull(ctx context.Context, pageParams paging.Params, params LoadBalancerListParams) (paging.PagedResponse[LoadBalancerGetResponse], error) {
	resp := paging.PagedResponse[LoadBalancerGetResponse]{}
	wrap := internal.NewAllAttributesWrapper(params)
	err := c.transport.Get(ctx, "/api/LBaaS/v2/loadbalancer", &resp, pageParams, wrap)
	return resp, common.MapTransportError(err)
}

// ListFullPageFetcher returns a paging.PageFetcher for clusters with all attributes included.
func (c *LoadBalancerClient) ListFullPageFetcher(params LoadBalancerListParams) paging.PageFetcher[LoadBalancerGetResponse] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[LoadBalancerGetResponse], error) {
		return c.ListFull(ctx, pageParams, params)
	}
}
