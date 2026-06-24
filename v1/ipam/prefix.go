package ipam

import (
	"context"
	"fmt"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/v1/common"
)

type PrefixType int

const (
	PrefixTypePublic  PrefixType = 0
	PrefixTypePrivate PrefixType = 1
)

// PrefixListParams defines the available parameters for the prefix list endpoint.
type PrefixListParams struct {
	Search string `url:"search,omitempty"`
}

// PrefixFilteredParams defines the available parameters for the prefix filtered endpoint.
type PrefixFilteredParams struct {
	Search                 string          `url:"search,omitempty"`
	OrganizationIdentifier string          `url:"organization_identifier,omitempty"`
	Version                *AddressVersion `url:"version,omitempty"`
	Private                *bool           `url:"private,omitempty"`
	Location               string          `url:"location,omitempty"`
	VlanIdentifier         string          `url:"vlan_identifier,omitempty"`
}

// PrefixListItem is an item in the prefix list response.
type PrefixListItem struct {
	Identifier          string           `json:"identifier"`
	Name                string           `json:"name"`
	DescriptionCustomer *string          `json:"description_customer"`
	RouterRedundancy    *string          `json:"router_redundancy"`
	Vlans               []PrefixVlanItem `json:"vlans"`
}

// PrefixVlanItem represents a vlan that a prefix is inside of.
type PrefixVlanItem struct {
	Identifier          string `json:"identifier"`
	Name                string `json:"name"`
	DescriptionCustomer string `json:"description_customer"`
}

// PrefixGetResponse represents the response of the prefix get endpoint.
type PrefixGetResponse struct {
	Identifier          string                          `json:"identifier"`
	Name                string                          `json:"name"`
	DescriptionCustomer *string                         `json:"description_customer"`
	DescriptionInternal string                          `json:"description_internal"`
	Version             AddressVersion                  `json:"version"`
	Netmask             int                             `json:"netmask"`
	ZoneOrdered         string                          `json:"zone_ordered"`
	AggregateOrdered    string                          `json:"aggregate_ordered"`
	Aggregate           string                          `json:"aggregate"`
	RoleText            string                          `json:"role_text"`
	Status              string                          `json:"status"`
	Locations           []PrefixGetResponseLocationItem `json:"locations"`
	RouterRedundancy    *bool                           `json:"router_redundancy"`
	Vlans               []PrefixVlanItem                `json:"vlans"`
	Type                PrefixType                      `json:"type"`
	Utilization         int                             `json:"utilization"`
}

// PrefixGetResponseLocationItem represents a location in a full prefix response.
type PrefixGetResponseLocationItem struct {
	Identifier string `json:"identifier"`
	Code       string `json:"code"`
	Name       string `json:"name"`
	Country    string `json:"country"`
	Lat        string `json:"lat"`
	Lon        string `json:"lon"`
	CityCode   string `json:"city_code"`
}

// PrefixClient is an api client for managing prefixes.
type PrefixClient struct {
	transport *internal.Transport
}

// NewPrefixClient creates a new prefix client.
func NewPrefixClient(transport *internal.Transport) *PrefixClient {
	return &PrefixClient{
		transport: transport,
	}
}

// List returns a list of paged aggregates.
func (c *PrefixClient) List(ctx context.Context, pagingParams paging.Params, params PrefixListParams) (paging.PagedResponse[PrefixListItem], error) {
	resp := internal.RequestWrapper[paging.PagedResponse[PrefixListItem]]{}
	err := c.transport.Get(ctx, "/api/ipam/v1/prefix.json", &resp, pagingParams, params)
	return resp.Data, common.MapTransportError(err)
}

// ListPageFetcher returns a paging.PageFetcher for aggregates.
func (c *PrefixClient) ListPageFetcher(params PrefixListParams) paging.PageFetcher[PrefixListItem] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[PrefixListItem], error) {
		return c.List(ctx, pageParams, params)
	}
}

// ListFiltered returns a paged list of aggregates filtered by the provided parameters.
func (c *PrefixClient) ListFiltered(ctx context.Context, pagingParams paging.Params, params PrefixFilteredParams) (paging.PagedResponse[PrefixListItem], error) {
	resp := internal.RequestWrapper[paging.PagedResponse[PrefixListItem]]{}
	err := c.transport.Get(ctx, "/api/ipam/v1/prefix/filtered.json", &resp, pagingParams, params)
	return resp.Data, common.MapTransportError(err)
}

// ListFilteredPageFetcher returns a paging.PageFetcher for filtered aggregates.
func (c *PrefixClient) ListFilteredPageFetcher(params PrefixFilteredParams) paging.PageFetcher[PrefixListItem] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[PrefixListItem], error) {
		return c.ListFiltered(ctx, pageParams, params)
	}
}

// Get returns an aggregate by identifier.
func (c *PrefixClient) Get(ctx context.Context, identifier string) (PrefixGetResponse, error) {
	resp := PrefixGetResponse{}
	err := c.transport.GetSingle(ctx, fmt.Sprintf("/api/ipam/v1/prefix.json/%s", identifier), &resp)
	return resp, common.MapTransportError(err)
}
