package lbaas

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/v2/common"
)

// FrontendFilters is a struct that represents all filterable fields of a frontend.
type FrontendFilters struct {
	LoadBalancer *string
}

// EncodeValues implements the query.Encode interface for FrontendFilters.
func (f *FrontendFilters) EncodeValues(key string, v *url.Values) error {
	sb := strings.Builder{}

	if f.LoadBalancer != nil {
		_, _ = fmt.Fprintf(&sb, "load_balancer=%s", *f.LoadBalancer)
	}

	v.Add(key, sb.String())
	return nil
}

// FrontendListParams defines the available parameters for the frontend list endpoint.
type FrontendListParams struct {
	Search  string           `url:"search,omitempty"`
	Filters *FrontendFilters `url:"filters,omitempty"`
}

// FrontendListItem is an item in the frontend list response.
type FrontendListItem struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

// FrontendGetResponse represents the response of the frontend get endpoint.
type FrontendGetResponse struct {
	CustomerIdentifier         *string              `json:"customer_identifier"`
	ResellerIdentifier         string               `json:"reseller_identifier"`
	CriticalOperationPassword  *string              `json:"critical_operation_password"`
	CriticalOperationConfirmed bool                 `json:"critical_operation_confirmed"`
	Identifier                 string               `json:"identifier"`
	Name                       string               `json:"name"`
	State                      common.State[string] `json:"state"`
	Enable                     bool                 `json:"enable"`
	LoadBalancer               common.Resource      `json:"load_balancer"`
	DefaultBackend             common.Resource      `json:"default_backend"`
	Mode                       string               `json:"mode"`
	ClientTimeout              string               `json:"client_timeout"`
	Redeploy                   bool                 `json:"redeploy"`
	AutomationRules            []common.Resource    `json:"automation_rules"`
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
	resp := paging.PagedResponse[FrontendListItem]{}
	err := c.transport.Get(ctx, "/api/LBaaS/v2/frontend", &resp, pagingParams, params)
	return resp, common.MapTransportError(err)
}

// ListPageFetcher returns a paging.PageFetcher for frontends.
func (c *FrontendClient) ListPageFetcher(params FrontendListParams) paging.PageFetcher[FrontendListItem] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[FrontendListItem], error) {
		return c.List(ctx, pageParams, params)
	}
}

// ListFull returns a list of paged frontends with all fields.
func (c *FrontendClient) ListFull(ctx context.Context, pagingParams paging.Params, params FrontendListParams) (paging.PagedResponse[FrontendGetResponse], error) {
	resp := paging.PagedResponse[FrontendGetResponse]{}
	wrap := internal.NewAllAttributesWrapper(params)
	err := c.transport.Get(ctx, "/api/LBaaS/v2/frontend", &resp, pagingParams, wrap)
	return resp, common.MapTransportError(err)
}

// ListFullPageFetcher returns a paging.PageFetcher for frontends iwth all fields.
func (c *FrontendClient) ListFullPageFetcher(params FrontendListParams) paging.PageFetcher[FrontendGetResponse] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[FrontendGetResponse], error) {
		return c.ListFull(ctx, pageParams, params)
	}
}

// Get returns a frontend by identifier.
func (c *FrontendClient) Get(ctx context.Context, identifier string) (FrontendGetResponse, error) {
	resp := FrontendGetResponse{}
	err := c.transport.GetSingle(ctx, fmt.Sprintf("/api/LBaaS/v2/frontend/%s", identifier), &resp)
	return resp, common.MapTransportError(err)
}
