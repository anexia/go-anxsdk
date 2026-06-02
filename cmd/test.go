package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"code.anexia.com/se/ks/go-anxsdk"
	"code.anexia.com/se/ks/go-anxsdk/paging"
	"code.anexia.com/se/ks/go-anxsdk/utils"
	v1 "code.anexia.com/se/ks/go-anxsdk/v1"
)

func main() {
	ctx := context.Background()

	apiKey := os.Getenv("API_KEY")

	httpClient := &http.Client{
		Transport: utils.NewLoggingRoundTripper(http.DefaultTransport),
	}

	client := anxsdk.NewClient(
		anxsdk.WithAPIKey(apiKey),
		anxsdk.WithBaseURL("https://engine.anexia-it.com/"),
		anxsdk.WithHTTPClient(httpClient),
	)

	clusterClient := client.V1().DevClusters()
	clusters := paging.PaginateAndLoad(ctx, clusterClient.ListPageFetcher(v1.ClusterListParams{}), clusterClient.Get)
	for cluster, err := range clusters {
		if err != nil {
			panic(err)
		}

		fmt.Printf("%+v\n", cluster)
	}
}
