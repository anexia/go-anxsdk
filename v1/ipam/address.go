package ipam

import (
	"context"
	"fmt"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/v1/common"
)

// AddressRole represents a role of an address.
type AddressRole string

const (
	AddressRoleDefault  AddressRole = "Default"
	AddressRoleReserved AddressRole = "Reserved"
	AddressRoleGateway  AddressRole = "Gateway"
	AddressRoleRouter   AddressRole = "Router"
)

type AddressVersion int

const (
	AddressVersion4 AddressVersion = 4
	AddressVersion6 AddressVersion = 6
)

// AddressListParams defines the available parameters for the address list endpoint.
type AddressListParams struct {
	Search string `url:"search,omitempty"`
}

// AddressFilteredParams defines the available parameters for the address filter endpoint.
type AddressFilteredParams struct {
	Search                 string          `url:"search,omitempty"`
	OrganizationIdentifier string          `url:"organization_identifier,omitempty"`
	Prefix                 string          `url:"prefix,omitempty"`
	Vlan                   string          `url:"vlan,omitempty"`
	Version                *AddressVersion `url:"version,omitempty"`
	Private                *bool           `url:"private,omitempty"`
	Role                   AddressRole     `url:"role,omitempty"`
	Status                 string          `url:"status,omitempty"`
	Location               string          `url:"location,omitempty"`
	RdnsName               string          `url:"rdnsName,omitempty"`
}

// AddressListItem is an item in the address list response.
type AddressListItem struct {
	Identifier          string      `json:"identifier"`
	Name                string      `json:"name"`
	DescriptionCustomer string      `json:"description_customer"`
	Role                AddressRole `json:"role_text"`
	RdnsName            *string     `json:"rdns_name"`
}

// GetID returns the Identifier if the AddressListItem.
func (i AddressListItem) GetID() string {
	return i.Identifier
}

// AddressGetResponse represents the response of the address get endpoint.
type AddressGetResponse struct {
	Identifier           string      `json:"identifier"`
	Name                 string      `json:"name"`
	DescriptionCustomer  string      `json:"description_customer"`
	DescriptionInternal  string      `json:"description_internal"`
	Version              int         `json:"version"`
	Role                 AddressRole `json:"role_text"`
	Status               string      `json:"status"`
	Vlan                 string      `json:"vlan"`
	Prefix               string      `json:"prefix"`
	AssignedResourceName string      `json:"assigned_resource_name"`
	AssignedResourceId   string      `json:"assigned_resource_id"`
	RdnsName             interface{} `json:"rdns_name"`
}

// AddressClient is an api client for managing addresses.
type AddressClient struct {
	transport *internal.Transport
}

// NewAddressClient creates a new address client.
func NewAddressClient(transport *internal.Transport) *AddressClient {
	return &AddressClient{
		transport: transport,
	}
}

// List returns a list of paged addresses.
func (c *AddressClient) List(ctx context.Context, pagingParams paging.Params, params AddressListParams) (paging.PagedResponse[AddressListItem], error) {
	resp := internal.RequestWrapper[paging.PagedResponse[AddressListItem]]{}
	err := c.transport.Get(ctx, "/api/ipam/v1/address.json", &resp, pagingParams, params)
	return resp.Data, common.MapTransportError(err)
}

// ListPageFetcher returns a paging.PageFetcher for addresses.
func (c *AddressClient) ListPageFetcher(params AddressListParams) paging.PageFetcher[AddressListItem] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[AddressListItem], error) {
		return c.List(ctx, pageParams, params)
	}
}

// ListFiltered returns a paged list of addresses filtered by the provided parameters.
func (c *AddressClient) ListFiltered(ctx context.Context, pagingParams paging.Params, params AddressFilteredParams) (paging.PagedResponse[AddressListItem], error) {
	resp := internal.RequestWrapper[paging.PagedResponse[AddressListItem]]{}
	err := c.transport.Get(ctx, "/api/ipam/v1/address/filtered.json", &resp, pagingParams, params)
	return resp.Data, common.MapTransportError(err)
}

// ListFilteredPageFetcher returns a paging.PageFetcher for filtered addresses.
func (c *AddressClient) ListFilteredPageFetcher(params AddressFilteredParams) paging.PageFetcher[AddressListItem] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[AddressListItem], error) {
		return c.ListFiltered(ctx, pageParams, params)
	}
}

// Get returns an address by identifier.
func (c *AddressClient) Get(ctx context.Context, identifier string) (AddressGetResponse, error) {
	resp := AddressGetResponse{}
	err := c.transport.GetSingle(ctx, fmt.Sprintf("/api/ipam/v1/address.json/%s", identifier), &resp)
	return resp, common.MapTransportError(err)
}
