package client

import (
	"context"
	"errors"
	"net/url"
)

// A Planet is a large mass, planet or planetoid in the Star Wars Universe, at the time of 0 ABY.
type Planet struct {
	Name           string   `json:"name"`
	RotationPeriod string   `json:"rotation_period"`
	OrbitalPeriod  string   `json:"orbital_period"`
	Diameter       string   `json:"diameter"`
	Climate        string   `json:"climate"`
	Gravity        string   `json:"gravity"`
	Terrain        string   `json:"terrain"`
	SurfaceWater   string   `json:"surface_water"`
	Population     string   `json:"population"`
	ResidentURLs   []string `json:"residents"`
	FilmURLs       []string `json:"films"`
	Created        string   `json:"created"`
	Edited         string   `json:"edited"`
	URL            string   `json:"url"`
}

var ErrPlanetNotFound = errors.New("planet not found")

func (c *Client) GetPlanetByName(ctx context.Context, planetName string) (Planet, error) {
	req, err := c.newRequest(ctx, "/planets", url.Values{"search": {planetName}})
	if err != nil {
		return Planet{}, err
	}

	var result searchResult

	if _, err = c.do(req, &result); err != nil {
		return Planet{}, err
	}

	if len(result.Planets) < 1 {
		return Planet{}, ErrPlanetNotFound
	}

	return result.Planets[0], nil
}

func (c *Client) GetPlanetAppearances(ctx context.Context, planetName string) (int, error) {
	planet, err := c.GetPlanetByName(ctx, planetName)

	if err != nil {
		return 0, err
	}

	return len(planet.FilmURLs), nil
}
