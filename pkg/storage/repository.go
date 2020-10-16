package storage

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/renanferr/swapi-golang-rest-api/pkg/adding"
	"github.com/renanferr/swapi-golang-rest-api/pkg/listing"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	// CollectionPlanet identifier for the JSON collection of planets
	DatabaseName     = "swapi-golang-rest-api"
	CollectionPlanet = "planets"
)

// Storage stores planet data in JSON files
type Storage struct {
	uri       *url.URL
	client    *mongo.Client
	timeoutMS time.Duration
}

// NewStorage returns a new JSON  storage
func NewStorage(uri string) *Storage {
	s := &Storage{}

	u, err := url.Parse(uri)

	if err != nil {
		log.Fatalf("invalid MongoDB Connection URI: %v", err)
	}

	s.uri = u
	s.client, err = mongo.NewClient(options.Client().ApplyURI(uri))

	if err != nil {
		log.Fatalf("error creating Mongo Client: %v", err)
	}

	return s
}

func (s *Storage) WithTimeout(timeout time.Duration) *Storage {
	s.timeoutMS = timeout
	return s
}

func (s *Storage) Connect(ctx context.Context) {
	log.Println("connecting to MongoDB")
	ctx, cancelFunc := context.WithTimeout(ctx, s.timeoutMS)
	defer cancelFunc()
	s.client.Connect(ctx)

	if err := s.client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("MongoDB server did not respond successfully: %s", err.Error())
	}
}

func (s *Storage) Disconnect(ctx context.Context) {
	log.Println("disconnecting from MongoDB")
	if err := s.client.Disconnect(ctx); err != nil {
		panic(err)
	}
}

// AddPlanet saves the given planet to the repository
func (s *Storage) AddPlanet(ctx context.Context, p adding.Planet) (string, error) {

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
		return "", fmt.Errorf("error marshalling planet: %v", err)
	}

	if _, err := s.client.Database(DatabaseName).Collection(CollectionPlanet).InsertOne(ctx, b); err != nil {
		return "", err
	}
	return planet.ID.Hex(), err
}

// Get returns a planet with the specified ID
func (s *Storage) GetPlanet(ctx context.Context, id string) (listing.Planet, error) {
	var p Planet
	var planet listing.Planet

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return planet, listing.ErrNotFound

	}

	filter := bson.M{"_id": oid}

	ctx, cancel := context.WithTimeout(ctx, s.timeoutMS)
	defer cancel()

	collection := s.client.Database(DatabaseName).Collection(CollectionPlanet)
	if err = collection.FindOne(ctx, filter).Decode(&p); err != nil {
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
