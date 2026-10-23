package corefake_test

import (
	"context"
	"testing"

	core "github.com/anexia/go-anxsdk/fake/v1/core"
	"github.com/anexia/go-anxsdk/paging"
	v1core "github.com/anexia/go-anxsdk/v1/core"
)

func TestFakeLocationsAPI(t *testing.T) {
	fake := &core.FakeLocationsAPI{}
	fake.ListReturns(
		paging.PagedResponse[v1core.LocationListItem]{
			Data: []v1core.LocationListItem{{Identifier: "loc-1"}},
		},
		nil,
	)

	params := v1core.LocationListParams{Search: "anx"}
	got, err := fake.List(context.Background(), paging.DefaultParams(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got.Data) != 1 || got.Data[0].Identifier != "loc-1" {
		t.Fatalf("unexpected response: %+v", got.Data)
	}

	if fake.ListCallCount() != 1 {
		t.Fatalf("expected 1 List call, got %d", fake.ListCallCount())
	}

	if _, _, gotParams := fake.ListArgsForCall(0); gotParams != params {
		t.Fatalf("unexpected params: %+v", gotParams)
	}
}

func TestFakeAPITree(t *testing.T) {
	locations := &core.FakeLocationsAPI{}
	locations.ListReturns(
		paging.PagedResponse[v1core.LocationListItem]{
			Data: []v1core.LocationListItem{{Identifier: "loc-1"}},
		},
		nil,
	)

	api := &core.FakeAPI{}
	api.LocationsReturns(locations)

	var coreAPI v1core.API = api
	got, err := coreAPI.Locations().List(context.Background(), paging.DefaultParams(), v1core.LocationListParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got.Data) != 1 {
		t.Fatalf("unexpected response: %+v", got.Data)
	}

	if api.LocationsCallCount() != 1 {
		t.Fatalf("expected 1 Locations call, got %d", api.LocationsCallCount())
	}
}
