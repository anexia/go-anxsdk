package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	go_anx_sdk "github.com/kkostial/go-anx-sdk"
	"github.com/kkostial/go-anx-sdk/config"
	"github.com/kkostial/go-anx-sdk/paging"
	"github.com/kkostial/go-anx-sdk/utils"
	v1 "github.com/kkostial/go-anx-sdk/v1"
)

func main() {
	httpClient := http.DefaultClient
	httpClient.Transport = utils.NewRateLimitRoundTripper(httpClient.Transport, time.Second*10, 10)

	client := go_anx_sdk.NewClient(config.WithAPIKey(os.Getenv("ANEXIA_TOKEN")), config.WithHTTPClient(httpClient))

	iter := paging.PaginateAndLoad(context.Background(), client.V1().Clusters().ListPageFetcher(v1.ClusterListParams{}), client.V1().Clusters().Get)

	_, _ = fmt.Printf("%s\t%s\n", "Identifier                      ", "Name")

	for item, itemErr := range iter {
		if itemErr != nil {
			panic(itemErr)
		}
		_, _ = fmt.Printf("%s\t%s\n", item.Identifier, item.Name)
	}
}
