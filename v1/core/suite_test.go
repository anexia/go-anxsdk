package core_test

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/anexia/go-anxsdk/v1/core"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/anexia/go-anxsdk"
	"github.com/anexia/go-anxsdk/utils"
)

func TestCore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Core tests")
}

var (
	coreClient *core.Client
)

var _ = BeforeSuite(func() {
	apiKey := os.Getenv("ANEXIA_TOKEN")
	if apiKey == "" {
		Skip("ANEXIA_TOKEN is not set, skipping core integration suite")
	}

	httpClient := http.DefaultClient
	httpClient.Transport = utils.NewRateLimitRoundTripper(httpClient.Transport, 10*time.Second, 10)

	cl := anxsdk.NewClient(anxsdk.WithAPIKey(apiKey), anxsdk.WithHTTPClient(httpClient))

	coreClient = cl.V1().Core()
})
