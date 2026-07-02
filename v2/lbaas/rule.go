package lbaas

import (
	"context"
	"fmt"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/v2/common"
)

// RuleListParams defines the available parameters for the rule list endpoint.
type RuleListParams struct {
	Search   string `url:"search,omitempty"`
	Frontend string `url:"frontend,omitempty"`
	Backend  string `url:"backend,omitempty"`
}

// RuleListItem is an item in the rule list response.
type RuleListItem struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

// RuleGetResponse represents the response of the rule get endpoint.
type RuleGetResponse struct {
	CustomerIdentifier         *string                     `json:"customer_identifier"`
	ResellerIdentifier         string                      `json:"reseller_identifier"`
	CriticalOperationPassword  *string                     `json:"critical_operation_password"`
	CriticalOperationConfirmed bool                        `json:"critical_operation_confirmed"`
	Identifier                 string                      `json:"identifier"`
	Name                       string                      `json:"name"`
	State                      common.State[string]        `json:"state"`
	RuleType                   common.IDTitleTuple[string] `json:"rule_type"`
	ParentType                 common.IDTitleTuple[string] `json:"parent_type"`
	Frontend                   common.Resource             `json:"frontend"`
	Backend                    *common.Resource            `json:"backend"`
	Index                      int                         `json:"index"`
	Condition                  common.IDTitleTuple[string] `json:"condition"`
	ConditionTest              string                      `json:"condition_test"`
	Type                       common.IDTitleTuple[string] `json:"type"`
	Action                     common.IDTitleTuple[string] `json:"action"`
	Redeploy                   bool                        `json:"redeploy"`
	RedirectionType            common.IDTitleTuple[string] `json:"redirection_type"`
	RedirectionValue           string                      `json:"redirection_value"`
	RedirectionCode            *string                     `json:"redirection_code"`
	AutomationRules            []common.Resource           `json:"automation_rules"`
}

// RuleClient is an api client for managing load balancer rules.
type RuleClient struct {
	transport *internal.Transport
}

// NewRuleClient creates a new rule client.
func NewRuleClient(transport *internal.Transport) *RuleClient {
	return &RuleClient{
		transport: transport,
	}
}

// List returns a list of paged rules.
func (c *RuleClient) List(ctx context.Context, pagingParams paging.Params, params RuleListParams) (paging.PagedResponse[RuleListItem], error) {
	resp := paging.PagedResponse[RuleListItem]{}
	err := c.transport.Get(ctx, "/api/LBaaS/v2/rule", &resp, pagingParams, params)
	return resp, common.MapTransportError(err)
}

// ListPageFetcher returns a paging.PageFetcher for rules.
func (c *RuleClient) ListPageFetcher(params RuleListParams) paging.PageFetcher[RuleListItem] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[RuleListItem], error) {
		return c.List(ctx, pageParams, params)
	}
}

// ListFull returns a list of paged rules with all fields.
func (c *RuleClient) ListFull(ctx context.Context, pagingParams paging.Params, params RuleListParams) (paging.PagedResponse[RuleGetResponse], error) {
	resp := paging.PagedResponse[RuleGetResponse]{}
	wrap := internal.NewAllAttributesWrapper(params)
	err := c.transport.Get(ctx, "/api/LBaaS/v2/rule", &resp, pagingParams, wrap)
	return resp, common.MapTransportError(err)
}

// ListFullPageFetcher returns a paging.PageFetcher for rules with all fields.
func (c *RuleClient) ListFullPageFetcher(params RuleListParams) paging.PageFetcher[RuleGetResponse] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[RuleGetResponse], error) {
		return c.ListFull(ctx, pageParams, params)
	}
}

// Get returns a rule by identifier.
func (c *RuleClient) Get(ctx context.Context, identifier string) (RuleGetResponse, error) {
	resp := RuleGetResponse{}
	err := c.transport.GetSingle(ctx, fmt.Sprintf("/api/LBaaS/v2/rule/%s", identifier), &resp)
	return resp, common.MapTransportError(err)
}
