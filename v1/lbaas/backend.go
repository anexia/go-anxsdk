package lbaas

import (
	"context"
	"fmt"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/v1/common"
)

// BackendListParams defines the available parameters for the backend list endpoint.
type BackendListParams struct {
	Search string `url:"search,omitempty"`
}

// BackendListItem is an item in the backend list response.
type BackendListItem struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

// GetID returns the Itentifier of the BackendListItem.
func (i BackendListItem) GetID() string {
	return i.Identifier
}

// BackendGetResponse represents the response of the backend get endpoint.
type BackendGetResponse struct {
	CustomerIdentifier         *string              `json:"customer_identifier"`
	ResellerIdentifier         string               `json:"reseller_identifier"`
	CriticalOperationPassword  *string              `json:"critical_operation_password"`
	CriticalOperationConfirmed bool                 `json:"critical_operation_confirmed"`
	Identifier                 string               `json:"identifier"`
	Name                       string               `json:"name"`
	State                      common.State[string] `json:"state"`
	Enable                     bool                 `json:"enable"`
	LoadBalancer               common.Resource      `json:"load_balancer"`
	HealthCheck                string               `json:"health_check"`
	Mode                       string               `json:"mode"`
	ServerTimeout              int                  `json:"server_timeout"`
	Redeploy                   bool                 `json:"redeploy"`
	AutomationRules            []common.Resource    `json:"automation_rules"`
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
	resp := internal.RequestWrapper[paging.PagedResponse[BackendListItem]]{}
	err := c.transport.Get(ctx, "/api/LBaaS/v1/backend.json", &resp, pagingParams, params)
	return resp.Data, common.MapTransportError(err)
}

// ListPageFetcher returns a paging.PageFetcher for backends.
func (c *BackendClient) ListPageFetcher(params BackendListParams) paging.PageFetcher[BackendListItem] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[BackendListItem], error) {
		return c.List(ctx, pageParams, params)
	}
}

// Get returns a backend by identifier.
func (c *BackendClient) Get(ctx context.Context, identifier string) (BackendGetResponse, error) {
	resp := BackendGetResponse{}
	err := c.transport.GetSingle(ctx, fmt.Sprintf("/api/LBaaS/v1/backend.json/%s", identifier), &resp)
	return resp, common.MapTransportError(err)
}
