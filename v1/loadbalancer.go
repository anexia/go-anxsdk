package v1

import (
	"context"
	"fmt"

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

// GetID returns the Identifier of the LoadBalancerListItem.
func (i LoadBalancerListItem) GetID() string {
	return i.Identifier
}

type LoadBalancerGetResponse struct {
	CustomerIdentifier         *string    `json:"customer_identifier"`
	ResellerIdentifier         string     `json:"reseller_identifier"`
	CriticalOperationPassword  *string    `json:"critical_operation_password"`
	CriticalOperationConfirmed bool       `json:"critical_operation_confirmed"`
	Identifier                 string     `json:"identifier"`
	Name                       string     `json:"name"`
	State                      State      `json:"state"`
	IpAddress                  string     `json:"ip_address"`
	AutomationRules            []Resource `json:"automation_rules"`
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

// Get returns a load balancer by identifier.
func (c *LoadBalancerClient) Get(ctx context.Context, identifier string) (LoadBalancerGetResponse, error) {
	resp := LoadBalancerGetResponse{}
	err := c.transport.GetSingle(ctx, fmt.Sprintf("/api/LBaaS/v1/loadbalancer.json/%s", identifier), &resp)
	return resp, mapTransportError(err)
}
