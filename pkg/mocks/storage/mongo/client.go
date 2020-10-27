package mongo

import (
	"context"

	"github.com/renanferr/swapi-golang-rest-api/pkg/storage/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type ClientMock struct {
	DB           mongo.Database
	ConnectErr   error
	DisconectErr error
	PingErr      error
}

func NewClientMock(db mongo.Database, connectErr error, disconnectErr error, pingErr error) *ClientMock {
	return &ClientMock{db, connectErr, disconnectErr, pingErr}
}

func (m *ClientMock) Database(databaseName string) mongo.Database {
	return m.DB
}

func (m *ClientMock) Connect(ctx context.Context) error {
	return m.ConnectErr
}

func (m *ClientMock) Disconnect(ctx context.Context) error {
	return m.DisconectErr
}

func (m *ClientMock) Ping(ctx context.Context, readpref *readpref.ReadPref) error {
	return m.PingErr
}
