package lbaas

import (
	"context"
	"fmt"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/v1/common"
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

// GetID returns the Identifier of the MapListItem.
func (i MapListItem) GetID() string {
	return i.Identifier
}

// MapGetResponse represents the response of the map get endpoint.
type MapGetResponse struct {
	Identifier string `json:"identifier"`
}

//counterfeiter:generate -o ../../fake/v1/lbaas . MapAPI

// MapAPI is the interface implemented by [MapClient].
type MapAPI interface {
	// List returns a list of paged maps.
	List(ctx context.Context, pagingParams paging.Params, params MapListParams) (paging.PagedResponse[MapListItem], error)
	// ListPageFetcher returns a paging.PageFetcher for maps.
	ListPageFetcher(params MapListParams) paging.PageFetcher[MapListItem]
	// Get returns a map by identifier.
	Get(ctx context.Context, identifier string) (MapGetResponse, error)
}

// MapClient is an api client for managing load balancer maps.
type MapClient struct {
	transport *internal.Transport
}

// Check that MapClient implements the MapAPI interface.
var _ MapAPI = (*MapClient)(nil)

// NewMapClient creates a new map client.
func NewMapClient(transport *internal.Transport) *MapClient {
	return &MapClient{
		transport: transport,
	}
}

// List returns a list of paged maps.
func (c *MapClient) List(ctx context.Context, pagingParams paging.Params, params MapListParams) (paging.PagedResponse[MapListItem], error) {
	resp := internal.RequestWrapper[paging.PagedResponse[MapListItem]]{}
	err := c.transport.Get(ctx, "/api/LBaaS/v1/map.json", &resp, pagingParams, params)
	return resp.Data, common.MapTransportError(err)
}

// ListPageFetcher returns a paging.PageFetcher for maps.
func (c *MapClient) ListPageFetcher(params MapListParams) paging.PageFetcher[MapListItem] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[MapListItem], error) {
		return c.List(ctx, pageParams, params)
	}
}

// Get returns a map by identifier.
func (c *MapClient) Get(ctx context.Context, identifier string) (MapGetResponse, error) {
	resp := MapGetResponse{}
	err := c.transport.GetSingle(ctx, fmt.Sprintf("/api/LBaaS/v1/map.json/%s", identifier), &resp)
	return resp, common.MapTransportError(err)
}
