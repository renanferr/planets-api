package mongo

import (
	"github.com/renanferr/planets-api/pkg/mocks"
	"go.mongodb.org/mongo-driver/bson"
)

type SingleResultMock mocks.Mock

func (m *SingleResultMock) Decode(v interface{}) error {
	if m.Err != nil {
		return m.Err
	}
	return bson.Unmarshal(m.Value.([]byte), v)
}
