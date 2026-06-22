package lbaas

import (
	"context"
	"fmt"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/v2/common"
)

// BackendListParams defines the available parameters for the backend list endpoint.
type BackendListParams struct {
	Search       string `url:"search,omitempty"`
	LoadBalancer string `url:"load_balancer,omitempty"`
}

// BackendListItem is an item in the backend list response.
type BackendListItem struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

// BackendGetResponse represents the response of the backend get endpoint.
type BackendGetResponse struct {
	CustomerIdentifier         *string                     `json:"customer_identifier"`
	ResellerIdentifier         string                      `json:"reseller_identifier"`
	CriticalOperationPassword  *string                     `json:"critical_operation_password"`
	CriticalOperationConfirmed bool                        `json:"critical_operation_confirmed"`
	Identifier                 string                      `json:"identifier"`
	Name                       string                      `json:"name"`
	State                      common.State[string]        `json:"state"`
	Enable                     bool                        `json:"enable"`
	LoadBalancer               common.Resource             `json:"load_balancer"`
	HealthCheck                string                      `json:"health_check"`
	Mode                       common.IDTitleTuple[string] `json:"mode"`
	ServerTimeout              int                         `json:"server_timeout"`
	Redeploy                   bool                        `json:"redeploy"`
	AutomationRules            []common.Resource           `json:"automation_rules"`
}

// BackendClient is an api client for managing load balancer backends.
type BackendClient struct {
	transport *internal.Transport
}

// NewBackendClient creates a new backend client.
func NewBackendClient(transport *internal.Transport) *BackendClient {
	return &BackendClient{
		transport: transport,
	}
}

// List returns a list of paged backends.
func (c *BackendClient) List(ctx context.Context, pagingParams paging.Params, params BackendListParams) (paging.PagedResponse[BackendListItem], error) {
	resp := paging.PagedResponse[BackendListItem]{}
	err := c.transport.Get(ctx, "/api/LBaaS/v2/backend", &resp, pagingParams, params)
	return resp, common.MapTransportError(err)
}

// ListPageFetcher returns a paging.PageFetcher for backends.
func (c *BackendClient) ListPageFetcher(params BackendListParams) paging.PageFetcher[BackendListItem] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[BackendListItem], error) {
		return c.List(ctx, pageParams, params)
	}
}

// ListFull returns a list of paged backends with all fields.
func (c *BackendClient) ListFull(ctx context.Context, pagingParams paging.Params, params BackendListParams) (paging.PagedResponse[BackendGetResponse], error) {
	resp := paging.PagedResponse[BackendGetResponse]{}
	wrap := internal.NewAllAttributesWrapper(params)
	err := c.transport.Get(ctx, "/api/LBaaS/v2/backend", &resp, pagingParams, wrap)
	return resp, common.MapTransportError(err)
}

// ListFullPageFetcher returns a paging.PageFetcher for backends with all fields.
func (c *BackendClient) ListFullPageFetcher(params BackendListParams) paging.PageFetcher[BackendGetResponse] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[BackendGetResponse], error) {
		return c.ListFull(ctx, pageParams, params)
	}
}

// Get returns a backend by identifier.
func (c *BackendClient) Get(ctx context.Context, identifier string) (BackendGetResponse, error) {
	resp := BackendGetResponse{}
	err := c.transport.GetSingle(ctx, fmt.Sprintf("/api/LBaaS/v2/backend/%s", identifier), &resp)
	return resp, common.MapTransportError(err)
}
