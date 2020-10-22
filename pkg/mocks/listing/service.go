package listing

import (
	"context"
	"errors"
	"log"

	"github.com/renanferr/swapi-golang-rest-api/pkg/listing"
	"github.com/renanferr/swapi-golang-rest-api/pkg/mocks"
)

type ServiceMock mocks.Mock

func NewServiceMock(v interface{}, err error) *ServiceMock {
	return &ServiceMock{v, err}
}

func (m *ServiceMock) GetPlanet(ctx context.Context, id string) (listing.Planet, error) {
	v, ok := m.Value.(listing.Planet)
	if !ok {
		return listing.Planet{}, errors.New("unexpected value type. expected: listing.Planet")
	}

	return v, m.Err
}

func (m *ServiceMock) GetPlanets(ctx context.Context) []listing.Planet {
	v, ok := m.Value.([]listing.Planet)
	if !ok {
		log.Panicf("could not assert %v of type `listing.Planet`", m.Value)
	}

	return v
}
