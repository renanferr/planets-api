package adding

import (
	"context"
	"errors"
)

// Event defines possible outcomes from the "adding actor"
type Event int

const (
	// Done means finished processing successfully
	Done Event = iota

	// InvalidPlanet means the given planet is invalid
	InvalidPlanet

	// PlanetAlreadyExists means the given planet is a duplicate of an existing one
	PlanetAlreadyExists

	// Failed means processing did not finish successfully
	Failed
)

func (e Event) Get() string {
	if e == Done {
		return "Done"
	}

	if e == PlanetAlreadyExists {
		return "Duplicate planet"
	}

	if e == Failed {
		return "Failed"
	}

	return "Unknown result"
}

var ErrDuplicate = errors.New("planet already exists")

// Service provides planet adding operations.
type Service interface {
	// AddPlanet invokes the operations needed to save a planet
	AddPlanet(context.Context, Planet) error
}

// PlanetsClient provides an interface to fetch extra planets info from a third-party API
type PlanetsClient interface {
	GetPlanetAppearances(context.Context, string) (int, error)
}

// Repository provides access to planet repository.
type Repository interface {
	// AddPlanet saves a given planet to the repository.
	AddPlanet(context.Context, Planet) error
}

type service struct {
	repo    Repository
	planets PlanetsClient
}

// NewService creates an adding service with the necessary dependencies
func NewService(r Repository, p PlanetsClient) Service {
	return &service{r, p}
}

// AddPlanet adds the given planet(s) to the database
func (s *service) AddPlanet(ctx context.Context, p Planet) error {
	var err error
	p.Appearances, err = s.planets.GetPlanetAppearances(ctx, p.Name)
	if err != nil {
		panic(err)
	}
	return s.repo.AddPlanet(ctx, p)
}
