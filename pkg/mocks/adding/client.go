package adding

import (
	"context"
	"errors"

	"github.com/renanferr/planets-api/pkg/mocks"
)

type PlanetsClientMock mocks.Mock

func NewPlanetsClientMock(v interface{}, err error) *PlanetsClientMock {
	return &PlanetsClientMock{v, err}
}

func (m *PlanetsClientMock) GetPlanetAppearances(ctx context.Context, planetName string) (int, error) {
	v, ok := m.Value.(int)
	if !ok {
		return 0, errors.New("unexpected value type. expected: adding.Planet")
	}

	return v, m.Err
}
