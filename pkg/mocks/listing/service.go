package listing

import (
	"context"
	"errors"
	"log"

	"github.com/renanferr/swapi-golang-rest-api/pkg/listing"
)

type ServiceMock struct {
	Value interface{}
	Err   error
	Total int64
}

func NewServiceMock(v interface{}, err error, total int64) *ServiceMock {
	return &ServiceMock{v, err, total}
}

func (m *ServiceMock) GetPlanet(ctx context.Context, id string) (listing.Planet, error) {
	v, ok := m.Value.(listing.Planet)
	if !ok {
		return listing.Planet{}, errors.New("unexpected value type. expected: listing.Planet")
	}

	return v, m.Err
}

func (m *ServiceMock) GetPlanets(ctx context.Context, offset int64, limit int64) ([]listing.Planet, int64) {
	v, ok := m.Value.([]listing.Planet)
	if !ok {
		log.Panicf("could not assert %v of type `listing.Planet`", m.Value)
	}

	log.Printf("%d %d", offset, limit)
	if limit > m.Total {
		limit = m.Total
	}

	log.Printf("v[%d:%d]", offset, limit)

	return v[offset : offset+limit], m.Total
}
