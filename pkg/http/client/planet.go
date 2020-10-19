package client

import (
	"context"
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

const (
	ErrPlanetNotFound = iota
)

func (c *Client) GetPlanetByName(ctx context.Context, planetName string) (Planet, error) {
	req, err := c.newRequest(ctx, "/planets", query{"search": planetName})
	if err != nil {
		return Planet{}, err
	}

	var result searchResult

	if _, err = c.do(req, &result); err != nil {
		return Planet{}, err
	}

	if len(result.Planets) > 0 {
		return result.Planets[0], nil
	}
	return Planet{}, nil
}

func (c *Client) GetPlanetAppearances(ctx context.Context, planetName string) (int, error) {
	planet, err := c.GetPlanetByName(ctx, planetName)

	if err != nil {
		return 0, err
	}

	return len(planet.FilmURLs), nil
}
