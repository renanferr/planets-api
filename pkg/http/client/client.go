package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	pathlib "path"
	"strings"
)

type searchResult struct {
	Planets []Planet `json:"results"`
}

// HTTPClient interface
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// A Client communicates with SWAPI
type Client struct {
	// baseURL is the base url for SWAPI
	baseURL *url.URL

	// HTTP client used to communicate with the SWAPI
	httpClient HTTPClient
}

func NewClient(baseURL string) (*Client, error) {
	parsed, err := url.Parse(baseURL)

	if err != nil {
		return &Client{}, err
	}
	c := &Client{
		baseURL:    parsed,
		httpClient: http.DefaultClient,
	}

	return c, nil
}

func (c *Client) newRequest(ctx context.Context, path string, query url.Values) (*http.Request, error) {
	url := c.baseURL
	url.Path = pathlib.Join(url.Path, path)

	q := url.Query()
	q.Set("format", "json")
	for k, v := range query {
		q.Set(k, strings.Join(v, ","))
	}

	url.RawQuery = q.Encode()

	return http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
}

func (c *Client) do(req *http.Request, target interface{}) (*http.Response, error) {
	req.Close = true
	log.Printf("[%s] %s", req.Method, req.URL)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(target)

	if err != nil {
		return nil, fmt.Errorf("error reading response from %s %s: %s", req.Method, req.URL.RequestURI(), err)
	}

	return resp, nil
}
