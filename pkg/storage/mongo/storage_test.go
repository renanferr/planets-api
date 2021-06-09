package mongo_test

import (
	"testing"
	"time"

	"github.com/renanferr/planets-api/pkg/storage/mongo"
)

func TestStorage(t *testing.T) {
	tt := []struct {
		uri                  string
		timeout              time.Duration
		expectedDatabaseName string
	}{
		{
			"mongodb://user:pass@mock.test:1234/mydb?foo=bar",
			time.Duration(5000) * time.Millisecond,
			"mydb",
		},
	}
	for _, tc := range tt {
		s := mongo.NewStorage(tc.uri)

		if s.URI.String() != tc.uri {
			t.Errorf("Connection URIs do not match. got: %s want: %s", s.URI.String(), tc.uri)
		}

		s = s.WithTimeout(tc.timeout)

		if s.TimeoutMS != tc.timeout {
			t.Errorf("timeouts do not match. got: %s want: %s", s.TimeoutMS, tc.timeout)
		}

		if s.DatabaseName() != tc.expectedDatabaseName {
			t.Errorf("unexpected database name. god: %s want :%s", s.DatabaseName(), tc.expectedDatabaseName)
		}
	}
}
