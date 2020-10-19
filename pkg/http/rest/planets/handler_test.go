package planets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/renanferr/swapi-golang-rest-api/pkg/adding"
	"github.com/renanferr/swapi-golang-rest-api/pkg/listing"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestRequest struct {
	Method string
	Path   string
	Body   io.Reader
}

type ExpectedResponse struct {
	Status  int
	Headers http.Header
	Body    io.Reader
}

type TestCase struct {
	request  TestRequest
	response ExpectedResponse
}

func marshalBody(t *testing.T, v interface{}) io.Reader {
	b, err := json.Marshal(v)
	if err != nil {
		t.Errorf("error marshalling planet: %w", err)
	}

	return bytes.NewReader(b)
}

func TestAddPlanet(t *testing.T) {
	rr := httptest.NewRecorder()

	oid := primitive.NewObjectID()
	a := adding.NewAddingMock(oid.Hex(), nil)
	l := listing.NewListingMock([]listing.Planet{}, nil)
	handler := Handler(a, l)

	tt := []TestCase{
		{
			TestRequest{
				"POST",
				"/",
				marshalBody(t, &adding.Planet{
					Name:        "tatooine",
					Climate:     "arid",
					Terrain:     "desert",
					Appearances: 0,
				}),
			},
			ExpectedResponse{
				http.StatusCreated,
				http.Header{"Location": []string{fmt.Sprintf("/%s", oid.Hex())}},
				nil,
			},
		},
	}

	for _, tc := range tt {
		req, err := http.NewRequest(tc.request.Method, tc.request.Path, tc.request.Body)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("[%s] %s", req.Method, req.URL)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != tc.response.Status {
			t.Errorf("handler returned wrong status code: got %v want %v", status, tc.response.Status)
		}

		expectedHeaders := fmt.Sprintf("%s", tc.response.Headers)

		if headers := fmt.Sprintf("%s", rr.HeaderMap); headers != expectedHeaders {
			t.Errorf("handler returned unexpected header: got %s want %s", headers, expectedHeaders)
		}

		var expectedBody string
		if tc.response.Body != nil {
			b, err := ioutil.ReadAll(tc.response.Body)
			if err != nil {
				t.Errorf("error reading expected response: %w", err)
			}
			expectedBody = string(b)
		}

		if rr.Body.String() != expectedBody {
			t.Errorf("handler returned unexpected body: got %q want %q", rr.Body.String(), expectedBody)
		}
	}

}
