package lbaas

import (
	"context"
	"fmt"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/v2/common"
)

// BackendSwitchingRuleListParams defines the available parameters for the backend switching rule list endpoint.
type BackendSwitchingRuleListParams struct {
	Search string `url:"search,omitempty"`
}

// BackendSwitchingRuleListItem is an item in the backend switching rule list response.
type BackendSwitchingRuleListItem struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

// BackendSwitchingRuleGetResponse represents the response of the backend switching rule get endpoint.
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
func (c *BackendSwitchingRuleClient) List(
	ctx context.Context,
	pagingParams paging.Params,
	params BackendSwitchingRuleListParams,
) (paging.PagedResponse[BackendSwitchingRuleListItem], error) {
	resp := paging.PagedResponse[BackendSwitchingRuleListItem]{}
	err := c.transport.Get(ctx, "/api/LBaaS/v2/backendswitchingrule", &resp, pagingParams, params)
	return resp, common.MapTransportError(err)
}

// ListPageFetcher returns a paging.PageFetcher for backendSwitchingRules.
func (c *BackendSwitchingRuleClient) ListPageFetcher(params BackendSwitchingRuleListParams) paging.PageFetcher[BackendSwitchingRuleListItem] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[BackendSwitchingRuleListItem], error) {
		return c.List(ctx, pageParams, params)
	}
}

// ListFull returns a list of paged backendSwitchingRules with all fields.
func (c *BackendSwitchingRuleClient) ListFull(
	ctx context.Context,
	pagingParams paging.Params,
	params BackendSwitchingRuleListParams,
) (paging.PagedResponse[BackendSwitchingRuleGetResponse], error) {
	resp := paging.PagedResponse[BackendSwitchingRuleGetResponse]{}
	wrap := internal.NewAllAttributesWrapper(params)
	err := c.transport.Get(ctx, "/api/LBaaS/v2/backendswitchingrule", &resp, pagingParams, wrap)
	return resp, common.MapTransportError(err)
}

// ListFullPageFetcher returns a paging.PageFetcher for backendSwitchingRules with all fields.
func (c *BackendSwitchingRuleClient) ListFullPageFetcher(params BackendSwitchingRuleListParams) paging.PageFetcher[BackendSwitchingRuleGetResponse] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[BackendSwitchingRuleGetResponse], error) {
		return c.ListFull(ctx, pageParams, params)
	}
}

// Get returns a backendSwitchingRule by identifier.
func (c *BackendSwitchingRuleClient) Get(ctx context.Context, identifier string) (BackendSwitchingRuleGetResponse, error) {
	resp := BackendSwitchingRuleGetResponse{}
	err := c.transport.GetSingle(ctx, fmt.Sprintf("/api/LBaaS/v2/backendswitchingrule/%s", identifier), &resp)
	return resp, common.MapTransportError(err)
}
