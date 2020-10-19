package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/renanferr/swapi-golang-rest-api/pkg/adding"
	"github.com/renanferr/swapi-golang-rest-api/pkg/listing"
)

func TestHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/healthcheck", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	a := adding.NewAddingMock("", nil)
	l := listing.NewListingMock([]listing.Planet{}, nil)
	handler := Handler(a, l)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	buf := bytes.NewBuffer([]byte{})
	json.NewEncoder(buf).Encode(map[string]bool{"alive": true})
	expectedBody := buf.String()
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("handler returned unexpected body: got %q want %q", body, expectedBody)
	}
}
