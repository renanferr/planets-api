package adding

import (
	"context"
	"errors"
)

type Payload []*Planet

// Event defines possible outcomes from the "adding actor"
type Event int

const (
	// Done means finished processing successfully
	Done Event = iota

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
	AddPlanet(context.Context, ...Planet)
}

// Repository provides access to planet repository.
type Repository interface {
	// AddPlanet saves a given planet to the repository.
	AddPlanet(context.Context, Planet) error
}

type service struct {
	repo Repository
}

// NewService creates an adding service with the necessary dependencies
func NewService(r Repository) Service {
	return &service{r}
}

// AddPlanet adds the given planet(s) to the database
func (s *service) AddPlanet(ctx context.Context, p Planet) {

	validatePlanet(p)
	for _, planet := range p {
		_ = s.repo.AddPlanet(ctx, planet) // error handling omitted for simplicity
	}

}

func getPlanetAppearances(p *Planet) {

}
