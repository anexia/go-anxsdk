package paging

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//
// Helpers
//

type testItem struct {
	identifier string
}

func (i testItem) GetID() string {
	return i.identifier
}

type itemFetchError struct{}

func (itemFetchError) Error() string {
	return "fetching item failed"
}

var errItemFetch = itemFetchError{}

func newTestPageFetcher(pages [][]testItem) PageFetcher[testItem] {
	return func(_ context.Context, params Params) (PagedResponse[testItem], error) {
		return PagedResponse[testItem]{
			Page:       params.Page,
			TotalPages: len(pages),
			Limit:      params.Limit,
			Data:       pages[params.Page-EngineFirstPage],
		}, nil
	}
}

func newTestItemFetcher(failOn string) ItemFetcher[string] {
	return func(_ context.Context, id string) (string, error) {
		if id == failOn {
			return "", errItemFetch
		}

		return "loaded-" + id, nil
	}
}

//
// Tests
//

func TestPaginateAndLoad_StopsWhenConsumerBreaks(t *testing.T) {
	// arrange
	pages := [][]testItem{{{"a"}, {"b"}}, {{"c"}}}
	loaded := make([]string, 0, 1)

	// act
	broke := func() {
		for item := range PaginateAndLoad(context.Background(), newTestPageFetcher(pages), newTestItemFetcher("")) {
			loaded = append(loaded, item)
			break
		}
	}

	// assert
	assert.NotPanics(t, broke)
	assert.Equal(t, []string{"loaded-a"}, loaded)
}

func TestPaginateAndLoad_StopsOnItemError(t *testing.T) {
	// arrange
	pages := [][]testItem{{{"a"}, {"b"}}, {{"c"}}}

	// act
	items, err := CollectAll(PaginateAndLoad(context.Background(), newTestPageFetcher(pages), newTestItemFetcher("b")))

	// assert
	require.ErrorIs(t, err, errItemFetch)
	assert.Nil(t, items)
}

func TestPaginateAndLoad_LoadsEveryPage(t *testing.T) {
	// arrange
	pages := [][]testItem{{{"a"}, {"b"}}, {{"c"}}}

	// act
	items, err := CollectAll(PaginateAndLoad(context.Background(), newTestPageFetcher(pages), newTestItemFetcher("")))

	// assert
	require.NoError(t, err)
	assert.Equal(t, []string{"loaded-a", "loaded-b", "loaded-c"}, items)
}
