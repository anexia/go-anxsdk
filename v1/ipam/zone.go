package ipam

import (
	"context"
	"fmt"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/v1/common"
)

// ZoneListParams defines the available parameters for the zone list endpoint.
type ZoneListParams struct {
	Search string `url:"search,omitempty"`
}

// ZoneFilteredParams defines the available parameters for the zone filter endpoint.
type ZoneFilteredParams struct {
	Search   string `url:"search,omitempty"`
	Name     string `url:"name,omitempty"`
	Location string `url:"location,omitempty"`
}

// ZoneListItem is an item in the zone list response.
type ZoneListItem struct {
	Identifier string `url:"identifier"`
	Name       string `url:"name"`
}

// GetID returns the Identifier of the ZoneListItem.
func (i ZoneListItem) GetID() string {
	return i.Identifier
}

// ZoneGetResponse represents the response of the zone get endpoint.
type ZoneGetResponse struct {
	Identifier  string   `json:"identifier"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Location    Location `json:"location"`
}

// ZoneClient is an api client for managing zones.
type ZoneClient struct {
	transport *internal.Transport
}

// NewZoneClient creates a new zone client.
func NewZoneClient(transport *internal.Transport) *ZoneClient {
	return &ZoneClient{
		transport: transport,
	}
}

// List returns a list of paged zones.
func (c *ZoneClient) List(ctx context.Context, pagingParams paging.Params, params ZoneListParams) (paging.PagedResponse[ZoneListItem], error) {
	resp := internal.RequestWrapper[paging.PagedResponse[ZoneListItem]]{}
	err := c.transport.Get(ctx, "/api/ipam/v1/zone.json", &resp, pagingParams, params)
	return resp.Data, common.MapTransportError(err)
}

// ListPageFetcher returns a paging.PageFetcher for zones.
func (c *ZoneClient) ListPageFetcher(params ZoneListParams) paging.PageFetcher[ZoneListItem] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[ZoneListItem], error) {
		return c.List(ctx, pageParams, params)
	}
}

// ListFiltered returns a paged list of zones filtered by the provided parameters.
func (c *ZoneClient) ListFiltered(ctx context.Context, pagingParams paging.Params, params ZoneFilteredParams) (paging.PagedResponse[ZoneListItem], error) {
	resp := internal.RequestWrapper[paging.PagedResponse[ZoneListItem]]{}
	err := c.transport.Get(ctx, "/api/ipam/v1/zone/filtered.json", &resp, pagingParams, params)
	return resp.Data, common.MapTransportError(err)
}

// ListFilteredPageFetcher returns a paging.PageFetcher for filtered zones.
func (c *ZoneClient) ListFilteredPageFetcher(params ZoneFilteredParams) paging.PageFetcher[ZoneListItem] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[ZoneListItem], error) {
		return c.ListFiltered(ctx, pageParams, params)
	}
}

// Get returns an zone by identifier.
func (c *ZoneClient) Get(ctx context.Context, identifier string) (ZoneGetResponse, error) {
	resp := ZoneGetResponse{}
	err := c.transport.GetSingle(ctx, fmt.Sprintf("/api/ipam/v1/zone.json/%s", identifier), &resp)
	return resp, common.MapTransportError(err)
}
