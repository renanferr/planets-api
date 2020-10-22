package adding

import (
	"context"
	"errors"

	"github.com/renanferr/swapi-golang-rest-api/pkg/adding"
	"github.com/renanferr/swapi-golang-rest-api/pkg/mocks"
)

type ServiceMock mocks.Mock

func NewServiceMock(v interface{}, err error) *ServiceMock {
	return &ServiceMock{v, err}
}

func (m *ServiceMock) AddPlanet(ctx context.Context, planet adding.Planet) (string, error) {
	v, ok := m.Value.(string)
	if !ok {
		return "", errors.New("unexpected value type. expected: adding.Planet")
	}

	return v, m.Err
}
