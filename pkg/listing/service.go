package listing

import (
	"context"
	"errors"
)

var (
	ErrPlanetNotFound error = errors.New("planet not found")
)

// Repository provides access to the planet storage.
type Repository interface {
	// GetPlanet returns the planet with given ID.
	GetPlanet(context.Context, string) (Planet, error)
	// GetPlanets returns all planets saved in storage.
	GetPlanets(context.Context, int, int) []Planet
}

// Service provides planet listing operations.
type Service interface {
	GetPlanet(context.Context, string) (Planet, error)
	GetPlanets(ctx context.Context, limit int, offset int) []Planet
}

type service struct {
	r Repository
}

// NewService creates a listing service with the necessary dependencies
func NewService(r Repository) Service {
	return &service{r}
}

// GetPlanets returns all planets
func (s *service) GetPlanets(ctx context.Context, limit int, offset int) []Planet {
	return s.r.GetPlanets(ctx, limit, offset)
}

// GetPlanet returns a planet
func (s *service) GetPlanet(ctx context.Context, id string) (Planet, error) {
	return s.r.GetPlanet(ctx, id)
}
