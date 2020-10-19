package adding

import (
	"context"
	"errors"
)

type mock struct {
	Value interface{}
	Err   error
}

type AddingMock mock
type RepositoryMock mock
type ClientMock mock

func NewAddingMock(v interface{}, err error) *AddingMock {
	return &AddingMock{v, err}
}

func (m *AddingMock) AddPlanet(ctx context.Context, p Planet) (string, error) {
	v, ok := m.Value.(string)
	if !ok {
		return "", errors.New("unexpected value type. expected: string")
	}

	return v, m.Err
}

func NewRepositoryMock(v interface{}, err error) *RepositoryMock {
	return &RepositoryMock{v, err}
}

func (m *RepositoryMock) AddPlanet(ctx context.Context, p Planet) (string, error) {
	v, ok := m.Value.(string)
	if !ok {
		return "", errors.New("unexpected value type. expected: string")
	}

	return v, m.Err
}

func NewClientMock(v interface{}, err error) *ClientMock {
	return &ClientMock{v, err}
}

func (m *ClientMock) GetPlanetAppearances(ctx context.Context, planetName string) (int, error) {
	s, ok := m.Value.(int)
	if !ok {
		return 0, errors.New("unexpected value type. expected: int")
	}

	return s, m.Err
}
