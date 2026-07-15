package core

import (
	"context"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/v1/common"
)

// LocationListParams defines the available parameters for the location list endpoint.
type LocationListParams struct {
	Search string `url:"search,omitempty"`
}

// LocationListItem is an item in the location list response.
type LocationListItem struct {
	Identifier string  `json:"identifier"`
	Code       string  `json:"code"`
	Name       string  `json:"name"`
	CityCode   *string `json:"city_code"`
	Country    *string `json:"country"`
	Lat        *string `json:"lat"`
	Lon        *string `json:"lon"`
}

// GetID returns the Identifier of the [LocationListItem].
func (l LocationListItem) GetID() string {
	return l.Identifier
}

// LocationsClient is an api client for managing locations.
type LocationsClient struct {
	transport *internal.Transport
}

// newLocationsClient creates a new location client.
func newLocationsClient(transport *internal.Transport) *LocationsClient {
	return &LocationsClient{
		transport,
	}
}

// List returns a list of paged locations.
func (v *LocationsClient) List(ctx context.Context, pagingParams paging.Params, params LocationListParams) (paging.PagedResponse[LocationListItem], error) {
	resp := internal.RequestWrapper[paging.PagedResponse[LocationListItem]]{}
	err := v.transport.Get(ctx, "/api/core/v1/location.json", &resp, pagingParams, params)
	return resp.Data, common.MapTransportError(err)
}

// ListPageFetcher returns a paging.PageFetcher for locations.
func (v *LocationsClient) ListPageFetcher(params LocationListParams) paging.PageFetcher[LocationListItem] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[LocationListItem], error) {
		return v.List(ctx, pageParams, params)
	}
}
