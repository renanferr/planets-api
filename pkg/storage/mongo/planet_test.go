package mongo_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/renanferr/swapi-golang-rest-api/pkg/adding"
	"github.com/renanferr/swapi-golang-rest-api/pkg/listing"
	mocks "github.com/renanferr/swapi-golang-rest-api/pkg/mocks/storage/mongo"
	"github.com/renanferr/swapi-golang-rest-api/pkg/storage/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodb "go.mongodb.org/mongo-driver/mongo"
)

func TestAddPlanet(t *testing.T) {

	storage := mockStorage(&mongodb.InsertOneResult{})

	var p adding.Planet
	p.Name = "tatooine"
	p.Climate = "arid"
	p.Terrain = "desert"
	p.Appearances = 5

	id, err := storage.AddPlanet(context.Background(), p)
	if err != nil {
		t.Fatal(err)
	}

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		t.Errorf("error casting \"%s\" to ObjectID: %s", id, err.Error())
	}

	if id != oid.Hex() {
		t.Errorf("unexpected inserted id. got: %s want %s", id, oid.Hex())
	}
}

func mockStorage(value interface{}) *mongo.Storage {
	coll := &mocks.CollectionMock{value, nil}
	db := &mocks.DatabaseMock{coll}
	cli := &mocks.ClientMock{db, nil, nil, nil}
	return mocks.NewStorageMock("mongodb://mock.test", cli)
}

func mockSingleResult(t *testing.T, p *mongo.Planet, err error) *mocks.SingleResultMock {
	b, e := bson.Marshal(p)
	if e != nil {
		t.Fatalf("error marshalling planet: %s", e.Error())
	}
	return &mocks.SingleResultMock{
		Value: b,
		Err:   err,
	}
}

func mockCursor(t *testing.T, planets []*mongo.Planet, err error) *mocks.CursorMock {
	return mocks.NewCursorMock(planets, err)
}

func TestGetPlanet(t *testing.T) {

	type TestCase struct {
		Name     string
		PlanetID string
		Mock     *mocks.SingleResultMock
		Err      error
	}

	oid := primitive.NewObjectID()
	tt := []*TestCase{
		{
			Name:     "get planet successfully",
			PlanetID: oid.Hex(),
			Mock: mockSingleResult(
				t,
				&mongo.Planet{
					ID:          oid,
					Name:        "tatooine",
					Climate:     "arid",
					Terrain:     "desert",
					Appearances: 5,
				},
				nil,
			),
		},
		{
			Name:     "invalid id",
			PlanetID: "test",
			Mock:     mockSingleResult(t, &mongo.Planet{}, nil),
			Err:      listing.ErrPlanetNotFound,
		},
		{
			Name:     "planet not found",
			PlanetID: primitive.NewObjectID().Hex(),
			Mock: mockSingleResult(
				t,
				&mongo.Planet{},
				errors.New("decoding error"),
			),
			Err: listing.ErrPlanetNotFound,
		},
	}

	for _, tc := range tt {

		storage := mockStorage(tc.Mock)

		planet, err := storage.GetPlanet(context.Background(), tc.PlanetID)

		expectedPlanet := &listing.Planet{}

		if !errors.Is(err, tc.Err) {
			t.Errorf("<%s> unexpected error. got: %s want: %s", tc.Name, err, tc.Err)
		}

		var p mongo.Planet
		bson.Unmarshal(tc.Mock.Value.([]byte), &p)

		if err == nil {
			expectedPlanet.ID = p.ID.Hex()
			expectedPlanet.Name = p.Name
			expectedPlanet.Climate = p.Climate
			expectedPlanet.Terrain = p.Terrain
			expectedPlanet.Appearances = p.Appearances
		}

		if !reflect.DeepEqual(planet, *expectedPlanet) {
			t.Errorf("<%s> planets do not match. got: %v want: %v", tc.Name, planet, *expectedPlanet)
		}
	}

}

func TestGetPlanets(t *testing.T) {
	type TestCase struct {
		Name string
		Mock *mocks.CursorMock
		Err  error
	}

	oid := primitive.NewObjectID()
	oid2 := primitive.NewObjectID()
	tt := []*TestCase{
		{
			Name: "get 1 planet",
			Mock: mockCursor(
				t,
				[]*mongo.Planet{
					{
						ID:          oid,
						Name:        "tatooine",
						Climate:     "arid",
						Terrain:     "desert",
						Appearances: 5,
					},
				},
				nil,
			),
		},
		{
			Name: "get 2 planets",
			Mock: mockCursor(
				t,
				[]*mongo.Planet{
					{
						ID:          oid,
						Name:        "tatooine",
						Climate:     "arid",
						Terrain:     "desert",
						Appearances: 5,
					},
					{
						ID:          oid2,
						Name:        "alderaan",
						Climate:     "temperate",
						Terrain:     "grasslands",
						Appearances: 2,
					},
				},
				nil,
			),
		},
		{
			Name: "get 0 planets",
			Mock: mockCursor(
				t,
				[]*mongo.Planet{},
				nil,
			),
		},
		{
			Name: "cursor error",
			Mock: mockCursor(
				t,
				[]*mongo.Planet{},
				errors.New("cursor error"),
			),
			Err: errors.New("cursor error"),
		},
	}

	for _, tc := range tt {
		storage := mockStorage(tc.Mock)

		if tc.Err != nil {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("GetPlanets did not panic")
				}
			}()
		}

		out := storage.GetPlanets(context.Background(), 20, 0)
		var expectedOutput []listing.Planet
		for _, p := range tc.Mock.Values.([]*mongo.Planet) {
			var planet listing.Planet
			planet.ID = p.ID.Hex()
			planet.Name = p.Name
			planet.Climate = p.Climate
			planet.Terrain = p.Terrain
			planet.Appearances = p.Appearances
			expectedOutput = append(expectedOutput, planet)
		}

		if len(out) != 0 &&
			len(expectedOutput) != 0 &&
			!reflect.DeepEqual(out, expectedOutput) {

			t.Errorf("<%s> planets do not match. got: %v want: %v", tc.Name, out, expectedOutput)
		}
	}
}
