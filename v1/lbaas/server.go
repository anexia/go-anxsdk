package lbaas

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/v1/common"
)

// ServerFilters is a struct that represents all filterable fields of a server.
type ServerFilters struct {
	Backend *string
}

// EncodeValues implements the query.Encode interface for ServerFilters.
func (f *ServerFilters) EncodeValues(key string, v *url.Values) error {
	sb := strings.Builder{}

	if f.Backend != nil {
		_, _ = fmt.Fprintf(&sb, "backend=%s", *f.Backend)
	}

	v.Add(key, sb.String())
	return nil
}

// ServerListParams defines the available parameters for the server list endpoint.
type ServerListParams struct {
	Search  string         `url:"search,omitempty"`
	Filters *ServerFilters `url:"filters,omitempty"`
}

// ServerListItem is an item in the server list response.
type ServerListItem struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

// GetID returns the Itentifier of the ServerListItem.
func (i ServerListItem) GetID() string {
	return i.Identifier
}

// ServerGetResponse represents the response of the server get endpoint.
type ServerGetResponse struct {
	CustomerIdentifier         *string              `json:"customer_identifier"`
	ResellerIdentifier         string               `json:"reseller_identifier"`
	CriticalOperationPassword  *string              `json:"critical_operation_password"`
	CriticalOperationConfirmed bool                 `json:"critical_operation_confirmed"`
	Identifier                 string               `json:"identifier"`
	Name                       string               `json:"name"`
	State                      common.State[string] `json:"state"`
	IP                         string               `json:"ip"`
	Port                       int                  `json:"port"`
	Backend                    common.Resource      `json:"backend"`
	Check                      string               `json:"check"`
	Redeploy                   bool                 `json:"redeploy"`
	AutomationRules            []common.Resource    `json:"automation_rules"`
}

// ServerClient is an api client for managing load balancer servers.
type ServerClient struct {
	transport *internal.Transport
}

// NewServerClient creates a new server client.
func NewServerClient(transport *internal.Transport) *ServerClient {
	return &ServerClient{
		transport: transport,
	}
}

// List returns a list of paged servers.
func (c *ServerClient) List(ctx context.Context, pagingParams paging.Params, params ServerListParams) (paging.PagedResponse[ServerListItem], error) {
	resp := internal.RequestWrapper[paging.PagedResponse[ServerListItem]]{}
	err := c.transport.Get(ctx, "/api/LBaaS/v1/server.json", &resp, pagingParams, params)
	return resp.Data, common.MapTransportError(err)
}

// ListPageFetcher returns a paging.PageFetcher for servers.
func (c *ServerClient) ListPageFetcher(params ServerListParams) paging.PageFetcher[ServerListItem] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[ServerListItem], error) {
		return c.List(ctx, pageParams, params)
	}
}

// Get returns a server by identifier.
func (c *ServerClient) Get(ctx context.Context, identifier string) (ServerGetResponse, error) {
	resp := ServerGetResponse{}
	err := c.transport.GetSingle(ctx, fmt.Sprintf("/api/LBaaS/v1/server.json/%s", identifier), &resp)
	return resp, common.MapTransportError(err)
}
