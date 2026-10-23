package main

import (
	"context"
	"fmt"

	anxsdkfake "github.com/anexia/go-anxsdk/fake"
	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/v1/core"
)

func listFakeLocations() {
	ctx := context.Background()

	// NewClient returns a fully pre-wired fake client tree. Configure the fake
	// clients you care about via the getters.
	client := anxsdkfake.NewClient()
	client.Locations().ListReturns(
		paging.PagedResponse[core.LocationListItem]{
			Data: []core.LocationListItem{
				{Identifier: "anx01-identifier", Name: "ANX01 - Klagenfurt"},
				{Identifier: "anx04-identifier", Name: "ANX04 - Wien"},
			},
		},
		nil,
	)

	// The fake client implements the same interfaces as the real client.
	page, err := client.V1().Core().Locations().List(ctx, paging.DefaultParams(), core.LocationListParams{})
	if err != nil {
		panic(err)
	}

	for _, loc := range page.Data {
		_, _ = fmt.Printf("%s\t%s\n", loc.Identifier, loc.Name)
	}
}
