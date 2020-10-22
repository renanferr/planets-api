package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
)

type MockHttpClient struct {
	Response *http.Response
	Err      error
}

func NewMockHttpClient(status int, body io.ReadCloser, err error) *MockHttpClient {
	return &MockHttpClient{
		&http.Response{
			StatusCode: status,
			Body:       body,
		},
		err,
	}
}

func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	return m.Response, m.Err
}

func TestGetPlanetByName(t *testing.T) {
	planet := Planet{
		Name:           "tatooine",
		RotationPeriod: "",
		OrbitalPeriod:  "",
		Diameter:       "",
		Climate:        "",
		Gravity:        "",
		Terrain:        "",
		SurfaceWater:   "",
		Population:     "",
		ResidentURLs:   []string{},
		FilmURLs:       []string{"", "", "", "", ""},
		Created:        "",
		Edited:         "",
		URL:            "",
	}
	result := &searchResult{
		Planets: []Planet{planet},
	}
	b, err := json.Marshal(result)
	r := ioutil.NopCloser(bytes.NewReader(b))

	c, err := NewClient("http://mock.test")
	c.httpClient = NewMockHttpClient(http.StatusOK, r, nil)
	if err != nil {
		t.Fatal(err)
	}

	p, err := c.GetPlanetByName(context.Background(), "tatooine")
	if err != nil {
		t.Fatal(err)
	}

	if p.Name != planet.Name {
		t.Errorf("planets names do not match. got: %s want %s", p.Name, planet.Name)
	}

}

func TestGetPlanetAppearances(t *testing.T) {
	planet := Planet{
		Name:           "tatooine",
		RotationPeriod: "",
		OrbitalPeriod:  "",
		Diameter:       "",
		Climate:        "",
		Gravity:        "",
		Terrain:        "",
		SurfaceWater:   "",
		Population:     "",
		ResidentURLs:   []string{},
		FilmURLs:       []string{"", "", "", "", ""},
		Created:        "",
		Edited:         "",
		URL:            "",
	}
	result := &searchResult{
		Planets: []Planet{planet},
	}
	b, err := json.Marshal(result)
	r := ioutil.NopCloser(bytes.NewReader(b))

	c, err := NewClient("http://mock.test")
	c.httpClient = NewMockHttpClient(http.StatusOK, r, nil)
	if err != nil {
		t.Fatal(err)
	}

	appearances, err := c.GetPlanetAppearances(context.Background(), "tatooine")
	if err != nil {
		t.Fatal(err)
	}

	if appearances != len(planet.FilmURLs) {
		t.Errorf("appearances do not match. got: %d want: %d", appearances, len(planet.FilmURLs))
	}
}
