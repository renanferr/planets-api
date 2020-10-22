package adding

import (
	"context"
	"errors"

	"github.com/renanferr/swapi-golang-rest-api/pkg/adding"
	"github.com/renanferr/swapi-golang-rest-api/pkg/mocks"
)

type RepositoryMock mocks.Mock

func NewRepositoryMock(v interface{}, err error) *RepositoryMock {
	return &RepositoryMock{v, err}
}

func (m *RepositoryMock) AddPlanet(ctx context.Context, planet adding.Planet) (string, error) {
	v, ok := m.Value.(string)
	if !ok {
		return "", errors.New("unexpected value type. expected: string")
	}

	return v, m.Err
}
