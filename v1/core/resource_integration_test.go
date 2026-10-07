package core_test

import (
	"context"
	"net/http"
	"time"

	"github.com/anexia/go-anxsdk/v1/common"
	"github.com/anexia/go-anxsdk/v1/core"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/anexia/go-anxsdk/paging"
)

var _ = Describe("CoreClient", Ordered, func() {
	var (
		ctx        context.Context
		cancel     context.CancelFunc
		resourceID string
	)

	const (
		tagName    = "k8s-svm"
		newTestTag = "go-anx-integration-test-tag"
	)

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(context.Background(), time.Minute)
		DeferCleanup(cancel)
	})

	Describe("ListResource", func() {
		It("lists resources", func() {
			pageResp, err := coreClient.Resources().List(ctx, paging.DefaultParams(), core.ResourceListParams{
				Query:   new("e2e%dev"),
				TagName: new(tagName),
			})

			Expect(err).NotTo(HaveOccurred())

			for _, r := range pageResp.Data {
				Expect(common.IsEngineIdentifier(r.Identifier)).To(BeTrue())
				Expect(r.Name).ToNot(BeEmpty())
			}
			resourceID = pageResp.Data[0].Identifier
		})
	})

	Describe("GetResource", func() {
		It("gets resource", func() {
			resp, err := coreClient.Resources().Get(ctx, resourceID)

			Expect(err).NotTo(HaveOccurred())
			Expect(resp).NotTo(BeNil())

			Expect(resp.Identifier).To(Equal(resourceID))
			Expect(resp.ServiceName).To(Equal("Dynamic Compute"))
		})
	})

	Describe("Tags", func() {
		It("gets tags", func() {
			tagsResp, err := coreClient.Resources().GetTags(ctx, resourceID)

			Expect(err).NotTo(HaveOccurred())
			Expect(tagsResp).NotTo(BeEmpty())
		})

		It("tries to assign already existing tag", func() {
			err := coreClient.Resources().AssignTag(ctx, resourceID, tagName)

			Expect(err).To(HaveOccurred())
			Expect(common.IsErrorWithStatusCode(err, http.StatusUnprocessableEntity)).To(BeTrue())
		})

		It("assigns a new test tag", func() {
			err := coreClient.Resources().AssignTag(ctx, resourceID, newTestTag)

			Expect(err).ToNot(HaveOccurred())
		})
	})

	AfterAll(func() {
		err := coreClient.Resources().RemovesTag(ctx, resourceID, newTestTag)

		Expect(err).ToNot(HaveOccurred())
	})
})
