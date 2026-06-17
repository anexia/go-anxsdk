package lbaas

import (
	"context"
	"fmt"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/paging"
)

type BackendSwitchingRuleListParams struct {
	Search string `url:"search,omitempty"`
}

type BackendSwitchingRuleListItem struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

// GetID returns the Itentifier of the BackendSwitchingRuleListItem.
func (i BackendSwitchingRuleListItem) GetID() string {
	return i.Identifier
}

type BackendSwitchingRuleGetResponse struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

// BackendSwitchingRuleClient is an api client for managing load balancer backendSwitchingRules.
type BackendSwitchingRuleClient struct {
	transport *internal.Transport
}

// NewBackendSwitchingRuleClient creates a new backendSwitchingRule client.
func NewBackendSwitchingRuleClient(transport *internal.Transport) *BackendSwitchingRuleClient {
	return &BackendSwitchingRuleClient{
		transport: transport,
	}
}

// List returns a list of paged backendSwitchingRules.
func (c *BackendSwitchingRuleClient) List(ctx context.Context, pagingParams paging.Params, params BackendSwitchingRuleListParams) (paging.PagedResponse[BackendSwitchingRuleListItem], error) {
	resp := internal.RequestWrapper[paging.PagedResponse[BackendSwitchingRuleListItem]]{}
	err := c.transport.Get(ctx, "/api/LBaaS/v1/backendswitchingrule.json", &resp, pagingParams, params)
	return resp.Data, v1.mapTransportError(err)
}

// ListPageFetcher returns a paging.PageFetcher for backendSwitchingRules.
func (c *BackendSwitchingRuleClient) ListPageFetcher(params BackendSwitchingRuleListParams) paging.PageFetcher[BackendSwitchingRuleListItem] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[BackendSwitchingRuleListItem], error) {
		return c.List(ctx, pageParams, params)
	}
}

// Get returns a backendSwitchingRule by identifier.
func (c *BackendSwitchingRuleClient) Get(ctx context.Context, identifier string) (BackendSwitchingRuleGetResponse, error) {
	resp := BackendSwitchingRuleGetResponse{}
	err := c.transport.GetSingle(ctx, fmt.Sprintf("/api/LBaaS/v1/backendswitchingrule.json/%s", identifier), &resp)
	return resp, v1.mapTransportError(err)
}
