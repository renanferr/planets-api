package mongo

import (
	"context"
	"errors"
	"log"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
)

type CursorMock struct {
	Cursor  int
	Current interface{}
	Values  interface{}
	Error   error
	len     int64
}

func NewCursorMock(values interface{}, err error) *CursorMock {
	return &CursorMock{
		Cursor:  0,
		Current: nil,
		Values:  values,
		Error:   err,
		len:     int64(reflect.ValueOf(values).Len()),
	}
}

func (m *CursorMock) Len() int64 {
	return m.len
}

func (m *CursorMock) Next(ctx context.Context) bool {
	hasNext := m.Cursor < reflect.ValueOf(m.Values).Len()
	if hasNext {
		var err error
		m.Current = reflect.ValueOf(m.Values).Index(m.Cursor).Interface()
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

func (m *CursorMock) Close(ctx context.Context) error {
	return nil
}
