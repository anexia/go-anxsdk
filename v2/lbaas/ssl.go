package lbaas

import (
	"context"
	"fmt"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/v2/common"
)

// SslListParams defines the available parameters for the ssl list endpoint.
type SslListParams struct {
	Search string `url:"search,omitempty"`
}

// SslListItem is an item in the ssl list response.
type SslListItem struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

// SslGetResponse represents the response of the ssl get endpoint.
type SslGetResponse struct {
	Identifier string `json:"identifier"`
}

// SslClient is an api client for managing load balancer ssls.
type SslClient struct {
	transport *internal.Transport
}

// NewSslClient creates a new ssl client.
func NewSslClient(transport *internal.Transport) *SslClient {
	return &SslClient{
		transport: transport,
	}
}

// List returns a list of paged ssls.
func (c *SslClient) List(ctx context.Context, pagingParams paging.Params, params SslListParams) (paging.PagedResponse[SslListItem], error) {
	resp := paging.PagedResponse[SslListItem]{}
	err := c.transport.Get(ctx, "/api/LBaaS/v2/ssl", &resp, pagingParams, params)
	return resp, common.MapTransportError(err)
}

// ListPageFetcher returns a paging.PageFetcher for ssls.
func (c *SslClient) ListPageFetcher(params SslListParams) paging.PageFetcher[SslListItem] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[SslListItem], error) {
		return c.List(ctx, pageParams, params)
	}
}

// ListFull returns a list of paged ssls with all fields.
func (c *SslClient) ListFull(ctx context.Context, pagingParams paging.Params, params SslListParams) (paging.PagedResponse[SslGetResponse], error) {
	resp := paging.PagedResponse[SslGetResponse]{}
	wrap := internal.NewAllAttributesWrapper(params)
	err := c.transport.Get(ctx, "/api/LBaaS/v2/ssl", &resp, pagingParams, wrap)
	return resp, common.MapTransportError(err)
}

// ListFullPageFetcher returns a paging.PageFetcher for ssls with all fields.
func (c *SslClient) ListFullPageFetcher(params SslListParams) paging.PageFetcher[SslGetResponse] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[SslGetResponse], error) {
		return c.ListFull(ctx, pageParams, params)
	}
}

// Get returns a ssl by identifier.
func (c *SslClient) Get(ctx context.Context, identifier string) (SslGetResponse, error) {
	resp := SslGetResponse{}
	err := c.transport.GetSingle(ctx, fmt.Sprintf("/api/LBaaS/v2/ssl/%s", identifier), &resp)
	return resp, common.MapTransportError(err)
}
