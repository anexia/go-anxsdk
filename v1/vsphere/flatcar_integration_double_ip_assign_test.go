package vsphere_test

import (
	"context"
	"fmt"
	"sync"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/v1/vsphere"
)

var _ = Describe("provisioning two flatcar vms concurrently with the same ip is rejected", Ordered, func() {
	// provisionedVM identifies a provisioned vm and carries any provisioning error.
	type provisionedVM struct {
		name string
		id   string
		err  error
	}

	var (
		ctx                context.Context
		cancel             context.CancelFunc
		provisioningClient vsphere.ProvisioningAPI
		infoClient         vsphere.InfoAPI
		vms                []provisionedVM
	)

	const (
		flatcarVMSSHKey      = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOJtesttesttesttesttesttesttesttest test@go-anxsdk"
		requestedIP          = "10.244.35.211"
		rejectDoubleIPAssign = false // set to true once the backend rejects the duplicate ip assignment
	)

	BeforeAll(func() {
		// The backend currently accepts BOTH concurrent provisioning requests and
		// assigns the same ip to both vms instead of rejecting the duplicate.
		// This test expects exactly one of the requests to error, so it is skipped
		// until the engine returns an error for the second allocation.
		if !rejectDoubleIPAssign {
			Skip("backend does not yet reject a double ip assignment, set rejectDoubleIPAssign=true once fixed")
		}

		ctx, cancel = context.WithTimeout(context.Background(), 15*time.Minute)
		DeferCleanup(cancel)

		provisioningClient = vsphereClient.Provisioning()
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

		// Provision both vms concurrently, both claiming the same ip. Exactly one
		// of them is expected to error out.
		vmNames := []string{
			"go-anxsdk-flatcar-e2e-first",
			"go-anxsdk-flatcar-e2e-second",
		}

		// Each goroutine writes to its own slice slot (distinct memory, no data
		// race); results are asserted afterwards on the main goroutine since
		// Ginkgo must not be failed from worker goroutines.
		results := make([]provisionedVM, len(vmNames))

		var wg sync.WaitGroup
		for i, vmName := range vmNames {
			wg.Add(1)
			go func(i int, vmName string) {
				defer wg.Done()

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
							IPs:            []string{requestedIP},
						},
					},
				}

				provisionResponse, err := provisioningClient.ProvisionTemplate(ctx, location.ID, template.ID, request)
				if err != nil {
					results[i] = provisionedVM{name: vmName, err: fmt.Errorf("provisioning vm: %w", err)}
					return
				}

				vmID, err := provisioningClient.AwaitCompletion(ctx, provisionResponse.TaskIdentifier)
				if err != nil {
					results[i] = provisionedVM{name: vmName, err: fmt.Errorf("awaiting completion: %w", err)}
					return
				}

				results[i] = provisionedVM{name: vmName, id: vmID}
			}(i, vmName)
		}
		wg.Wait()

		vms = results
		for _, vm := range vms {
			if vm.err != nil {
				GinkgoWriter.Printf("vm %q was rejected during provisioning: %v\n", vm.name, vm.err)
			} else {
				GinkgoWriter.Printf("provisioned flatcar vm %q with identifier %q\n", vm.name, vm.id)
			}
		}
	})

	AfterAll(func() {
		for _, vm := range vms {
			if vm.id == "" {
				continue // this vm was rejected during provisioning
			}
			_, err := provisioningClient.Deprovision(ctx, vm.id, false)
			Expect(err).NotTo(HaveOccurred())
		}
	})

	It("reports exactly one error for the duplicate ip request", func() {
		var (
			rejected  []provisionedVM
			survivors []provisionedVM
		)
		for _, vm := range vms {
			if vm.err != nil {
				rejected = append(rejected, vm)
			} else {
				survivors = append(survivors, vm)
			}
		}

		Expect(rejected).To(HaveLen(1), "expected exactly one of the concurrently provisioned vms to be rejected for ip %q, got rejected=%d survived=%d: %+v", requestedIP, len(rejected), len(survivors), vms)
		Expect(survivors).To(HaveLen(1), "expected exactly one of the concurrently provisioned vms to survive for ip %q", requestedIP)
	})

	It("the surviving vm is running and holds the requested ip", func() {
		var survivor *provisionedVM
		for i := range vms {
			if vms[i].err == nil {
				survivor = &vms[i]
			}
		}
		Expect(survivor).NotTo(BeNil())
		Expect(survivor.id).NotTo(BeEmpty())

		info, err := infoClient.Get(ctx, survivor.id)
		Expect(err).NotTo(HaveOccurred())
		Expect(info.Identifier).To(Equal(survivor.id))
		Expect(info.Status).To(Equal(vsphere.PowerStatePoweredOn))
		Expect(info.Network).NotTo(BeEmpty())
		Expect(info.DiskInfo).NotTo(BeEmpty())

		var vmIPs []string
		for _, n := range info.Network {
			vmIPs = append(vmIPs, n.IPsv4...)
		}
		GinkgoWriter.Printf("info: vm %q (id %q) is running with ipv4 ips %v\n", survivor.name, survivor.id, vmIPs)
		Expect(vmIPs).To(ContainElement(requestedIP), "surviving vm %q does not hold the requested ip %s: %v", survivor.name, requestedIP, vmIPs)
	})
})
