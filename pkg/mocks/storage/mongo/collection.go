package mongo

import (
	"context"

	"github.com/renanferr/swapi-golang-rest-api/pkg/mocks"
	"github.com/renanferr/swapi-golang-rest-api/pkg/storage/mongo"
)

type CollectionMock mocks.Mock

func NewCollectionMock(value interface{}, err error) *CollectionMock {
	return &CollectionMock{value, err}
}

func (m *CollectionMock) FindOne(ctx context.Context, filter interface{}) mongo.SingleResult {
	return m.Value.(mongo.SingleResult)
}

func (m *CollectionMock) Find(ctx context.Context, filter interface{}) (mongo.Cursor, error) {
	return m.Value.(mongo.Cursor), m.Err
}

func (m *CollectionMock) InsertOne(ctx context.Context, document interface{}) (interface{}, error) {
	return m.Value, m.Err
}
