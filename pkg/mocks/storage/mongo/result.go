package mongo

import (
	"github.com/renanferr/swapi-golang-rest-api/pkg/mocks"
	"gopkg.in/mgo.v2/bson"
)

type SingleResultMock mocks.Mock

func (m *SingleResultMock) Decode(v interface{}) error {
	if m.Err != nil {
		return m.Err
	}
	return bson.Unmarshal(m.Value.([]byte), v)
}
