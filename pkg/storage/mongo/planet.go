package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/renanferr/swapi-golang-rest-api/pkg/adding"
	"github.com/renanferr/swapi-golang-rest-api/pkg/listing"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

// Planet defines the storage form of a planet
type Planet struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
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
		Name:        p.Name,
		Climate:     p.Climate,
		Terrain:     p.Terrain,
		Appearances: p.Appearances,
	}

	b, err := bson.Marshal(planet)

	if err != nil {
		return "", fmt.Errorf("error marshalling planet: %v", err)
	}

	if _, err := s.Client.Database(DatabaseName).Collection(PlanetsCollection).InsertOne(ctx, b); err != nil {
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
		return planet, listing.ErrPlanetNotFound

	}

	filter := bson.M{"_id": oid}

	ctx, cancel := context.WithTimeout(ctx, s.TimeoutMS)
	defer cancel()

	collection := s.Client.Database(DatabaseName).Collection(PlanetsCollection)
	err = collection.FindOne(ctx, filter).Decode(&p)

	if err != nil {
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
func (s *Storage) GetPlanets(ctx context.Context) []listing.Planet {
	list := []listing.Planet{}

	cur, err := s.Client.Database(DatabaseName).Collection(PlanetsCollection).Find(ctx, bson.M{})

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
		log.Printf("decoded %v", p)

		if err != nil {
			log.Printf("error decoding planet: %s", err.Error())
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
