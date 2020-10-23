package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/renanferr/swapi-golang-rest-api/pkg/adding"
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

func isPlanetInSlice(p *Planet, planets *[]Planet) bool {
	for _, planet := range *planets {
		if p.Name == planet.Name {
			return true
		}
	}
	return false
}
func TestGetPlanetByName(t *testing.T) {
	type TestCase struct {
		Name        string
		Planets     []Planet
		Query       string
		ResponseErr error
		ExpectedErr error
	}

	tt := []TestCase{
		{
			"get planet successfully",
			[]Planet{
				{"tatooine", []string{"", "", "", "", ""}},
			},
			"tatooine",
			nil,
			nil,
		},
		{
			"planet not found",
			[]Planet{},
			"tatooine",
			nil,
			adding.ErrPlanetNotFound,
		},
	}

	for _, tc := range tt {

		result := &searchResult{
			Planets: tc.Planets,
		}
		b, err := json.Marshal(result)
		r := ioutil.NopCloser(bytes.NewReader(b))

		c, err := NewClient("http://mock.test")
		if err != nil {
			t.Fatalf("<%s> unexpected err creating client: %s", tc.Name, err.Error())
		}
		mockClient := NewMockHttpClient(http.StatusOK, r, tc.ResponseErr)
		c.httpClient = mockClient

		p, err := c.GetPlanetByName(context.Background(), tc.Query)
		if err != nil {
			if !errors.Is(err, tc.ResponseErr) && !errors.Is(err, tc.ExpectedErr) {
				t.Fatalf("<%s> unexpected error: %s", tc.Name, err.Error())
			}
		} else {
			if !isPlanetInSlice(&p, &tc.Planets) {
				t.Errorf("<%s> planet with name \"%s\" is not in the expected response %s", tc.Name, p.Name, tc.Planets)
			}
		}

	}

}

func TestGetPlanetAppearances(t *testing.T) {

	type TestCase struct {
		Name          string
		Planets       []Planet
		ExpectedValue int
		ExpectedErr   error
	}

	tt := []TestCase{
		{
			"get tatooine's 5 appearances",
			[]Planet{
				{"tatooine", []string{"", "", "", "", ""}},
			},
			5,
			nil,
		},
		{
			"planet not found",
			[]Planet{},
			0,
			adding.ErrPlanetNotFound,
		},
	}

	for _, tc := range tt {

		result := &searchResult{
			Planets: tc.Planets,
		}
		b, err := json.Marshal(result)
		r := ioutil.NopCloser(bytes.NewReader(b))

		c, err := NewClient("http://mock.test")
		if err != nil {
			t.Fatalf("<%s> unexpected error creating client: %s", tc.Name, err.Error())
		}
		c.httpClient = NewMockHttpClient(http.StatusOK, r, nil)

		appearances, err := c.GetPlanetAppearances(context.Background(), "tatooine")
		if err != nil {
			if !errors.Is(err, tc.ExpectedErr) {
				t.Fatalf("<%s> unexpected error: %s", tc.Name, err.Error())
			}
		} else {
			if appearances != tc.ExpectedValue {
				t.Errorf("appearances do not match. got: %d want: %d", appearances, tc.ExpectedValue)
			}

		}
	}
}
