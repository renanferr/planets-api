package mongo

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/renanferr/planets-api/pkg/adding"
	"github.com/renanferr/planets-api/pkg/listing"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Planet defines the storage form of a planet
type Planet struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Climate     string             `json:"climate" bson:"climate"`
	Terrain     string             `json:"terrain" bson:"terrain"`
	Appearances int                `json:"appearances" bson:"appearances"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}

const PlanetsCollection = "planets"

// AddPlanet saves the given planet to the repository
func (s *Storage) AddPlanet(ctx context.Context, p adding.Planet) (string, error) {

	planet := Planet{
		ID:          primitive.NewObjectID(),
		CreatedAt:   time.Now(),
		Name:        strings.ToLower(p.Name),
		Climate:     strings.ToLower(p.Climate),
		Terrain:     strings.ToLower(p.Terrain),
		Appearances: p.Appearances,
	}

	b, err := bson.Marshal(planet)

	if err != nil {
		return "", fmt.Errorf("error marshalling planet: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, s.TimeoutMS)
	defer cancel()

	if _, err := s.Client.Database(s.DatabaseName()).Collection(PlanetsCollection).InsertOne(ctx, b); err != nil {
		return "", err
	}

	log.Printf("Inserted Planet with ID %s\n", planet.ID.Hex())

	return planet.ID.Hex(), err
}

// Get returns a planet with the specified ID
func (s *Storage) GetPlanet(ctx context.Context, id string) (listing.Planet, error) {
	var p Planet
	var planet listing.Planet

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("error casting \"%s\" to ObjectId: %s", id, err.Error())
		return planet, listing.ErrPlanetNotFound

	}

	filter := bson.M{"_id": oid}

	ctx, cancel := context.WithTimeout(ctx, s.TimeoutMS)
	defer cancel()

	collection := s.Client.Database(s.DatabaseName()).Collection(PlanetsCollection)
	err = collection.FindOne(ctx, filter).Decode(&p)

	if err != nil {
		log.Printf("error finding planet: %s", err.Error())
		return planet, listing.ErrPlanetNotFound

	}

	planet.ID = p.ID.Hex()
	planet.Name = p.Name
	planet.Climate = p.Climate
	planet.Terrain = p.Terrain
	planet.Appearances = p.Appearances

	return planet, nil
}

// GetPlanets returns all planets
func (s *Storage) GetPlanets(ctx context.Context, offset int64, limit int64) ([]listing.Planet, int64) {
	list := []listing.Planet{}

	ctx, cancel := context.WithTimeout(ctx, s.TimeoutMS)
	defer cancel()

	collection := s.Client.Database(s.DatabaseName()).Collection(PlanetsCollection)
	filter := bson.M{}

	count, err := collection.Count(ctx, filter)

	if err != nil {
		panic(err)
	}

	opts := options.Find().
		SetLimit(limit).
		SetSkip(offset)

	cur, err := collection.Find(ctx, filter, opts)

	if err != nil {
		panic(err)
	}

	defer cur.Close(ctx)

	if err = cur.Err(); err != nil {
		panic(err)
	}

	for cur.Next(ctx) {
		var p Planet
		err := cur.Decode(&p)

		if err != nil {
			log.Printf("error decoding planet: %s", err.Error())
			return list, 0
		}
		var planet listing.Planet

		planet.ID = p.ID.Hex()
		planet.Name = p.Name
		planet.Climate = p.Climate
		planet.Terrain = p.Terrain
		planet.Appearances = p.Appearances

		list = append(list, planet)
	}

	return list, count
}
