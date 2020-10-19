package adding

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/idna"
)

func TestAddPlanet(t *testing.T) {
	oid := primitive.NewObjectID()
	r := NewRepositoryMock(oid.Hex(), nil)
	c := NewClientMock(1, nil)
	s := NewService(r, c)

	var p Planet
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

	r := NewRepositoryMock("", nil)
	c := NewClientMock(1, nil)
	s := NewService(r, c)

	tt := []struct {
		in  Planet
		out map[string]string
	}{
		{
			Planet{"", "arid", "desert", 0},
			map[string]string{"name": "Missing required field"},
		},
		{
			Planet{"tatooine", "", "desert", 0},
			map[string]string{"climate": "Missing required field"},
		},
		{
			Planet{"tatooine", "arid", "", 0},
			map[string]string{"terrain": "Missing required field"},
		},
		{
			Planet{"a", "", "desert", 0},
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
		var e *ValidationError
		if errors.As(err, &e) {
			if fmt.Sprint(e.Fields) != fmt.Sprint(tc.out) {
				t.Errorf("unexpected result output. expected: %s; got: %s", fmt.Sprint(tc.out), fmt.Sprint(e.Fields))
			}
		} else {
			t.Errorf("error is not a validation error: %s", err.Error())
		}

	}

}
