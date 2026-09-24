package vsphere_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/v1/vsphere"
)

var _ = Describe("SearchClient", func() {
	var (
		ctx          context.Context
		cancel       context.CancelFunc
		searchClient *vsphere.SearchClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(context.Background(), time.Minute)
		DeferCleanup(cancel)

		searchClient = vsphereClient.Search()
	})

	Describe("ByName", func() {
		It("finds VMs matching a name via the page fetcher", func() {
			fetcher := searchClient.ByNamePageFetcher(vsphere.SearchByNameParams{Name: "%servicevm%"})

			results, err := paging.CollectAll(paging.Paginate(ctx, fetcher))
			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())

			for _, r := range results {
				Expect(r.Name).To(ContainSubstring("servicevm"))
			}
		})
	})

	Describe("ByTags", func() {
		It("finds VMs matching a tag via the page fetcher", func() {
			fetcher := searchClient.ByTagsPageFetcher(vsphere.SearchByTagsParams{Tags: []string{"k8s"}})

			results, err := paging.CollectAll(paging.Paginate(ctx, fetcher))
			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())

			for _, r := range results {
				Expect(r.Tags).NotTo(BeEmpty())
				Expect(r.Tags).To(ContainElement("k8s"))
			}
		})
	})
})
