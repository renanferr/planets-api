package listing

import (
	"context"
	"errors"

	"github.com/renanferr/swapi-golang-rest-api/pkg/listing"
	"github.com/renanferr/swapi-golang-rest-api/pkg/mocks"
)

type RepositoryMock mocks.Mock

func NewRepositoryMock(v interface{}, err error) *RepositoryMock {
	return &RepositoryMock{v, err}
}

func (m *RepositoryMock) GetPlanet(ctx context.Context, planetID string) (listing.Planet, error) {
	v, ok := m.Value.(listing.Planet)
	if !ok {
		return listing.Planet{}, errors.New("unexpected value type. expected: string")
	}

	return v, m.Err
}

func (m *RepositoryMock) GetPlanets(ctx context.Context, limit int, offset int) []listing.Planet {
	return m.Value.([]listing.Planet)
}
