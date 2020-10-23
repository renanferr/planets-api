package mongo

import (
	"context"
	"log"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const DatabaseName = "swapi-golang-rest-api"

type Client interface {
	Ping(context.Context, *readpref.ReadPref) error
	Connect(context.Context) error
	Database(string) Database
	Disconnect(context.Context) error
}

type ClientWrapping struct {
	cli *mongo.Client
}

func (c *ClientWrapping) Database(name string) Database {
	return &DatabaseWrapping{c.cli.Database(name)}
}

func (c *ClientWrapping) Connect(ctx context.Context) error {
	return c.cli.Connect(ctx)
}

func (c *ClientWrapping) Disconnect(ctx context.Context) error {
	return c.cli.Disconnect(ctx)
}

func (c *ClientWrapping) Ping(ctx context.Context, readpref *readpref.ReadPref) error {
	return c.cli.Ping(ctx, readpref)
}

type Database interface {
	Collection(string) Collection
}

type DatabaseWrapping struct {
	db *mongo.Database
}

func (d *DatabaseWrapping) Collection(name string) Collection {
	return &CollectionWrapping{d.db.Collection(name)}
}

type Collection interface {
	Find(context.Context, interface{}) (Cursor, error)
	FindOne(context.Context, interface{}) SingleResult
	InsertOne(context.Context, interface{}) (*mongo.InsertOneResult, error)
}

type CollectionWrapping struct {
	coll *mongo.Collection
}

func (c *CollectionWrapping) Find(ctx context.Context, filter interface{}) (Cursor, error) {
	return c.coll.Find(ctx, filter)
}
func (c *CollectionWrapping) FindOne(ctx context.Context, filter interface{}) SingleResult {
	return c.coll.FindOne(ctx, filter)
}
func (c *CollectionWrapping) InsertOne(ctx context.Context, doc interface{}) (*mongo.InsertOneResult, error) {
	return c.coll.InsertOne(ctx, doc)
}

type Cursor interface {
	Next(context.Context) bool
	Decode(interface{}) error
	Close(context.Context) error
	Err() error
}

type CursorWrapping struct {
	cursor *mongo.Cursor
}

func (c *CursorWrapping) Next(ctx context.Context) bool {
	return c.cursor.Next(ctx)
}

func (c *CursorWrapping) Decode(v interface{}) error {
	return c.cursor.Decode(v)
}

func (c *CursorWrapping) Close(ctx context.Context) error {
	return c.cursor.Close(ctx)
}

func (c *CursorWrapping) Err() error {
	return c.cursor.Err()
}

type SingleResult interface {
	Decode(interface{}) error
}

// Storage handles MongoDB transactions and its connections and contexts
type Storage struct {
	URI       *url.URL
	Client    Client
	TimeoutMS time.Duration
}

// NewStorage returns a new MongoDB Storage
func NewStorage(uri string) *Storage {
	s := &Storage{}

	u, err := url.Parse(uri)

	if err != nil {
		log.Fatalf("invalid MongoDB Connection URI: %v", err)
	}

	s.URI = u

	cli, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("error creating Mongo Client: %v", err)
	}
	client := &ClientWrapping{cli}
	s.Client = client
	return s
}

func (s *Storage) WithTimeout(timeout time.Duration) *Storage {
	s.TimeoutMS = timeout
	return s
}

func (s *Storage) Connect(ctx context.Context) {
	log.Println("connecting to MongoDB")
	ctx, cancelFunc := context.WithTimeout(ctx, s.TimeoutMS)
	defer cancelFunc()
	if err := s.Client.Connect(ctx); err != nil {
		log.Fatalf("error connecting to MongoDB: %s", err.Error())
	}

	if err := s.Client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("server did not respond successfully: %s", err.Error())
	}
}

func (s *Storage) Disconnect(ctx context.Context) {
	log.Println("disconnecting from MongoDB")
	if err := s.Client.Disconnect(ctx); err != nil {
		panic(err)
	}
}
