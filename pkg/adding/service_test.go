package adding_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/renanferr/swapi-golang-rest-api/pkg/adding"
	mocks "github.com/renanferr/swapi-golang-rest-api/pkg/mocks/adding"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// func TestAddPlanet(t *testing.T) {
// 	oid := primitive.NewObjectID()
// 	r := mocks.NewRepositoryMock(oid.Hex(), nil)
// 	c := mocks.NewPlanetsClientMock(1, nil)
// 	s := adding.NewService(r, c)

// 	var p adding.Planet
// 	p.Name = "tatooine"
// 	p.Climate = "arid"
// 	p.Terrain = "desert"
// 	id, err := s.AddPlanet(context.Background(), p)

// 	if err != nil {
// 		t.Errorf("unexpected error: %s", err.Error())
// 	}

// 	if id == "" {
// 		t.Error("inserted id is empty")
// 	}

// 	if id != oid.Hex() {
// 		t.Errorf("%s is not equal to %s", oid.Hex(), idna.New())
// 	}
// }

func TestAddPlanet(t *testing.T) {

	type TestCase struct {
		Name             string
		Planet           adding.Planet
		Repository       *mocks.RepositoryMock
		PlanetsClient    *mocks.PlanetsClientMock
		ValidationErrMap map[string]string
	}

	tt := []*TestCase{
		{
			"Add Planet successfully",
			adding.Planet{"tatooine", "arid", "desert", 0},
			mocks.NewRepositoryMock(primitive.NewObjectID().Hex(), nil),
			mocks.NewPlanetsClientMock(5, nil),
			map[string]string{},
		},
		{
			"Missing name field",
			adding.Planet{"", "arid", "desert", 0},
			mocks.NewRepositoryMock(primitive.NewObjectID().Hex(), nil),
			mocks.NewPlanetsClientMock(0, nil),
			map[string]string{"name": "Missing required field"},
		},
		{
			"Missing climate field",
			adding.Planet{"tatooine", "", "desert", 0},
			mocks.NewRepositoryMock(primitive.NewObjectID().Hex(), nil),
			mocks.NewPlanetsClientMock(0, nil),
			map[string]string{"climate": "Missing required field"},
		},
		{
			"Missing terrain field",
			adding.Planet{"tatooine", "arid", "", 0},
			mocks.NewRepositoryMock(primitive.NewObjectID().Hex(), nil),
			mocks.NewPlanetsClientMock(0, nil),
			map[string]string{"terrain": "Missing required field"},
		},
		{
			"Invalid name and missing climate field",
			adding.Planet{"a", "", "desert", 0},
			mocks.NewRepositoryMock(primitive.NewObjectID().Hex(), nil),
			mocks.NewPlanetsClientMock(5, nil),
			map[string]string{
				"name":    "a does not validate as length(2|128)",
				"climate": "Missing required field",
			},
		},
		{
			"Repository error",
			adding.Planet{"tatooine", "arid", "desert", 0},
			mocks.NewRepositoryMock("", errors.New("repository error")),
			mocks.NewPlanetsClientMock(5, nil),
			map[string]string{},
		},
		{
			"Planets client error",
			adding.Planet{"tatooine", "arid", "desert", 0},
			mocks.NewRepositoryMock("", nil),
			mocks.NewPlanetsClientMock(0, errors.New("planets client error")),
			map[string]string{},
		},
	}

	for _, tc := range tt {

		s := adding.NewService(tc.Repository, tc.PlanetsClient)

		id, err := s.AddPlanet(context.Background(), tc.Planet)

		if err != nil {
			if len(tc.ValidationErrMap) > 0 {
				var e *adding.ValidationError
				if errors.As(err, &e) {
					if fmt.Sprint(e.Fields) != fmt.Sprint(tc.ValidationErrMap) {
						t.Errorf("<%s> unexpected result output. expected: %s; got: %s", tc.Name, fmt.Sprint(tc.ValidationErrMap), fmt.Sprint(e.Fields))
					}
				} else {
					t.Errorf("<%s> error is not a validation error: %s", tc.Name, err.Error())
				}
			} else if !errors.Is(err, tc.Repository.Err) && !errors.Is(err, tc.PlanetsClient.Err) {
				t.Fatalf("<%s> unexpected err: %s", tc.Name, err.Error())
			}
		} else {
			if id == "" {
				t.Errorf("<%s> inserted id is empty", tc.Name)
			}

			if _, err = primitive.ObjectIDFromHex(id); err != nil {
				t.Fatalf("error casting inserted id to ObjectID: %s", err.Error())
			}
		}

	}

}

func TestValidationErr(t *testing.T) {
	e := errors.New("mock error")
	validationErr := adding.NewValidationError(e)
	if !errors.Is(validationErr.Err, e) || e.Error() != validationErr.Error() {
		t.Errorf("validation error and error do not match")
	}

}
