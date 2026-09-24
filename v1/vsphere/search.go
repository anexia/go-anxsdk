package vsphere

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/v1/common"
)

// SearchResult represents a single vm matched by a search request.
type SearchResult struct {
	Name            string `json:"name"`
	Identifier      string `json:"identifier"`
	LocationCode    string `json:"location_code"`
	LocationCountry string `json:"location_country"`
	LocationName    string `json:"location_name"`
	PrimaryIPv4     string `json:"ip_v4_primary"`
	PrimaryIPv6     string `json:"ip_v6_primary"`
	OSName          string `json:"os_name"`
	OSFamily        string `json:"os_family"`
	// Tags are deserialized from a comma-separated string returned by the api.
	Tags []string `json:"-"`
}

// UnmarshalJSON deserializes a SearchResult, splitting the api's comma-separated tags string into a slice.
func (s *SearchResult) UnmarshalJSON(data []byte) error {
	type searchResult SearchResult
	if err := json.Unmarshal(data, (*searchResult)(s)); err != nil {
		return fmt.Errorf("unmarshalling SearchResult: %w", err)
	}

	var tags struct {
		Tags string `json:"tags"`
	}
	if err := json.Unmarshal(data, &tags); err != nil {
		return fmt.Errorf("unmarshalling SearchResult: %w", err)
	}

	if tags.Tags != "" {
		s.Tags = strings.Split(tags.Tags, ",")
	}

	return nil
}

// SearchByNameParams defines the available parameters for the by_name search endpoint.
type SearchByNameParams struct {
	// Name is matched against the vm name and may contain wildcards '%' and '_'.
	Name string `url:"name,omitempty"`
}

// SearchByTagsParams defines the available parameters for the by_tags search endpoint.
type SearchByTagsParams struct {
	Tags []string `url:"tags,comma,omitempty"`
}

// SearchClient is an api client for searching vms.
type SearchClient struct {
	transport *internal.Transport
}

func newSearchClient(transport *internal.Transport) *SearchClient {
	return &SearchClient{
		transport: transport,
	}
}

// ByName returns a paged response of vms matching the given name parameters.
func (c *SearchClient) ByName(ctx context.Context, pageParams paging.Params, params SearchByNameParams) (paging.PagedResponse[SearchResult], error) {
	resp := paging.PagedResponse[SearchResult]{}
	err := c.transport.Get(ctx, "/api/vsphere/v1/search/by_name.json", &resp, pageParams, params)
	return resp, common.MapTransportError(err)
}

// ByNamePageFetcher returns a paging.PageFetcher for vms matching the given name parameters.
func (c *SearchClient) ByNamePageFetcher(params SearchByNameParams) paging.PageFetcher[SearchResult] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[SearchResult], error) {
		return c.ByName(ctx, pageParams, params)
	}
}

// ByTags returns a paged response of vms matching the given tags parameters.
func (c *SearchClient) ByTags(ctx context.Context, pageParams paging.Params, params SearchByTagsParams) (paging.PagedResponse[SearchResult], error) {
	resp := paging.PagedResponse[SearchResult]{}
	err := c.transport.Get(ctx, "/api/vsphere/v1/search/by_tags.json", &resp, pageParams, params)
	return resp, common.MapTransportError(err)
}

// ByTagsPageFetcher returns a paging.PageFetcher for vms matching the given tags parameters.
func (c *SearchClient) ByTagsPageFetcher(params SearchByTagsParams) paging.PageFetcher[SearchResult] {
	return func(ctx context.Context, pageParams paging.Params) (paging.PagedResponse[SearchResult], error) {
		return c.ByTags(ctx, pageParams, params)
	}
}
