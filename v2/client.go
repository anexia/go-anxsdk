package v2

import (
	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/v2/kubernetes"
	"github.com/anexia/go-anxsdk/v2/lbaas"
)

// API is the interface implemented by [Client].
type API interface {
	// KubernetesByEnv returns a clusters client based on the provided env.
	KubernetesByEnv(env kubernetes.Env) kubernetes.API
	// Kubernetes returns a clusters client.
	Kubernetes() kubernetes.API
	// KubernetesStage returns a stage clusters client.
	KubernetesStage() kubernetes.API
	// KubernetesDev returns a dev clusters client.
	KubernetesDev() kubernetes.API
	// LBaaS returns an entrypoint for getting lbaas resource clients.
	LBaaS() lbaas.API
}

// Client is an anexia v2 api client.
type Client struct {
	transport *internal.Transport
}

// Check that Client implements the API interface.
var _ API = (*Client)(nil)

// NewClient creates a new v2 api client.
func NewClient(transport *internal.Transport) *Client {
	return &Client{
		transport: transport,
	}
}

// KubernetesByEnv returns a clusters client based on the provided env.
func (c *Client) KubernetesByEnv(env kubernetes.Env) kubernetes.API {
	return kubernetes.NewClient(c.transport, env)
}

// Kubernetes returns a clusters client.
func (c *Client) Kubernetes() kubernetes.API {
	return kubernetes.NewClient(c.transport, kubernetes.EnvProduction)
}

// KubernetesStage returns a stage clusters client.
func (c *Client) KubernetesStage() kubernetes.API {
	return kubernetes.NewClient(c.transport, kubernetes.EnvStaging)
}

// KubernetesDev returns a dev clusters client.
func (c *Client) KubernetesDev() kubernetes.API {
	return kubernetes.NewClient(c.transport, kubernetes.EnvDevelopment)
}

// LBaaS returns an entrypoint for getting lbaas resource clients.
func (c *Client) LBaaS() lbaas.API {
	return lbaas.NewClient(c.transport)
}
