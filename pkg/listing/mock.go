package listing

import (
	"context"
	"errors"
	"log"
)

type mock struct {
	Value interface{}
	Err   error
}

type ListingMock mock
type RepositoryMock mock

func NewListingMock(v interface{}, err error) *ListingMock {
	return &ListingMock{v, err}
}

func (m *ListingMock) GetPlanet(ctx context.Context, id string) (Planet, error) {
	v, ok := m.Value.(Planet)
	if !ok {
		return Planet{}, errors.New("unexpected value type. expected: Planet")
	}

	return v, m.Err
}

func (m *ListingMock) GetPlanets(ctx context.Context) []Planet {
	v, ok := m.Value.([]Planet)
	if !ok {
		log.Panicf("could not assert %v of type `Planet`", m.Value)
	}

	return v
}
func NewRepositoryMock(v interface{}, err error) *RepositoryMock {
	return &RepositoryMock{v, err}
}

func (m *RepositoryMock) GetPlanet(ctx context.Context, planetID string) (Planet, error) {
	v, ok := m.Value.(Planet)
	if !ok {
		return Planet{}, errors.New("unexpected value type. expected: string")
	}

	return v, m.Err
}

func (m *RepositoryMock) GetPlanets(ctx context.Context) []Planet {
	return m.Value.([]Planet)
}
