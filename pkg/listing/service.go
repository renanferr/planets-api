package listing

import (
	"context"
	"errors"
)

// ErrNotFound is used when a planet could not be found.
var (
	ErrNotFound = errors.New("planet not found")
)

// Repository provides access to the planet storage.
type Repository interface {
	// GetPlanet returns the planet with given ID.
	GetPlanet(context.Context, string) (Planet, error)
	// GetPlanets returns all planets saved in storage.
	GetPlanets(context.Context) []Planet
}

// Service provides planet listing operations.
type Service interface {
	GetPlanet(context.Context, string) (Planet, error)
	GetPlanets(context.Context) []Planet
}

type service struct {
	r Repository
}

// NewService creates a listing service with the necessary dependencies
func NewService(r Repository) Service {
	return &service{r}
}

// GetPlanets returns all planets
func (s *service) GetPlanets(ctx context.Context) []Planet {
	return s.r.GetPlanets(ctx)
}

// GetPlanet returns a planet
func (s *service) GetPlanet(ctx context.Context, id string) (Planet, error) {
	return s.r.GetPlanet(ctx, id)
}