package v1

import (
	"context"
	"fmt"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/paging"
)

// VlanCreateRequest defines all fields available when creating a new vlan.
type VlanCreateRequest struct {
	LocationIdentifier  string `json:"location"`
	VMProvisioning      bool   `json:"vm_provisioning"`
	DescriptionCustomer string `json:"description_customer"`
}

// VlanCreateResponse defines the response for a vlan create request.
type VlanCreateResponse struct {
	Identifier          string `json:"identifier"`
	Name                string `json:"name"`
	DescriptionCustomer string `json:"description_customer"`
}

// VlanGetResponse represents a full vlan response.
type VlanGetResponse struct {
	Identifier          string                        `json:"identifier"`
	Name                string                        `json:"name"`
	DescriptionCustomer string                        `json:"description_customer"`
	DescriptionInternal string                        `json:"description_internal"`
	RoleText            string                        `json:"role_text"`
	Status              string                        `json:"status"`
	Locations           []VlanGetResponseLocationItem `json:"locations"`
	VMProvisioning      bool                          `json:"vm_provisioning"`
}

// VlanGetResponseLocationItem represents a location in a full vlan response.
type VlanGetResponseLocationItem struct {
	Identifier string `json:"identifier"`
	Code       string `json:"code"`
	Name       string `json:"name"`
	Country    string `json:"country"`
	Lat        string `json:"lat"`
	Lon        string `json:"lon"`
	CityCode   string `json:"city_code"`
}

// VlanUpdateRequest defines the possible values that can be updated in a vlan. Nil values are ignored.
type VlanUpdateRequest struct {
	VMProvisioning      *bool   `json:"vm_provisioning,omitempty"`
	DescriptionCustomer *string `json:"description_customer,omitempty"`
}

// VlanUpdateResponse is the response of a vlan update operation.
type VlanUpdateResponse struct {
	Identifier          string `json:"identifier"`
	Name                string `json:"name"`
	DescriptionCustomer string `json:"description_customer"`
}

// VlanListParams defines the available parameters for the vlan list endpoint.
type VlanListParams struct {
	Search string `url:"search,omitempty"`
}

// VlanFilteredParams defines the available parameters for the vlan filter endpoint.
type VlanFilteredParams struct {
	Search                 string `url:"search,omitempty"`
	OrganizationIdentifier string `url:"organization_identifier,omitempty"`
	RoleText               string `url:"role_text,omitempty"`
	Status                 string `url:"status,omitempty"`
	LocationIdentifier     string `url:"location,omitempty"`
}

// VlanListItem is an item in vlan list responses.
type VlanListItem struct {
	Identifier          string `json:"identifier"`
	Name                string `json:"name"`
	DescriptionCustomer string `json:"description_customer"`
}

// GetID returns the Identifier of the [VlanListItem].
func (v VlanListItem) GetID() string {
	return v.Identifier
}

// VlansClient is an api client for managing vlans.
type VlansClient struct {
	transport *internal.Transport
}

// NewVlansClient creates a new vlans client.
func NewVlansClient(transport *internal.Transport) *VlansClient {
	return &VlansClient{
		transport: transport,
	}
}

// Create creates a new vlan.
func (v *VlansClient) Create(ctx context.Context, request VlanCreateRequest) (VlanCreateResponse, error) {
	resp := VlanCreateResponse{}
	err := v.transport.Post(ctx, "/api/vlan/v1/vlan.json", request, &resp)
	return resp, mapTransportError(err)
}

// Get returns a vlan by identifier.
func (v *VlansClient) Get(ctx context.Context, identifier string) (VlanGetResponse, error) {
	resp := VlanGetResponse{}
	err := v.transport.GetSingle(ctx, fmt.Sprintf("api/vlan/v1/vlan.json/%s", identifier), &resp)
	return resp, mapTransportError(err)
}

// Update updates a vlan by identifier.
func (v *VlansClient) Update(ctx context.Context, identifier string, request VlanUpdateRequest) (VlanUpdateResponse, error) {
	resp := VlanUpdateResponse{}
	err := v.transport.Put(ctx, fmt.Sprintf("/api/vlan/v1/vlan.json/%s", identifier), request, &resp)
	return resp, mapTransportError(err)
}

// List returns a paged list of vlans.
func (v *VlansClient) List(ctx context.Context, pageParams paging.Params, params VlanListParams) (paging.PagedResponse[VlanListItem], error) {
	resp := internal.RequestWrapper[paging.PagedResponse[VlanListItem]]{}
	err := v.transport.Get(ctx, "/api/vlan/v1/vlan.json", &resp, pageParams, params)
	return resp.Data, mapTransportError(err)
}

// ListPageFetcher returns a paging.PageFetcher for vlans.
func (v *VlansClient) ListPageFetcher(params VlanListParams) paging.PageFetcher[VlanListItem] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[VlanListItem], error) {
		return v.List(ctx, pageParams, params)
	}
}

// ListFiltered returns a paged list of vlans filtered by the provided parameters.
func (v *VlansClient) ListFiltered(ctx context.Context, pageParams paging.Params, params VlanFilteredParams) (paging.PagedResponse[VlanListItem], error) {
	resp := internal.RequestWrapper[paging.PagedResponse[VlanListItem]]{}
	err := v.transport.Get(ctx, "/api/vlan/v1/vlan/filtered.json", &resp, pageParams, params)
	return resp.Data, mapTransportError(err)
}

// ListFilteredPageFetcher returns a paging.PageFetcher for filtered vlans.
func (v *VlansClient) ListFilteredPageFetcher(params VlanFilteredParams) paging.PageFetcher[VlanListItem] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[VlanListItem], error) {
		return v.ListFiltered(ctx, pageParams, params)
	}
}

// Delete deletes a vlan by identifier.
func (v *VlansClient) Delete(ctx context.Context, identifier string) error {
	err := v.transport.Delete(ctx, fmt.Sprintf("/api/vlan/v1/vlan.json/%s", identifier))
	return mapTransportError(err)
}
