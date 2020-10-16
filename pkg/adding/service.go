package adding

import (
	"context"
	"log"

	"github.com/asaskevich/govalidator"
)

type ErrInvalidPlanet = error

// Service provides planet adding operations.
type Service interface {
	// AddPlanet invokes the operations needed to save a planet
	AddPlanet(context.Context, Planet) (string, error)
}

// PlanetsClient provides an interface to fetch extra planets info from a third-party API
type PlanetsClient interface {
	GetPlanetAppearances(context.Context, string) (int, error)
}

// Repository provides access to planet repository.
type Repository interface {
	// AddPlanet saves a given planet to the repository.
	AddPlanet(context.Context, Planet) (string, error)
}

type service struct {
	repo    Repository
	planets PlanetsClient
}

// NewService creates an adding service with the necessary dependencies
func NewService(r Repository, p PlanetsClient) Service {
	govalidator.SetFieldsRequiredByDefault(true)
	return &service{r, p}
}

// AddPlanet adds the given planet(s) to the database
func (s *service) AddPlanet(ctx context.Context, p Planet) (string, error) {

	validationErrChan := make(chan error)
	go func() {
		isValid, validationErr := govalidator.ValidateStruct(p)
		if !isValid {
			validationErrChan <- validationErr
		}
		validationErrChan <- nil
	}()

	appearancesChan := make(chan int)
	go func() {
		a, e := s.planets.GetPlanetAppearances(ctx, p.Name)
		if e != nil {
			log.Printf("error fetching planet info: %s", e.Error())
		}
		appearancesChan <- a
	}()

	validationErr := <-validationErrChan
	if validationErr != nil {
		return "", validationErr.(ErrInvalidPlanet)
	}

	p.Appearances = <-appearancesChan

	return s.repo.AddPlanet(ctx, p)
}
