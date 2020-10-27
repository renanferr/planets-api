package listing

import (
	"context"
	"errors"

	"github.com/renanferr/swapi-golang-rest-api/pkg/listing"
)

type RepositoryMock struct {
	Value interface{}
	Err   error
	Total int64
}

func NewRepositoryMock(v interface{}, err error, total int64) *RepositoryMock {
	return &RepositoryMock{v, err, total}
}

func (m *RepositoryMock) GetPlanet(ctx context.Context, planetID string) (listing.Planet, error) {
	v, ok := m.Value.(listing.Planet)
	if !ok {
		return listing.Planet{}, errors.New("unexpected value type. expected: string")
	}

	return v, m.Err
}

func (m *RepositoryMock) GetPlanets(ctx context.Context, limit int64, offset int64) ([]listing.Planet, int64) {
	return m.Value.([]listing.Planet), m.Total
}
