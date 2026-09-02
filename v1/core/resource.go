package core

import (
	"context"
	"fmt"
	"time"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/v1/common"
)

// ResourceListParams defines the available parameters for the resource list endpoint.
type ResourceListParams struct {
	IncludeSoftDelete      *bool   `url:"include_soft_delete"`
	ResellerIdentifier     *string `url:"reseller_identifier"`
	CustomerIdentifier     *string `url:"customer_identifier"`
	TagName                *string `url:"tag_name"`
	ResourcePoolIdentifier *string `url:"resource_pool_identifier"`
}

// ResourceListItem is an item in the resource list response.
type ResourceListItem struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

// ResourceGetResponse represents the details of a core resource.
type ResourceGetResponse struct {
	Name            string            `json:"name"`
	Identifier      string            `json:"identifier"`
	ResourceType    common.Resource   `json:"resource_type"`
	ServiceName     string            `json:"service_name"`
	CreatedAt       time.Time         `json:"created_at"`
	DeletedAt       *time.Time        `json:"deleted_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	Reseller        Reseller          `json:"reseller"`
	Customer        Customer          `json:"customer"`
	BillingContract *string           `json:"billing_contract"`
	ManagedStatus   string            `json:"managed_status"`
	SharedBy        *string           `json:"shared_by"`
	SharedAt        *string           `json:"shared_at"`
	ResourcePools   []common.Resource `json:"resource_pools"`
	Tags            []common.Resource `json:"tags"`
}

// Reseller represents the reseller from the resource.
type Reseller struct {
	CustomerID string  `json:"customer_id"`
	Demo       bool    `json:"demo"`
	Identifier string  `json:"identifier"`
	Name       string  `json:"name"`
	NameSlug   string  `json:"name_slug"`
	Reseller   *string `json:"reseller"`
}

// Customer represents the customer from the resource.
type Customer struct {
	CustomerID *string `json:"customer_id"`
	Demo       bool    `json:"demo"`
	Identifier string  `json:"identifier"`
	Name       string  `json:"name"`
	NameSlug   string  `json:"name_slug"`
	Reseller   *string `json:"reseller"`
}

// ResourceClient is an api client for managing resources.
type ResourceClient struct {
	transport *internal.Transport
}

// newResourceClient creates a new resource client.
func newResourceClient(transport *internal.Transport) *ResourceClient {
	return &ResourceClient{
		transport: transport,
	}
}

// List returns a list of paged resources.
func (v *ResourceClient) List(ctx context.Context, pagingParams paging.Params, params ResourceListParams) (paging.PagedResponse[ResourceListItem], error) {
	resp := paging.PagedResponse[ResourceListItem]{}
	err := v.transport.Get(ctx, "api/core/v1/resource.json", &resp, pagingParams, params)
	return resp, common.MapTransportError(err)
}

// Get returns a single resource by its identifier.
func (v *ResourceClient) Get(ctx context.Context, identifier string) (ResourceGetResponse, error) {
	resp := ResourceGetResponse{}
	err := v.transport.GetSingle(ctx, fmt.Sprintf("api/core/v1/resource.json/%s", identifier), &resp)
	return resp, common.MapTransportError(err)
}

// GetTags returns a list of all tags on a resource.
func (v *ResourceClient) GetTags(ctx context.Context, identifier string) ([]common.Resource, error) {
	var resp []common.Resource
	err := v.transport.GetSingle(ctx, fmt.Sprintf("api/core/v1/resource.json/%s/tags", identifier), &resp)
	return resp, common.MapTransportError(err)
}

// Tag tags a resource with the provided tag name.
func (v *ResourceClient) Tag(ctx context.Context, identifier string, tagName string) error {
	err := v.transport.Post(ctx, fmt.Sprintf("api/core/v1/resource.json/%s/tags/%s", identifier, tagName), nil, nil)
	return common.MapTransportError(err)
}

// Untag removes a tag from a resource.
func (v *ResourceClient) Untag(ctx context.Context, identifier string, tagName string) error {
	err := v.transport.Delete(ctx, fmt.Sprintf("api/core/v1/resource.json/%s/tags/%s", identifier, tagName))
	return common.MapTransportError(err)
}
