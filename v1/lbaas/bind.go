package lbaas

import (
	"context"
	"fmt"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/v1/common"
)

type BindListParams struct {
	Search string `url:"search,omitempty"`
}

type BindListItem struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

// GetID returns the Itentifier of the BindListItem.
func (i BindListItem) GetID() string {
	return i.Identifier
}

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
	resp := internal.RequestWrapper[paging.PagedResponse[BindListItem]]{}
	err := c.transport.Get(ctx, "/api/LBaaS/v1/bind.json", &resp, pagingParams, params)
	return resp.Data, common.MapTransportError(err)
}

// ListPageFetcher returns a paging.PageFetcher for binds.
func (c *BindClient) ListPageFetcher(params BindListParams) paging.PageFetcher[BindListItem] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[BindListItem], error) {
		return c.List(ctx, pageParams, params)
	}
}

// Get returns a bind by identifier.
func (c *BindClient) Get(ctx context.Context, identifier string) (BindGetResponse, error) {
	resp := BindGetResponse{}
	err := c.transport.GetSingle(ctx, fmt.Sprintf("/api/LBaaS/v1/bind.json/%s", identifier), &resp)
	return resp, common.MapTransportError(err)
}
