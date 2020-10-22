package mongo

import "github.com/renanferr/swapi-golang-rest-api/pkg/storage/mongo"

type DatabaseMock struct {
	Coll *CollectionMock
}

func (m *DatabaseMock) Collection(name string) mongo.Collection {
	return m.Coll
}
