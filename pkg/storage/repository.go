package storage

import (
	"context"
	"log"
	"time"

	"github.com/renanferr/swapi-golang-rest-api/pkg/adding"
	"github.com/renanferr/swapi-golang-rest-api/pkg/listing"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// CollectionPlanet identifier for the JSON collection of planets
	DatabaseName     = "tech-challenge"
	CollectionPlanet = "planets"
)

// Storage stores planet data in JSON files
type Storage struct {
	client *mongo.Client
}

// NewStorage returns a new JSON  storage
func NewStorage() (*Storage, error) {
	var err error

	s := new(Storage)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	defer func() {
		if err = s.client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	if err != nil {
		return nil, err
	}

	return s, nil
}

// AddPlanet saves the given planet to the repository
func (s *Storage) AddPlanet(ctx context.Context, p adding.Planet) error {

	planet := Planet{
		ID:          primitive.NewObjectID(),
		CreatedAt:   time.Now(),
		Name:        p.Name,
		Climate:     p.Climate,
		Terrain:     p.Terrain,
		Appearances: p.Appearances,
	}

	b, err := bson.Marshal(planet)

	if err != nil {
		log.Fatalf("error marshaling planet: %v", err)
	}

	if _, err := s.client.Database(DatabaseName).Collection(CollectionPlanet).InsertOne(ctx, b); err != nil {
		return err
	}
	return nil
}

// Get returns a planet with the specified ID
func (s *Storage) GetPlanet(ctx context.Context, id string) (listing.Planet, error) {
	var p Planet
	var planet listing.Planet
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return planet, listing.ErrNotFound

	}
	filter := bson.M{"id": oid}
	if err = s.client.Database(DatabaseName).Collection(CollectionPlanet).FindOne(ctx, filter).Decode(&p); err != nil {
		// err handling omitted for simplicity
		return planet, listing.ErrNotFound
	}

	planet.ID = p.ID.Hex()
	planet.Name = p.Name
	planet.Climate = p.Climate
	planet.Terrain = p.Terrain
	planet.Appearances = p.Appearances

	return planet, nil
}

// GetPlanets returns all planets
func (s *Storage) GetPlanets(ctx context.Context) []listing.Planet {
	list := []listing.Planet{}

	cur, err := s.client.Database(DatabaseName).Collection(CollectionPlanet).Find(ctx, bson.M{})

	if err != nil {
		panic(err)
	}

	defer cur.Close(ctx)

	if cur.Err() != nil {
		panic(cur.Err())
	}

	for cur.Next(ctx) {
		var p Planet
		err := cur.Decode(&p)
		if err != nil {
			return list
		}
		var planet listing.Planet

		planet.ID = p.ID.Hex()
		planet.Name = p.Name
		planet.Climate = p.Climate
		planet.Terrain = p.Terrain
		planet.Appearances = p.Appearances

		list = append(list, planet)
	}

	return list
}
