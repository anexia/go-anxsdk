package v1

import (
	"context"
	"fmt"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/paging"
)

type FrontendListParams struct {
	Search string `url:"search,omitempty"`
}

type FrontendListItem struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

// GetID returns the Itentifier of the FrontendListItem.
func (i FrontendListItem) GetID() string {
	return i.Identifier
}

type FrontendGetResponse struct {
	CustomerIdentifier         *string    `json:"customer_identifier"`
	ResellerIdentifier         string     `json:"reseller_identifier"`
	CriticalOperationPassword  *string    `json:"critical_operation_password"`
	CriticalOperationConfirmed bool       `json:"critical_operation_confirmed"`
	Identifier                 string     `json:"identifier"`
	Name                       string     `json:"name"`
	State                      State      `json:"state"`
	Enable                     bool       `json:"enable"`
	LoadBalancer               Resource   `json:"load_balancer"`
	DefaultBackend             Resource   `json:"default_backend"`
	Mode                       string     `json:"mode"`
	ClientTimeout              string     `json:"client_timeout"`
	Redeploy                   bool       `json:"redeploy"`
	AutomationRules            []Resource `json:"automation_rules"`
}

// FrontendClient is an api client for managing load balancer frontends.
type FrontendClient struct {
	transport *internal.Transport
}

// NewFrontendClient creates a new frontend client.
func NewFrontendClient(transport *internal.Transport) *FrontendClient {
	return &FrontendClient{
		transport: transport,
	}
}

// List returns a list of paged frontends.
func (c *FrontendClient) List(ctx context.Context, pagingParams paging.Params, params FrontendListParams) (paging.PagedResponse[FrontendListItem], error) {
	resp := internal.RequestWrapper[paging.PagedResponse[FrontendListItem]]{}
	err := c.transport.Get(ctx, "/api/LBaaS/v1/frontend.json", &resp, pagingParams, params)
	return resp.Data, mapTransportError(err)
}

// ListPageFetcher returns a paging.PageFetcher for frontends.
func (c *FrontendClient) ListPageFetcher(params FrontendListParams) paging.PageFetcher[FrontendListItem] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[FrontendListItem], error) {
		return c.List(ctx, pageParams, params)
	}
}

// Get returns a frontend by identifier.
func (c *FrontendClient) Get(ctx context.Context, identifier string) (FrontendGetResponse, error) {
	resp := FrontendGetResponse{}
	err := c.transport.GetSingle(ctx, fmt.Sprintf("/api/LBaaS/v1/frontend.json/%s", identifier), &resp)
	return resp, mapTransportError(err)
}
