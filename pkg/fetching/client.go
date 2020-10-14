package fetching

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	pathlib "path"
)

// A Client communicates with SWAPI
type Client struct {
	// baseURL is the base url for SWAPI
	baseURL *url.URL

	// HTTP client used to communicate with the SWAPI
	httpClient *http.Client
}

type query = map[string]string

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

func (c *Client) newRequest(ctx context.Context, path string, query query) (*http.Request, error) {
	u, err := url.Parse(pathlib.Join(c.baseURL.String(), path))
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("format", "json")
	for k, v := range query {
		q.Set(k, v)
	}

	u.RawQuery = q.Encode()

	return http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	req.Close = true

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}

	if err != nil {
		return nil, fmt.Errorf("error reading response from %s %s: %s", req.Method, req.URL.RequestURI(), err)
	}

	return resp, nil
}

func (c *Client) GetPlanetByName(ctx context.Context, planetName string) (Planet, error) {
	req, err := c.newRequest(ctx, "/planets", query{"name": planetName})
	if err != nil {
		return Planet{}, err
	}

	var planet Planet

	if _, err = c.do(req, &planet); err != nil {
		return Planet{}, err
	}

	return planet, nil
}

func (c *Client) GetPlanetAppearances(ctx context.Context, planetName string) (int, error) {
	planet, err := c.GetPlanetByName(ctx, planetName)
	if err != nil {
		return 0, err
	}

	return len(planet.FilmURLs), nil
}
