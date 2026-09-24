package vsphere_test

import (
	"net/http"
	"os"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/anexia/go-anxsdk"
	"github.com/anexia/go-anxsdk/utils"
	"github.com/anexia/go-anxsdk/v1/vsphere"
)

func TestVSphere(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VSphere tests")
}

var (
	vsphereClient *vsphere.Client
)

var _ = BeforeSuite(func() {
	apiKey := os.Getenv("ANEXIA_TOKEN")
	if apiKey == "" {
		Skip("ANEXIA_TOKEN is not set, skipping vsphere integration suite")
	}

	httpClient := http.DefaultClient
	httpClient.Transport = utils.NewRateLimitRoundTripper(httpClient.Transport, 10*time.Second, 10)

	cl := anxsdk.NewClient(anxsdk.WithAPIKey(apiKey), anxsdk.WithHTTPClient(httpClient))

	vsphereClient = cl.V1().VSphere()
})
