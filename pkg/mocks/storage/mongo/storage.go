package mongo

import (
	"github.com/renanferr/swapi-golang-rest-api/pkg/storage/mongo"
)

type StorageMock mongo.Storage

func NewStorageMock(uri string, client mongo.Client) *mongo.Storage {
	s := mongo.NewStorage(uri)
	s.Client = client
	return s
}

func (m *StorageMock) WithDatabase(database *DatabaseMock) *StorageMock {
	m.Client.(*ClientMock).DB = database
	return m
}

func (m *StorageMock) WithCollection(collection *CollectionMock) *StorageMock {
	m.Client.(*ClientMock).DB.Coll = collection
	return m
}

func (m *StorageMock) WithSingleResult(result *SingleResultMock) *StorageMock {
	m.Client.(*ClientMock).DB.Coll.Value = result
	return m
}
