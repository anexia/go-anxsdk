package vsphere_test

import (
	"context"
	"fmt"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/v1/vsphere"
)

var _ = Describe("provisioning a flatcar vm and resolving it via search and info", Ordered, func() {
	var (
		ctx                context.Context
		cancel             context.CancelFunc
		provisioningClient vsphere.ProvisioningAPI
		searchClient       vsphere.SearchAPI
		infoClient         vsphere.InfoAPI
		vmID               string
		vmName             string
	)

	const (
		flatcarVMSSHKey = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOJtesttesttesttesttesttesttesttest test@go-anxsdk"
		searchTimeout   = 10 * time.Minute
		searchInterval  = 5 * time.Second
	)

	BeforeAll(func() {
		ctx, cancel = context.WithTimeout(context.Background(), 15*time.Minute)
		DeferCleanup(cancel)

		provisioningClient = vsphereClient.Provisioning()
		searchClient = vsphereClient.Search()
		infoClient = vsphereClient.Info()

		var location vsphere.Location
		locations, err := provisioningClient.ListLocations(ctx, paging.DefaultParams())
		Expect(err).NotTo(HaveOccurred())
		for _, loc := range locations.Data {
			if loc.Code == "ANX04" {
				location = loc
			}
		}
		Expect(location.ID).NotTo(BeEmpty())

		template, err := provisioningClient.FindNamedTemplate(ctx, location.ID, flatcarTemplateName, vsphere.LatestTemplateBuild)
		Expect(err).NotTo(HaveOccurred())

		vmName = fmt.Sprintf("go-anxsdk-flatcar-e2e-%d", time.Now().UnixNano())

		request := vsphere.ProvisioningRequest{
			Hostname:           vmName,
			CustomName:         &vmName,
			CPUs:               new(2),
			MemoryMB:           new(2048),
			DiskGB:             new(16),
			DiskType:           new(vsphere.DiskTypeENT3),
			CPUPerformanceType: new(vsphere.CPUPerformanceTypeStandard),
			SSH:                new(flatcarVMSSHKey),
			Network: []vsphere.ProvisioningRequestNetwork{
				{
					NicType:        vsphere.NicTypeVirtIO,
					VLan:           "70c0ac17cf234553949558b4fae7ee8f",
					BandwidthLimit: 1000,
					IPs:            []string{"10.244.35.222"},
				},
			},
		}

		provisionResponse, err := provisioningClient.ProvisionTemplate(ctx, location.ID, template.ID, request)
		Expect(err).NotTo(HaveOccurred())

		_ = provisionResponse

		vmID, err = provisioningClient.AwaitCompletion(ctx, provisionResponse.TaskIdentifier)
		Expect(err).NotTo(HaveOccurred())
		Expect(vmID).NotTo(BeEmpty())

		GinkgoWriter.Printf("provisioned flatcar vm %q with identifier %q\n", vmName, vmID)
	})

	AfterAll(func() {
		_, err := provisioningClient.Deprovision(ctx, vmID, false)
		Expect(err).NotTo(HaveOccurred())
	})

	It("finds the vm via repeated searchByName", func() {
		var (
			found         *vsphere.SearchResult
			attempts      int
			emptyAttempts int
			lastErr       error
		)

		Eventually(func() bool {
			attempts++

			results, err := paging.CollectAll(paging.Paginate(ctx, searchClient.ByNamePageFetcher(vsphere.SearchByNameParams{Name: vmName})))
			if err != nil {
				lastErr = err
				return false
			}

			if len(results) == 0 {
				emptyAttempts++
			}

			for i := range results {
				if strings.HasSuffix(results[i].Name, vmName) {
					found = &results[i]
					return true
				}
			}

			return false
		}).WithTimeout(searchTimeout).WithPolling(searchInterval).Should(BeTrue(), "search by name never returned the vm: %v", lastErr)

		GinkgoWriter.Printf("search by name: found after %d attempt(s), %d of which returned no results (vm id %q)\n", attempts, emptyAttempts, vmID)

		Expect(found).NotTo(BeNil())
		Expect(found.Name[strings.Index(found.Name, "-")+1:]).To(Equal(vmName)) // cut the digits prefix
		Expect(found.Identifier).To(Equal(vmID))
	})

	It("resolves the full details via the info endpoint", func() {
		info, err := infoClient.Get(ctx, vmID)
		Expect(err).NotTo(HaveOccurred())
		Expect(info).NotTo(BeNil())

		Expect(info.Identifier).To(Equal(vmID))
		Expect(info.Name).To(ContainSubstring(vmName))
		Expect(info.Status).To(Equal(vsphere.PowerStatePoweredOn))
		Expect(info.LocationIdentifier).NotTo(BeEmpty())
		Expect(info.Network).NotTo(BeEmpty())
		Expect(info.DiskInfo).NotTo(BeEmpty())
	})
})
