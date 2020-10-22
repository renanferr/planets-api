package mongo

import (
	"context"
	"errors"
	"log"
	"reflect"

	"gopkg.in/mgo.v2/bson"
)

type CursorMock struct {
	Cursor  int
	Current interface{}
	Values  interface{}
	Error   error
}

func NewCursorMock(values interface{}, err error) *CursorMock {
	return &CursorMock{
		Cursor:  0,
		Current: nil,
		Values:  values,
		Error:   err,
	}
}

func (m *CursorMock) Next(ctx context.Context) bool {
	hasNext := m.Cursor < reflect.ValueOf(m.Values).Len()
	if hasNext {
		var err error
		m.Current = reflect.ValueOf(m.Values).Index(m.Cursor).Interface()
		log.Printf("current %v %s", m.Current, reflect.TypeOf(m.Current))
		if err != nil {
			panic(err)
		}
		m.Cursor++
	}
	return hasNext
}

func (m *CursorMock) Decode(v interface{}) error {
	// if m.Error
	if m.Current != nil {
		b, err := bson.Marshal(m.Current)

		if err != nil {
			log.Printf("encoding error: %v", err.Error())
			return err
		}
		return bson.Unmarshal(b, v)
	}

	return errors.New("current value is `nil`")
}

func (m *CursorMock) Err() error {
	return m.Error
}

func (m *CursorMock) Close(ctx context.Context) {}
