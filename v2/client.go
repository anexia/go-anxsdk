package v2

import (
	"github.com/anexia/go-anxsdk/internal"
	"github.com/anexia/go-anxsdk/v2/kubernetes"
	"github.com/anexia/go-anxsdk/v2/lbaas"
)

// Client is an anexia v2 api client.
type Client struct {
	transport *internal.Transport
}

// NewClient creates a new v2 api client.
func NewClient(transport *internal.Transport) *Client {
	return &Client{
		transport: transport,
	}
}

// KubernetesByEnv returns a clusters client based on the provided env.
func (c *Client) KubernetesByEnv(env kubernetes.Env) *kubernetes.Client {
	return kubernetes.NewClient(c.transport, env)
}

// Kubernetes returns a clusters client.
func (c *Client) Kubernetes() *kubernetes.Client {
	return kubernetes.NewClient(c.transport, kubernetes.EnvProduction)
}

// KubernetesStage returns a stage clusters client.
func (c *Client) KubernetesStage() *kubernetes.Client {
	return kubernetes.NewClient(c.transport, kubernetes.EnvStaging)
}

// KubernetesDev returns a dev clusters client.
func (c *Client) KubernetesDev() *kubernetes.Client {
	return kubernetes.NewClient(c.transport, kubernetes.EnvDevelopment)
}

// LBaaS returns an entrypoint for getting lbaas resource clients.
func (c *Client) LBaaS() *lbaas.Client {
	return lbaas.NewClient(c.transport)
}
