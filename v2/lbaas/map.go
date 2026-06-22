package lbaas

import (
	"context"
	"fmt"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/v2/common"
)

// MapListParams defines the available parameters for the map list endpoint.
type MapListParams struct {
	Search string `url:"search,omitempty"`
}

// MapListItem is an item in the map list response.
type MapListItem struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

// MapGetResponse represents the response of the map get endpoint.
type MapGetResponse struct {
	Identifier string `json:"identifier"`
}

// MapClient is an api client for managing load balancer maps.
type MapClient struct {
	transport *internal.Transport
}

// NewMapClient creates a new map client.
func NewMapClient(transport *internal.Transport) *MapClient {
	return &MapClient{
		transport: transport,
	}
}

// List returns a list of paged maps.
func (c *MapClient) List(ctx context.Context, pagingParams paging.Params, params MapListParams) (paging.PagedResponse[MapListItem], error) {
	resp := paging.PagedResponse[MapListItem]{}
	err := c.transport.Get(ctx, "/api/LBaaS/v2/map", &resp, pagingParams, params)
	return resp, common.MapTransportError(err)
}

// ListPageFetcher returns a paging.PageFetcher for maps.
func (c *MapClient) ListPageFetcher(params MapListParams) paging.PageFetcher[MapListItem] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[MapListItem], error) {
		return c.List(ctx, pageParams, params)
	}
}

// ListFull returns a list of paged maps with all fields.
func (c *MapClient) ListFull(ctx context.Context, pagingParams paging.Params, params MapListParams) (paging.PagedResponse[MapGetResponse], error) {
	resp := paging.PagedResponse[MapGetResponse]{}
	wrap := internal.NewAllAttributesWrapper(params)
	err := c.transport.Get(ctx, "/api/LBaaS/v2/map", &resp, pagingParams, wrap)
	return resp, common.MapTransportError(err)
}

// ListFullPageFetcher returns a paging.PageFetcher for maps with all fields.
func (c *MapClient) ListFullPageFetcher(params MapListParams) paging.PageFetcher[MapGetResponse] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[MapGetResponse], error) {
		return c.ListFull(ctx, pageParams, params)
	}
}

// Get returns a map by identifier.
func (c *MapClient) Get(ctx context.Context, identifier string) (MapGetResponse, error) {
	resp := MapGetResponse{}
	err := c.transport.GetSingle(ctx, fmt.Sprintf("/api/LBaaS/v2/map/%s", identifier), &resp)
	return resp, common.MapTransportError(err)
}
