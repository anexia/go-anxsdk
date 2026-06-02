package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	anxsdk "code.anexia.com/se/ks/go-anxsdk"
	"code.anexia.com/se/ks/go-anxsdk/config"
	"code.anexia.com/se/ks/go-anxsdk/paging"
	"code.anexia.com/se/ks/go-anxsdk/utils"
	v2 "code.anexia.com/se/ks/go-anxsdk/v2"
)

func main() {
	httpClient := http.DefaultClient
	httpClient.Transport = utils.NewRateLimitRoundTripper(httpClient.Transport, time.Second*10, 10)

	client := anxsdk.NewClient(config.WithAPIKey(os.Getenv("ANEXIA_TOKEN")), config.WithHTTPClient(httpClient))

	ctx := context.Background()

	list, err := client.V2().Clusters().List(ctx, paging.Params{Page: 1, Limit: 3}, v2.ClusterListParams{})

	_, _ = fmt.Printf("%s\t%s\n", "Identifier                      ", "Name")
	if err != nil {
		panic(err)
	}

	for _, item := range list.Data {
		_, _ = fmt.Printf("%s\t%s\n", item.Identifier, item.Name)
	}

	iter := paging.PaginateAndLoad(context.Background(), client.V2().DevClusters().ListPageFetcher(v2.ClusterListParams{}), client.V2().DevClusters().Get)

	_, _ = fmt.Printf("%s\t%s\n", "Identifier                      ", "Name")

	count := 0

	for item, itemErr := range iter {
		if itemErr != nil {
			panic(itemErr)
		}
		_, _ = fmt.Printf("%s\t%s\n", item.Identifier, item.Name)

		count++

		if count > 5 {
			break
		}
	}
}
