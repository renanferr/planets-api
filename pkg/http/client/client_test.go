package client

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"testing"
)

func TestClient(t *testing.T) {
	baseURL := "http://mock.test"
	c, err := NewClient(baseURL)
	if err != nil {
		t.Fatal(err)
	}

	if c.baseURL.String() != baseURL {
		t.Errorf("%s is not equal to %s", baseURL, c.baseURL)
	}
}

func TestNewRequest(t *testing.T) {
	baseURL := "http://mock.test"
	path := "/test"

	q := url.Values{}
	q.Add("foo", "bar")

	c, err := NewClient(baseURL)
	if err != nil {
		t.Fatal(err)
	}

	req, err := c.newRequest(context.Background(), path, q)

	expectedQuery := url.Values{}
	expectedQuery.Add("format", "json")
	for k, v := range q {
		expectedQuery.Add(k, strings.Join(v, ","))
	}
	expectedURL := fmt.Sprintf("%s%s?%s", baseURL, path, expectedQuery.Encode())
	if req.URL.String() != expectedURL {
		t.Errorf("URLs do not match. got: %s want: %s", req.URL.String(), expectedURL)
	}
}
