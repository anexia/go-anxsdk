package lbaas

import (
	"context"
	"fmt"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/v2/common"
)

// BindListParams defines the available parameters for the bind list endpoint.
type BindListParams struct {
	Search string `url:"search,omitempty"`
}

// BindListItem is an item in the bind list response.
type BindListItem struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

// BindGetResponse represents the response of the bind get endpoint.
type BindGetResponse struct {
	CustomerIdentifier         *string              `json:"customer_identifier"`
	ResellerIdentifier         string               `json:"reseller_identifier"`
	CriticalOperationPassword  *string              `json:"critical_operation_password"`
	CriticalOperationConfirmed bool                 `json:"critical_operation_confirmed"`
	Identifier                 string               `json:"identifier"`
	Name                       string               `json:"name"`
	State                      common.State[string] `json:"state"`
	Frontend                   common.Resource      `json:"frontend"`
	Address                    *string              `json:"address"`
	Port                       int                  `json:"port"`
	Ssl                        bool                 `json:"ssl"`
	SslCertificatePath         string               `json:"ssl_certificate_path"`
	Redeploy                   bool                 `json:"redeploy"`
	AutomationRules            []common.Resource    `json:"automation_rules"`
}

// BindClient is an api client for managing load balancer binds.
type BindClient struct {
	transport *internal.Transport
}

// NewBindClient creates a new bind client.
func NewBindClient(transport *internal.Transport) *BindClient {
	return &BindClient{
		transport: transport,
	}
}

// List returns a list of paged binds.
func (c *BindClient) List(ctx context.Context, pagingParams paging.Params, params BindListParams) (paging.PagedResponse[BindListItem], error) {
	resp := paging.PagedResponse[BindListItem]{}
	err := c.transport.Get(ctx, "/api/LBaaS/v2/bind", &resp, pagingParams, params)
	return resp, common.MapTransportError(err)
}

// ListPageFetcher returns a paging.PageFetcher for binds.
func (c *BindClient) ListPageFetcher(params BindListParams) paging.PageFetcher[BindListItem] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[BindListItem], error) {
		return c.List(ctx, pageParams, params)
	}
}

// ListFull returns a list of paged binds with all fields.
func (c *BindClient) ListFull(ctx context.Context, pagingParams paging.Params, params BindListParams) (paging.PagedResponse[BindGetResponse], error) {
	resp := paging.PagedResponse[BindGetResponse]{}
	wrap := internal.NewAllAttributesWrapper(params)
	err := c.transport.Get(ctx, "/api/LBaaS/v2/bind", &resp, pagingParams, wrap)
	return resp, common.MapTransportError(err)
}

// ListFullPageFetcher returns a paging.PageFetcher for binds with all fields.
func (c *BindClient) ListFullPageFetcher(params BindListParams) paging.PageFetcher[BindGetResponse] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[BindGetResponse], error) {
		return c.ListFull(ctx, pageParams, params)
	}
}

// Get returns a bind by identifier.
func (c *BindClient) Get(ctx context.Context, identifier string) (BindGetResponse, error) {
	resp := BindGetResponse{}
	err := c.transport.GetSingle(ctx, fmt.Sprintf("/api/LBaaS/v2/bind/%s", identifier), &resp)
	return resp, common.MapTransportError(err)
}
