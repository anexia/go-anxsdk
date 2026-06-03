package anxsdk

import (
	"net/http"

	"github.com/anexia/go-anxsdk/internal"
)

// Config is a collection of all configuration options for a new client.
// All fields are intentionally unexported and should be configured via the options pattern (see ClientOption).
type Config struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// newConfig creates a new Config struct with sensible defaults (url, http client, ...).
// Config settings can be overridden by passing in a corresponding config option.
// A user should always configure an api key via WithAPIKey.
func newConfig(opts ...Option) Config {
	cfg := Config{
		baseURL: "https://engine.anexia-it.com",
		client:  http.DefaultClient,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	return cfg
}

// createTransport creates an internal.Transport helper from the config.
func (c Config) createTransport() *internal.Transport {
	return internal.NewTransport(
		c.baseURL,
		c.apiKey,
		c.client,
	)
}

// Option is a named type to improve code clarity and intent for configuration options.
type Option func(*Config)

// WithAPIKey provides a config option to use the given api key.
func WithAPIKey(apiKey string) Option {
	return func(config *Config) {
		config.apiKey = apiKey
	}
}

// WithBaseURL provides a config option to use the given url as the base url.
func WithBaseURL(url string) Option {
	return func(c *Config) {
		c.baseURL = url
	}
}

// WithHTTPClient provides a config option to use the given http client.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Config) {
		c.client = client
	}
}
