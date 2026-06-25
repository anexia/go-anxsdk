package ipam

import (
	"context"
	"fmt"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/v1/common"
)

// AggregateStatus represnts an aggregates status.
type AggregateStatus int

const (
	// AggregateStatusActive means that the aggregate is active.
	AggregateStatusActive AggregateStatus = 1
	// AggregateStatusInactive means that the aggregate is inactive.
	AggregateStatusInactive AggregateStatus = 2
)

// AggregateListParams defines the available parameters for the aggegate list endpoint.
type AggregateListParams struct {
	Search string `url:"search,omitempty"`
}

// AggregateFilteredParams defines the available parameters for the aggregate filtered endpoint.
type AggregateFilteredParams struct {
	Search                 string           `url:"search,omitempty"`
	Name                   string           `url:"name,omitempty"`
	Private                *bool            `url:"private,omitempty"`
	OrganizationIdentifier string           `url:"organization_identifier,omitempty"`
	Status                 *AggregateStatus `url:"status,omitempty"`
	IsFull                 *bool            `url:"is_full,omitempty"`
}

// AggregateListItem is an item in the aggregate list response.
type AggregateListItem struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

// GetID returns the Identifier of the AggregateListItem.
func (i AggregateListItem) GetID() string {
	return i.Identifier
}

// AggregateGetResponse represents the response of the aggregate get endpoint.
type AggregateGetResponse struct {
	Identifier  string         `json:"identifier"`
	Name        string         `json:"name"`
	Status      string         `json:"status"`
	Description string         `json:"description"`
	IsFull      bool           `json:"isFull"`
	Role        string         `json:"role"`
	Locations   []Location     `json:"locations"`
	Private     bool           `json:"private"`
	Version     AddressVersion `json:"version"`
}

// AggregateClient is an api client for managing aggregates.
type AggregateClient struct {
	transport *internal.Transport
}

// NewAggregateClient creates a new aggregate client.
func NewAggregateClient(transport *internal.Transport) *AggregateClient {
	return &AggregateClient{
		transport: transport,
	}
}

// List returns a list of paged aggregates.
func (c *AggregateClient) List(ctx context.Context, pagingParams paging.Params, params AggregateListParams) (paging.PagedResponse[AggregateListItem], error) {
	resp := internal.RequestWrapper[paging.PagedResponse[AggregateListItem]]{}
	err := c.transport.Get(ctx, "/api/ipam/v1/aggregate.json", &resp, pagingParams, params)
	return resp.Data, common.MapTransportError(err)
}

// ListPageFetcher returns a paging.PageFetcher for aggregates.
func (c *AggregateClient) ListPageFetcher(params AggregateListParams) paging.PageFetcher[AggregateListItem] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[AggregateListItem], error) {
		return c.List(ctx, pageParams, params)
	}
}

// ListFiltered returns a paged list of aggregates filtered by the provided parameters.
func (c *AggregateClient) ListFiltered(ctx context.Context, pagingParams paging.Params, params AggregateFilteredParams) (paging.PagedResponse[AggregateListItem], error) {
	resp := internal.RequestWrapper[paging.PagedResponse[AggregateListItem]]{}
	err := c.transport.Get(ctx, "/api/ipam/v1/aggregate/filtered.json", &resp, pagingParams, params)
	return resp.Data, common.MapTransportError(err)
}

// ListFilteredPageFetcher returns a paging.PageFetcher for filtered aggregates.
func (c *AggregateClient) ListFilteredPageFetcher(params AggregateFilteredParams) paging.PageFetcher[AggregateListItem] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[AggregateListItem], error) {
		return c.ListFiltered(ctx, pageParams, params)
	}
}

// Get returns an aggregate by identifier.
func (c *AggregateClient) Get(ctx context.Context, identifier string) (AggregateGetResponse, error) {
	resp := AggregateGetResponse{}
	err := c.transport.GetSingle(ctx, fmt.Sprintf("/api/ipam/v1/aggregate.json/%s", identifier), &resp)
	return resp, common.MapTransportError(err)
}
