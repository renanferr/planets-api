package listing_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/renanferr/swapi-golang-rest-api/pkg/listing"
	mocks "github.com/renanferr/swapi-golang-rest-api/pkg/mocks/listing"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetPlanet(t *testing.T) {
	oid := primitive.NewObjectID()

	var p listing.Planet
	p.Name = "tatooine"
	p.Climate = "arid"
	p.Terrain = "desert"

	r := mocks.NewRepositoryMock(p, nil)

	s := listing.NewService(r)

	planet, err := s.GetPlanet(context.Background(), oid.Hex())

	if err != nil {
		t.Error("unexpected error: %w", err)
	}
	if planet.ID != p.ID {
		t.Errorf("planet `ID` does not match: %s != %s", p.ID, planet.ID)
	}
	if planet.Name != p.Name {
		t.Errorf("planet `Name` does not match: %s != %s", p.Name, planet.Name)
	}
	if planet.Climate != p.Climate {
		t.Errorf("planet `Climate` does not match: %s != %s", p.Climate, planet.Climate)
	}
	if planet.Terrain != p.Terrain {
		t.Errorf("planet `Terrain` does not match: %s != %s", p.Terrain, planet.Terrain)
	}
}

func TestGetPlanetNotFound(t *testing.T) {
	tt := []string{
		"test",
		primitive.NewObjectID().Hex(),
	}

	r := mocks.NewRepositoryMock(listing.Planet{}, listing.ErrPlanetNotFound)
	s := listing.NewService(r)

	for _, id := range tt {

		_, err := s.GetPlanet(context.Background(), id)

		if err == nil {
			t.Error("expected error is nil")
		}

		if !errors.Is(err, listing.ErrPlanetNotFound) {
			t.Errorf("unexpected error: %w", err)
		}
	}

}

func TestGetPlanets(t *testing.T) {

	tt := [][]listing.Planet{
		{},
		{
			{
				primitive.NewObjectID().Hex(),
				"tatooine",
				"arid",
				"desert",
				5,
			},
		},
		{
			{
				primitive.NewObjectID().Hex(),
				"tatooine",
				"arid",
				"desert",
				5,
			},
			{
				primitive.NewObjectID().Hex(),
				"alderaan",
				"temperate",
				"grasslands",
				2,
			},
		},
	}

	for i, tc := range tt {
		r := mocks.NewRepositoryMock(tt[i], nil)
		s := listing.NewService(r)

		planets := s.GetPlanets(context.Background())
		if !reflect.DeepEqual(tc, planets) {
			t.Errorf("planet lists does not match: %v != %v", tc, planets)
		}
	}
}
