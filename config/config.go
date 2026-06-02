package config

import (
	"net/http"

	"code.anexia.com/se/ks/go-anxsdk/internal"
)

// Config is a collection of all configuration options for a new client.
// All fields are intentionally unexported and should be configured via the options pattern (see ClientOption).
type Config struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// NewConfig creates a new Config struct with sensible defaults (url, http client, ...).
// Config settings can be overridden by passing in a corresponding config option.
// A user should always configure an api key via WithAPIKey.
func NewConfig(opts ...Option) Config {
	cfg := Config{
		baseURL: "https://engine.anexia-it.com",
		client:  http.DefaultClient,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	return cfg
}

// CreateTransport creates an internal transport helper from the config.
func (c Config) CreateTransport() *internal.Transport {
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

// WithBaseURL ovides a config option to use the given url as the base url.
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
