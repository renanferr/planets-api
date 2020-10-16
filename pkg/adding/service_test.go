package adding

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mock struct {
	Value interface{}
	Err   error
}

type repositoryMock = mock

type clientMock = mock

func (m *repositoryMock) AddPlanet(ctx context.Context, p Planet) (string, error) {
	return m.Value.(string), m.Err
}

func (m *clientMock) GetPlanetAppearances(ctx context.Context, planetName string) (int, error) {
	return m.Value.(int), m.Err
}

func TestAddingService(t *testing.T) {
	oid := primitive.NewObjectID()
	r := &mock{oid.Hex(), nil}
	c := &clientMock{1, nil}
	s := NewService(r, c)
	var p Planet
	p.Name = "tatooine"
	p.Climate = "arid"
	p.Terrain = "desert"
	id, err := s.AddPlanet(context.Background(), p)
	if err != nil {
		t.Error(err.Error())
	}

	if id == "" {
		t.Error("id is empty")
	}

	if oid.Hex() != id {
		t.Errorf("%s is not equal to %s", oid.Hex(), id)
	}
}
