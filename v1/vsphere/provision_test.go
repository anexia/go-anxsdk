package vsphere_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/anexia/go-anxsdk/v1/vsphere"
)

var _ = Describe("TemplateResponse", func() {
	Describe("BuildNumber", func() {
		It("parses a valid build number", func() {
			tmpl := vsphere.TemplateResponse{Build: "b123"}

			buildNumber, err := tmpl.BuildNumber()
			Expect(err).NotTo(HaveOccurred())
			Expect(buildNumber).To(Equal(123))
		})

		It("errors when the build is empty", func() {
			tmpl := vsphere.TemplateResponse{Build: ""}

			_, err := tmpl.BuildNumber()
			Expect(err).To(HaveOccurred())
		})

		It("errors when the build does not start with 'b'", func() {
			tmpl := vsphere.TemplateResponse{Build: "123"}

			_, err := tmpl.BuildNumber()
			Expect(err).To(HaveOccurred())
		})

		It("errors when the build number is not numeric", func() {
			tmpl := vsphere.TemplateResponse{Build: "bxyz"}

			_, err := tmpl.BuildNumber()
			Expect(err).To(HaveOccurred())
		})
	})
})
