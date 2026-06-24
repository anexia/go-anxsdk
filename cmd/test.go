package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/anexia/go-anxsdk"
	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/utils"
	v2 "github.com/anexia/go-anxsdk/v2"
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

	clusterClient := client.V2().KubernetesDev()
	clusters := paging.Paginate(ctx, clusterClient.ListFullPageFetcher(v2.ClusterListParams{
		Name:  "a",
		State: "Foo",
	}))

	for cluster, err := range clusters {
		if err != nil {
			panic(err)
		}

		fmt.Printf("%+v\n", cluster)
	}
}
