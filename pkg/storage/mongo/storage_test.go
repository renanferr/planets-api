package mongo_test

import (
	"testing"
	"time"

	"github.com/renanferr/planets-api/pkg/storage/mongo"
)

func TestStorage(t *testing.T) {
	uri := "mongodb://user:pass@mock.test:1234/mydb?foo=bar"
	s := mongo.NewStorage(uri)

	if s.URI.String() != uri {
		t.Errorf("Connection URIs do not match. got: %s want: %s", s.URI.String(), uri)
	}

	ms := 5000
	timeout := time.Duration(ms) * time.Millisecond
	s = s.WithTimeout(timeout)

	if s.TimeoutMS != timeout {
		t.Errorf("timeouts do not match. got: %s want: %s", s.TimeoutMS, timeout)
	}
}
