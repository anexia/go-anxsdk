package lbaas

import (
	"context"
	"fmt"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/v1/common"
)

// ACLListParams defines the available parameters for the acl list endpoint.
type ACLListParams struct {
	Search string `url:"search,omitempty"`
}

// ACLListItem is an item in the acl list response.
type ACLListItem struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

// GetID returns the Identifier of the ACLListItem.
func (i ACLListItem) GetID() string {
	return i.Identifier
}

// ACLGetResponse represents the response of the acl get endpoint.
type ACLGetResponse struct {
	CustomerIdentifier         *string              `json:"customer_identifier"`
	ResellerIdentifier         string               `json:"reseller_identifier"`
	CriticalOperationPassword  *string              `json:"critical_operation_password"`
	CriticalOperationConfirmed bool                 `json:"critical_operation_confirmed"`
	Identifier                 string               `json:"identifier"`
	Name                       string               `json:"name"`
	State                      common.State[string] `json:"state"`
	ParentType                 string               `json:"parent_type"`
	Frontend                   common.Resource      `json:"frontend"`
	Backend                    *common.Resource     `json:"backend"`
	Criterion                  string               `json:"criterion"`
	Index                      int                  `json:"index"`
	Value                      string               `json:"value"`
	Redeploy                   bool                 `json:"redeploy"`
	AutomationRules            []common.Resource    `json:"automation_rules"`
}

// ACLClient is an api client for managing load balancer acls.
type ACLClient struct {
	transport *internal.Transport
}

// NewACLClient creates a new acl client.
func NewACLClient(transport *internal.Transport) *ACLClient {
	return &ACLClient{
		transport: transport,
	}
}

// List returns a list of paged acls.
func (c *ACLClient) List(ctx context.Context, pagingParams paging.Params, params ACLListParams) (paging.PagedResponse[ACLListItem], error) {
	resp := internal.RequestWrapper[paging.PagedResponse[ACLListItem]]{}
	err := c.transport.Get(ctx, "/api/LBaaS/v1/acl.json", &resp, pagingParams, params)
	return resp.Data, common.MapTransportError(err)
}

// ListPageFetcher returns a paging.PageFetcher for acls.
func (c *ACLClient) ListPageFetcher(params ACLListParams) paging.PageFetcher[ACLListItem] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[ACLListItem], error) {
		return c.List(ctx, pageParams, params)
	}
}

// Get returns a acl by identifier.
func (c *ACLClient) Get(ctx context.Context, identifier string) (ACLGetResponse, error) {
	resp := ACLGetResponse{}
	err := c.transport.GetSingle(ctx, fmt.Sprintf("/api/LBaaS/v1/acl.json/%s", identifier), &resp)
	return resp, common.MapTransportError(err)
}
