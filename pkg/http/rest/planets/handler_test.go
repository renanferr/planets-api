package planets

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/asaskevich/govalidator"
	"github.com/renanferr/planets-api/pkg/adding"
	"github.com/renanferr/planets-api/pkg/listing"
	adding_mocks "github.com/renanferr/planets-api/pkg/mocks/adding"
	listing_mocks "github.com/renanferr/planets-api/pkg/mocks/listing"
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
	Name     string
	Adding   *adding_mocks.ServiceMock
	Listing  *listing_mocks.ServiceMock
	Request  TestRequest
	Response ExpectedResponse
}

func marshalBody(t *testing.T, v interface{}) io.Reader {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(v)
	if err != nil {
		t.Fatalf("error marshalling planet: %s", err.Error())
	}
	return buf
	// return bytes.NewReader(b)
}

func runTestCase(t *testing.T, tc *TestCase) {
	rr := httptest.NewRecorder()
	handler := Handler(tc.Adding, tc.Listing)
	req, err := http.NewRequest(tc.Request.Method, tc.Request.Path, tc.Request.Body)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("[%s] %s", req.Method, req.URL)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != tc.Response.Status {
		t.Errorf("<%s> handler returned wrong status code: got %v want %v", tc.Name, status, tc.Response.Status)
	}

	expectedHeaders := fmt.Sprintf("%s", tc.Response.Headers)

	if headers := fmt.Sprintf("%s", rr.HeaderMap); headers != expectedHeaders {
		t.Errorf("<%s> handler returned unexpected header: got %q want %q", tc.Name, headers, expectedHeaders)
	}

	var expectedBody string
	if tc.Response.Body != nil {
		b, err := ioutil.ReadAll(tc.Response.Body)
		if err != nil {
			t.Errorf("<%s> error reading expected response: %w", tc.Name, err)
		}
		expectedBody = string(b)
		fmt.Println(expectedBody)
	}

	if rr.Body.String() != expectedBody {
		t.Errorf("<%s> handler returned unexpected body: got %q want %q", tc.Name, rr.Body.String(), expectedBody)
	}
}

func runTestTable(t *testing.T, tt *[]TestCase) {
	for _, tc := range *tt {
		runTestCase(t, &tc)
	}
}

func TestAddPlanet(t *testing.T) {
	oid := primitive.NewObjectID()

	tt := []TestCase{
		{
			"add planet",
			adding_mocks.NewServiceMock(oid.Hex(), nil),
			listing_mocks.NewServiceMock([]listing.Planet{}, nil, 0),
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
		{
			"JSON decoding error",
			adding_mocks.NewServiceMock("", nil),
			listing_mocks.NewServiceMock(listing.Planet{}, nil, 0),
			TestRequest{
				"POST",
				"/",
				bytes.NewBuffer([]byte{}),
			},
			ExpectedResponse{
				http.StatusBadRequest,
				http.Header{"Content-Type": []string{"application/json"}},
				marshalBody(t, &errorResponse{"JSON decoding error"}),
			},
		},
		{
			"validation error",
			adding_mocks.NewServiceMock("", adding.NewValidationError(govalidator.Error{
				Name:                     "name",
				Err:                      errors.New("Missing required field"),
				CustomErrorMessageExists: false,
				Validator:                "required",
				Path:                     []string{},
			})),
			listing_mocks.NewServiceMock([]listing.Planet{}, nil, 0),
			TestRequest{
				"POST",
				"/",
				marshalBody(t, &adding.Planet{
					Name:        "",
					Climate:     "arid",
					Terrain:     "desert",
					Appearances: 0,
				}),
			},
			ExpectedResponse{
				http.StatusBadRequest,
				http.Header{"Content-Type": []string{"application/json"}},
				marshalBody(t, &validationErrorResponse{"validation error", map[string]string{"name": "Missing required field"}}),
			},
		},
	}

	runTestTable(t, &tt)

}

func TestGetPlanets(t *testing.T) {
	oid, oid2 := primitive.NewObjectID(), primitive.NewObjectID()
	planets := []listing.Planet{
		{
			ID:          oid.Hex(),
			Name:        "tatooine",
			Climate:     "arid",
			Terrain:     "desert",
			Appearances: 5,
		},
		{
			ID:          oid2.Hex(),
			Name:        "alderaan",
			Climate:     "temperate",
			Terrain:     "grasslands",
			Appearances: 2,
		},
	}

	tt := []TestCase{
		{
			"get planets",
			adding_mocks.NewServiceMock(oid.Hex(), nil),
			listing_mocks.NewServiceMock(planets, nil, int64(len(planets))),
			TestRequest{
				"GET",
				"/",
				nil,
			},
			ExpectedResponse{
				http.StatusOK,
				http.Header{
					"Content-Type":  []string{"application/json"},
					"X-Total-Count": []string{strconv.Itoa(len(planets))},
				},
				marshalBody(t, &planets),
			},
		},
		{
			"get 1st planet in 1st page",
			adding_mocks.NewServiceMock(oid.Hex(), nil),
			listing_mocks.NewServiceMock(planets, nil, int64(len(planets))),
			TestRequest{
				"GET",
				"/?page=1&limit=1",
				nil,
			},
			ExpectedResponse{
				http.StatusOK,
				http.Header{
					"Content-Type":  []string{"application/json"},
					"X-Total-Count": []string{strconv.Itoa(len(planets))},
				},
				marshalBody(t, planets[0:1]),
			},
		},
		{
			"get 2nd planet in 2nd page",
			adding_mocks.NewServiceMock(oid.Hex(), nil),
			listing_mocks.NewServiceMock(planets, nil, int64(len(planets))),
			TestRequest{
				"GET",
				"/?page=2&limit=1",
				nil,
			},
			ExpectedResponse{
				http.StatusOK,
				http.Header{
					"Content-Type":  []string{"application/json"},
					"X-Total-Count": []string{strconv.Itoa(len(planets))},
				},
				marshalBody(t, planets[1:2]),
			},
		},
	}

	runTestTable(t, &tt)

}

func TestGetPlanet(t *testing.T) {
	oid := primitive.NewObjectID()
	p := listing.Planet{
		ID:          oid.Hex(),
		Name:        "tatooine",
		Climate:     "arid",
		Terrain:     "desert",
		Appearances: 5,
	}

	tt := []TestCase{
		{
			"get planet",
			adding_mocks.NewServiceMock(oid.Hex(), nil),
			listing_mocks.NewServiceMock(p, nil, 0),
			TestRequest{
				"GET",
				fmt.Sprintf("/%s", oid.Hex()),
				nil,
			},
			ExpectedResponse{
				http.StatusOK,
				http.Header{"Content-Type": []string{"application/json"}},
				marshalBody(t, &p),
			},
		},
		{
			"planet not found",
			adding_mocks.NewServiceMock("", nil),
			listing_mocks.NewServiceMock(listing.Planet{}, listing.ErrPlanetNotFound, 0),
			TestRequest{
				"GET",
				fmt.Sprintf("/%s", primitive.NewObjectID().Hex()),
				nil,
			},
			ExpectedResponse{
				http.StatusNotFound,
				http.Header{},
				nil,
			},
		},
	}

	runTestTable(t, &tt)

}
