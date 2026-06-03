package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/anexia/go-anxsdk/paging"
	"github.com/google/go-querystring/query"
)

// TransportError is an internal http error.
type TransportError struct {
	StatusCode int
	Status     string
	Body       string
}

// Error implements the error interface for TransportError.
func (a *TransportError) Error() string {
	return fmt.Sprintf("transport error: StatusCode=%d, Status=%s, Body=%s", a.StatusCode, a.Status, a.Body)
}

// Transport is an internal http transport helper.
// It is responsible for all http protocol layer logic and provides the reusable functions used in the client implementations.
type Transport struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

func (t *Transport) buildRequestURL(endpoint string, pageParams *paging.Params, filterParams any) (string, error) {
	base, err := url.Parse(t.baseURL)
	if err != nil {
		return "", fmt.Errorf("parsing base url: %w", err)
	}

	rel, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("parsing endpoint url: %w", err)
	}

	base = base.ResolveReference(rel)

	q, err := query.Values(filterParams)
	if err != nil {
		return "", fmt.Errorf("building query params: %w", err)
	}

	if pageParams != nil {
		err := pageParams.Validate()
		if err != nil {
			return "", fmt.Errorf("validating page params: %w", err)
		}

		q.Set("page", strconv.Itoa(pageParams.Page))
		q.Set("limit", strconv.Itoa(pageParams.Limit))
	}

	base.RawQuery = q.Encode()

	fullURL := base.String()
	return fullURL, nil
}

func (t *Transport) newRequest(ctx context.Context, method, endpoint string, body io.Reader, pageParams *paging.Params, filterParams any) (*http.Request, error) {
	fullURL, err := t.buildRequestURL(endpoint, pageParams, filterParams)
	if err != nil {
		return nil, fmt.Errorf("building request url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	if t.apiKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Token %s", t.apiKey))
	}

	if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (t *Transport) do(req *http.Request, response any) error {
	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("executing http request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode >= http.StatusBadRequest {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("reading response body: %w", err)
		}

		return &TransportError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Body:       string(body),
		}
	}

	if response != nil {
		err = json.NewDecoder(resp.Body).Decode(response)
		if err != nil {
			return fmt.Errorf("decoding body: %w", err)
		}
	}

	return nil
}

func (t *Transport) doRequest(ctx context.Context, method string, endpoint string, request any, response any, pageParams *paging.Params, params any) error {
	var body io.Reader

	if request != nil {
		var reqBody bytes.Buffer

		err := json.NewEncoder(&reqBody).Encode(request)
		if err != nil {
			return fmt.Errorf("marshalling request: %w", err)
		}

		body = &reqBody
	}

	req, err := t.newRequest(ctx, method, endpoint, body, pageParams, params)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	err = t.do(req, response)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}

	return nil
}

// Get executes a get request against the anexia api.
// endpoint is the relative url (to the configured base url) of the http api endpoint.
// response must be a pointer to a response struct in which the response body will be JSON deserialized if not nil.
// pageParams is the paging parameters
// params will be used to extract the query parameters for the request via `query` attributes.
func (t *Transport) Get(ctx context.Context, endpoint string, response any, pageParams paging.Params, params any) error {
	return t.doRequest(ctx, http.MethodGet, endpoint, nil, response, &pageParams, params)
}

// GetSingle executes a get request against the anexia api.
// endpoint is the relative url (to the configured base url) of the http api endpoint.
// response must be a pointer to a response struct in which the response body will be JSON deserialized if not nil.
func (t *Transport) GetSingle(ctx context.Context, endpoint string, response any) error {
	return t.doRequest(ctx, http.MethodGet, endpoint, nil, response, nil, nil)
}

// Delete executes a delete request against the anexia api.
// endpoint is the relative url (to the configured base url) of the http api endpoint.
func (t *Transport) Delete(ctx context.Context, endpoint string) error {
	return t.doRequest(ctx, http.MethodDelete, endpoint, nil, nil, nil, nil)
}

// Post executes a post request against the anexia api.
// endpoint is the relative url (to the configured base url) of the http api endpoint.
// request if not nil will be JSON serialized and sent as the request body.
// response must be a pointer to a response struct in which the response body will be JSON deserialized if not nil.
func (t *Transport) Post(ctx context.Context, endpoint string, request any, response any) error {
	return t.doRequest(ctx, http.MethodPost, endpoint, request, response, nil, nil)
}

// Put executes a put request against the anexia api.
// endpoint is the relative url (to the configured base url) of the http api endpoint.
// request if not nil will be JSON serialized and sent as the request body.
// response must be a pointer to a response struct in which the response body will be JSON deserialized if not nil.
func (t *Transport) Put(ctx context.Context, endpoint string, request any, response any) error {
	return t.doRequest(ctx, http.MethodPut, endpoint, request, response, nil, nil)
}

// NewTransport creates a new internal transport helper with the provided base url, api key and http client.
func NewTransport(baseURL string, apiKey string, client *http.Client) *Transport {
	if client == nil {
		client = http.DefaultClient
	}

	return &Transport{
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  client,
	}
}
