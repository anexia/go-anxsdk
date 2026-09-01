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
	"github.com/anexia/go-anxsdk/v1/vsphere"
)

func runVSphere() {
	httpClient := http.DefaultClient
	httpClient.Transport = utils.NewRateLimitRoundTripper(httpClient.Transport, time.Second*10, 10)

	client := anxsdk.NewClient(anxsdk.WithAPIKey(os.Getenv("ANEXIA_TOKEN")), anxsdk.WithHTTPClient(httpClient))

	ctx := context.Background()

	locations, err := client.V1().VSphere().Provisioning().ListLocations(ctx, paging.DefaultParams())
	if err != nil {
		panic(err)
	}

	var anx04Location vsphere.Location
	for _, loc := range locations.Data {
		if loc.Code == "ANX04" {
			anx04Location = loc
		}
	}

	outjson, _ := json.MarshalIndent(anx04Location, "", "  ")
	fmt.Println("anx04 location: ", string(outjson))

	templates, err := client.V1().VSphere().Provisioning().ListTemplates(ctx, anx04Location.ID, vsphere.TemplateTypeTemplates)
	if err != nil {
		panic(err)
	}

	var flatcarTemplate vsphere.TemplateResponse
	for _, template := range templates {
		if template.Name == "Flatcar Linux Stable UEFI" {
			flatcarTemplate = template
		}
	}

	outjson, _ = json.MarshalIndent(flatcarTemplate, "", "  ")
	fmt.Println("flatcar template: ", string(outjson))

	architectures, err := client.V1().VSphere().Provisioning().GetCPUArchitectures(ctx)
	if err != nil {
		panic(err)
	}

	outjson, _ = json.MarshalIndent(architectures, "", "  ")
	fmt.Println("CPU architectures: ", string(outjson))

	cpuPerfTypes, err := client.V1().VSphere().Provisioning().GetCPUPerformanceTypes(ctx)
	if err != nil {
		panic(err)
	}

	outjson, _ = json.MarshalIndent(cpuPerfTypes, "", "  ")
	fmt.Println("CPU performance types: ", string(outjson))

	diskTypes, err := client.V1().VSphere().Provisioning().GetDiskTypes(ctx, anx04Location.ID)
	if err != nil {
		panic(err)
	}

	outjson, _ = json.MarshalIndent(diskTypes, "", "  ")
	fmt.Println("anx04 disk types: ", string(outjson))

	availabilityZones, err := client.V1().VSphere().Provisioning().ListAvailabilityZones(ctx, anx04Location.ID)
	if err != nil {
		panic(err)
	}

	outjson, _ = json.MarshalIndent(availabilityZones, "", "  ")
	fmt.Println("anx04 availability zones: ", string(outjson))

	nicTypes, err := client.V1().VSphere().Provisioning().GetNicTypes(ctx)
	if err != nil {
		panic(err)
	}

	outjson, _ = json.MarshalIndent(nicTypes, "", "  ")
	fmt.Println("nic types: ", string(outjson))
}
