package storage

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
