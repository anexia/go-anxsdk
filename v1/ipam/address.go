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

//revive:disable // self-explanatory constants

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

//revive:enable

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
	AssignedResourceID   string      `json:"assigned_resource_id"`
	RdnsName             *string     `json:"rdns_name"`
}

// AddressCreateRequest defines all fields available when creating a new address.
type AddressCreateRequest struct {
	Prefix              string      `json:"prefix"`
	Name                string      `json:"name"`
	DescriptionCustomer string      `json:"description_customer,omitempty"`
	Role                AddressRole `json:"role"`
	RdnsName            string      `json:"rdns_name,omitempty"`
}

// AddressCreateResponse defines the response when creating or reserving an address.
type AddressCreateResponse struct {
	Identifier           string  `json:"identifier"`
	Name                 string  `json:"name"`
	DescriptionCustomer  *string `json:"description_customer"`
	RoleText             string  `json:"role_text"`
	AssignedResourceName *string `json:"assigned_resource_name"`
	AssignedResourceID   *string `json:"assigned_resource_id"`
	RdnsName             *string `json:"rdns_name"`
}

// AddressReserveRandomRequest defines all fields available when reserving random addresses.
type AddressReserveRandomRequest struct {
	Count              int             `json:"count"`
	LocationIdentifier string          `json:"location_identifier"`
	VlanIdentifier     string          `json:"vlan_identifier"`
	PrefixIdentifier   *string         `json:"prefix_identifier,omitempty"`
	IPVersion          *AddressVersion `json:"ip_version,omitempty"`
	ReservationPeriod  *int            `json:"reservation_period,omitempty"`
}

// AddressReserveSpecificRequest defines all fields available when reserving a specific address.
type AddressReserveSpecificRequest struct {
	LocationIdentifier string   `json:"location_identifier"`
	VlanIdentifier     string   `json:"vlan_identifier"`
	Ips                []string `json:"ips"`
	ReservationPeriod  *int     `json:"reservation_period,omitempty"`
}

// AddressReserveResponseItem represents the response of a single reserved address.
type AddressReserveResponseItem struct {
	Identifier    string `json:"identifier"`
	Text          string `json:"text"`
	Prefix        string `json:"prefix"`
	ManagedStatus string `json:"managed_status"`
}

// AddressUpdateRequest defines the possible values that can be updated in an address. Nil values are ignored.
type AddressUpdateRequest struct {
	DescriptionCustomer *string      `json:"description_customer,omitempty"`
	Role                *AddressRole `json:"role,omitempty"`
	RdnsName            *string      `json:"rdns_name,omitempty"`
}

// AddressUpdateResponse is the response of an address update operation.
type AddressUpdateResponse struct {
	Identifier string `json:"identifier"`
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

// Create creates a new address.
func (c *AddressClient) Create(ctx context.Context, request AddressCreateRequest) (AddressCreateResponse, error) {
	resp := AddressCreateResponse{}
	err := c.transport.Post(ctx, "/api/ipam/v1/address.json", request, &resp)
	return resp, common.MapTransportError(err)
}

// ReserveRandom reserves random ip addresses.
func (c *AddressClient) ReserveRandom(ctx context.Context, request AddressReserveRandomRequest) (paging.PagedResponse[AddressReserveResponseItem], error) {
	resp := paging.PagedResponse[AddressReserveResponseItem]{}
	err := c.transport.Post(ctx, "/api/ipam/v1/address/ip/count.json", request, &resp)
	return resp, common.MapTransportError(err)
}

// ReserveSpecific reserves specific ip addresses.
func (c *AddressClient) ReserveSpecific(ctx context.Context, request AddressReserveSpecificRequest) (paging.PagedResponse[AddressReserveResponseItem], error) {
	resp := paging.PagedResponse[AddressReserveResponseItem]{}
	err := c.transport.Post(ctx, "/api/ipam/v1/address/ip/specific.json", request, &resp)
	return resp, common.MapTransportError(err)
}

// Update updates an address by identifier.
func (c *AddressClient) Update(ctx context.Context, identifier string, request AddressUpdateRequest) (AddressUpdateResponse, error) {
	resp := AddressUpdateResponse{}
	err := c.transport.Put(ctx, fmt.Sprintf("/api/ipam/v1/address.json/%s", identifier), request, &resp)
	return resp, common.MapTransportError(err)
}

// Delete deletes an address by identifier.
func (c *AddressClient) Delete(ctx context.Context, identifier string) error {
	err := c.transport.Delete(ctx, fmt.Sprintf("/api/ipam/v1/address.json/%s", identifier))
	return common.MapTransportError(err)
}
