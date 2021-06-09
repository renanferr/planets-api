package mongo

import "github.com/renanferr/planets-api/pkg/storage/mongo"

type DatabaseMock struct {
	Coll mongo.Collection
}

func NewDatabaseMock(coll mongo.Collection) mongo.Database {
	return &DatabaseMock{Coll: coll}
}

func (m *DatabaseMock) Collection(name string) mongo.Collection {
	return m.Coll
}
