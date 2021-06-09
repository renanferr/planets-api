package mongo

import (
	"context"

	"github.com/renanferr/planets-api/pkg/storage/mongo"
	mongodb "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CollectionMock struct {
	Value interface{}
	Err   error
}

func NewCollectionMock(value interface{}, err error) *CollectionMock {
	return &CollectionMock{value, err}
}

func (m *CollectionMock) Count(ctx context.Context, filter interface{}) (int64, error) {
	return m.Value.(*CursorMock).Len(), m.Err
}

func (m *CollectionMock) FindOne(ctx context.Context, filter interface{}) mongo.SingleResult {
	return m.Value.(mongo.SingleResult)
}

func (m *CollectionMock) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (mongo.Cursor, error) {
	return m.Value.(*CursorMock), m.Err
}

func (m *CollectionMock) InsertOne(ctx context.Context, document interface{}) (*mongodb.InsertOneResult, error) {
	return m.Value.(*mongodb.InsertOneResult), m.Err
}
