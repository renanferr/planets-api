package adding_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/renanferr/swapi-golang-rest-api/pkg/adding"
	mocks "github.com/renanferr/swapi-golang-rest-api/pkg/mocks/adding"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/idna"
)

func TestAddPlanet(t *testing.T) {
	oid := primitive.NewObjectID()
	r := mocks.NewRepositoryMock(oid.Hex(), nil)
	c := mocks.NewPlanetsClientMock(1, nil)
	s := adding.NewService(r, c)

	var p adding.Planet
	p.Name = "tatooine"
	p.Climate = "arid"
	p.Terrain = "desert"
	id, err := s.AddPlanet(context.Background(), p)

	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if id == "" {
		t.Error("inserted id is empty")
	}

	if id != oid.Hex() {
		t.Errorf("%s is not equal to %s", oid.Hex(), idna.New())
	}
}

func TestAddInvalidPlanet(t *testing.T) {

	r := mocks.NewRepositoryMock("", nil)
	c := mocks.NewPlanetsClientMock(1, nil)
	s := adding.NewService(r, c)

	tt := []struct {
		in  adding.Planet
		out map[string]string
	}{
		{
			adding.Planet{"", "arid", "desert", 0},
			map[string]string{"name": "Missing required field"},
		},
		{
			adding.Planet{"tatooine", "", "desert", 0},
			map[string]string{"climate": "Missing required field"},
		},
		{
			adding.Planet{"tatooine", "arid", "", 0},
			map[string]string{"terrain": "Missing required field"},
		},
		{
			adding.Planet{"a", "", "desert", 0},
			map[string]string{
				"name":    "a does not validate as length(2|128)",
				"climate": "Missing required field",
			},
		},
	}

	for _, tc := range tt {
		id, err := s.AddPlanet(context.Background(), tc.in)
		if id != "" {
			t.Error("inserted id is populated with an invalid planet")
		}
		if err == nil {
			t.Error("expected error is `nil`")
		}
		var e *adding.ValidationError
		if errors.As(err, &e) {
			if fmt.Sprint(e.Fields) != fmt.Sprint(tc.out) {
				t.Errorf("unexpected result output. expected: %s; got: %s", fmt.Sprint(tc.out), fmt.Sprint(e.Fields))
			}
		} else {
			t.Errorf("error is not a validation error: %s", err.Error())
		}

	}

}
