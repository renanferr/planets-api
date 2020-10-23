package adding

import (
	"context"
	"errors"
	"log"

	"github.com/asaskevich/govalidator"
)

var ErrPlanetNotFound = errors.New("planet not found")

// ValidationError defines the type for a validation error
type ValidationError struct {
	Err    error
	Fields map[string]string
}

func NewValidationError(err error) *ValidationError {
	return &ValidationError{
		err,
		govalidator.ErrorsByField(err),
	}
}

func (e ValidationError) Error() string {
	return e.Err.Error()
}

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

	responseChan := make(chan struct {
		appearances int
		err         error
	})

	go func() {
		a, e := s.planets.GetPlanetAppearances(ctx, p.Name)

		responseChan <- struct {
			appearances int
			err         error
		}{a, e}
	}()

	validationErr := <-validationErrChan
	if validationErr != nil {
		return "", NewValidationError(validationErr)
	}

	response := <-responseChan

	if response.err != nil {
		log.Printf("error fetching planet info: %s", response.err.Error())
		if !errors.Is(ErrPlanetNotFound, response.err) {
			return "", response.err
		}
	}

	p.Appearances = response.appearances

	id, err := s.repo.AddPlanet(ctx, p)
	if err != nil {
		return "", err
	}
	return id, nil
}
