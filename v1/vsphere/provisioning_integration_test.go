package vsphere_test

import (
	"context"
	"time"

	"github.com/anexia/go-anxsdk/v1/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/v1/vsphere"
)

const flatcarTemplateName = "Flatcar Linux Stable UEFI"

var _ = Describe("ProvisioningClient", func() {
	var (
		ctx                context.Context
		cancel             context.CancelFunc
		anx04Location      vsphere.Location
		template           vsphere.TemplateResponse
		provisioningClient *vsphere.ProvisioningClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(context.Background(), time.Minute)
		DeferCleanup(cancel)

		provisioningClient = vsphereClient.Provisioning()
	})

	Describe("ListLocations", func() {
		It("lists locations and finds ANX04", func() {
			locations, err := provisioningClient.ListLocations(ctx, paging.DefaultParams())
			Expect(err).NotTo(HaveOccurred())
			Expect(locations.Data).NotTo(BeEmpty())

			for _, loc := range locations.Data {
				if loc.Code == "ANX04" {
					anx04Location = loc
				}
			}
			Expect(anx04Location.ID).NotTo(BeEmpty())
		})
	})

	Describe("ListTemplates", func() {
		It("lists templates for the location and finds the debian template", func() {
			templates, err := provisioningClient.ListTemplates(ctx, anx04Location.ID, vsphere.TemplateTypeTemplates)
			Expect(err).NotTo(HaveOccurred())
			Expect(templates).NotTo(BeEmpty())

			for _, tmpl := range templates {
				if tmpl.Name == "Debian 13 UEFI" {
					template = tmpl
				}
			}
			Expect(template.ID).NotTo(BeEmpty())
		})
	})

	Describe("FindNamedTemplate", func() {
		It("finds the latest build of the flatcar template", func() {
			tmpl, err := provisioningClient.FindNamedTemplate(ctx, anx04Location.ID, flatcarTemplateName, vsphere.LatestTemplateBuild)
			Expect(err).NotTo(HaveOccurred())
			Expect(tmpl).NotTo(BeNil())
			Expect(tmpl.Name).To(Equal(flatcarTemplateName))
		})

		It("finds a specified build of the flatcar template", func() {
			// this specified build may be removed in the future, change it if test fails
			tmpl, err := provisioningClient.FindNamedTemplate(ctx, anx04Location.ID, flatcarTemplateName, "b25")
			Expect(err).NotTo(HaveOccurred())
			Expect(tmpl).NotTo(BeNil())
			Expect(tmpl.Name).To(Equal(flatcarTemplateName))
		})

		It("cannot find non-existing templates", func() {
			tmpl, err := provisioningClient.FindNamedTemplate(ctx, anx04Location.ID, "i do not exist", vsphere.LatestTemplateBuild)
			Expect(err).To(HaveOccurred())
			Expect(common.IsNotFoundError(err)).To(BeTrue())
			Expect(tmpl).To(BeNil())
		})

		It("cannot find templates on wrong locations", func() {
			tmpl, err := provisioningClient.FindNamedTemplate(ctx, "wrong location id", "i do not exist", vsphere.LatestTemplateBuild)
			Expect(err).To(HaveOccurred())
			Expect(common.IsNotFoundError(err)).To(BeTrue())
			Expect(tmpl).To(BeNil())
		})
	})

	Describe("GetCPUArchitectures", func() {
		It("returns the available cpu architectures", func() {
			architectures, err := provisioningClient.GetCPUArchitectures(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(architectures).NotTo(BeEmpty())
		})
	})

	Describe("GetCPUPerformanceTypes", func() {
		It("returns the available cpu performance types", func() {
			cpuPerfTypes, err := provisioningClient.GetCPUPerformanceTypes(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(cpuPerfTypes).NotTo(BeEmpty())
		})
	})

	Describe("GetDiskTypes", func() {
		It("returns the available disk types for the location", func() {
			diskTypes, err := provisioningClient.GetDiskTypes(ctx, anx04Location.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(diskTypes).NotTo(BeEmpty())
		})
	})

	Describe("ListAvailabilityZones", func() {
		It("lists the availability zones for the location", func() {
			availabilityZones, err := provisioningClient.ListAvailabilityZones(ctx, anx04Location.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(availabilityZones).NotTo(BeEmpty())
		})
	})

	Describe("GetNicTypes", func() {
		It("returns the available nic types", func() {
			nicTypes, err := provisioningClient.GetNicTypes(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(nicTypes).NotTo(BeEmpty())
		})
	})
})
