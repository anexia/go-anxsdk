package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/anexia/go-anxsdk"
	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/utils"
	"github.com/anexia/go-anxsdk/v2/kubernetes"
)

func nodepools() {
	httpClient := http.DefaultClient
	httpClient.Transport = utils.NewRateLimitRoundTripper(httpClient.Transport, time.Second*10, 10)

	client := anxsdk.NewClient(anxsdk.WithAPIKey(os.Getenv("ANEXIA_TOKEN")), anxsdk.WithHTTPClient(httpClient))

	ctx := context.Background()

	iter := paging.Paginate(ctx, client.V2().KubernetesDev().Nodepools().ListFullPageFetcher(kubernetes.NodepoolListParams{
		DNSOverrideIPv4: new(false),
	}))

	allNodepools, err := paging.CollectAll(iter)
	if err != nil {
		panic(err)
	}

	outjson, _ := json.MarshalIndent(allNodepools, "", "  ")

	fmt.Println(string(outjson))
}
