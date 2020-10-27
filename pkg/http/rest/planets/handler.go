package planets

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/renanferr/planets-api/pkg/adding"
	"github.com/renanferr/planets-api/pkg/listing"
)

// Router returns the /planets router
func Handler(a adding.Service, l listing.Service) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/", addPlanet(a))
	r.Get("/", getPlanets(l))
	r.Get("/{planetID}", getPlanet(l))

	return r
}

type validationErrorResponse struct {
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields"`
}

type errorResponse struct {
	Message string `json:"message"`
}

func sendErrorResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(&errorResponse{message}); err != nil {
		log.Panicf("error encoding error response: %s", err)
	}
}

// addPlanet returns a handler for POST /planets requests
func addPlanet(s adding.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var newPlanet adding.Planet
		err := decoder.Decode(&newPlanet)
		if err != nil {
			log.Printf("error decoding planet: %s", err.Error())
			sendErrorResponse(w, http.StatusBadRequest, "JSON decoding error")
			return
		}

		id, err := s.AddPlanet(r.Context(), newPlanet)

		var e *adding.ValidationError
		if errors.As(err, &e) {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			err := json.NewEncoder(w).Encode(&validationErrorResponse{"validation error", e.Fields})
			if err != nil {
				log.Panicf("error encoding bad request response: %s", err.Error())
			}
			return
		}

		if err != nil {
			log.Panicf("error adding planet: %s", err)
		}

		w.Header().Set("Location", fmt.Sprintf("/%s", id))
		w.WriteHeader(http.StatusCreated)
	}
}

func getOffset(limit int64, page int64) int64 {
	if page < 1 {
		return 0
	}

	return (page - 1) * limit
}

func getPaginationInfo(query url.Values) (int64, int64, error) {

	info := map[string]int64{"limit": 20, "page": 1}

	var err error

	for k := range info {
		v := query.Get(k)
		if v != "" {
			val, err := strconv.ParseInt(v, 10, 64)
			info[k] = val
			if err != nil {
				log.Println(err)
				return 0, 0, fmt.Errorf("%s value must be an integer. got: %s", k, v)
			}
		}
	}

	return info["limit"], info["page"], err
}

// getPlanets returns a handler for GET /planets requests
func getPlanets(s listing.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		limit, page, err := getPaginationInfo(r.URL.Query())
		if err != nil {
			log.Println(err)
			sendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if limit < 1 || page < 1 {
			err = errors.New("pagination out of range")
			log.Println(err)
			sendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		offset := getOffset(limit, page)

		list, total := s.GetPlanets(r.Context(), offset, limit)
		w.Header().Set("X-Total-Count", strconv.FormatInt(total, 10))
		json.NewEncoder(w).Encode(list)
	}
}

// getPlanet returns a handler for GET /planets/:id requests
func getPlanet(s listing.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ID := chi.URLParam(r, "planetID")

		planet, err := s.GetPlanet(r.Context(), ID)
		if err != nil {

			if errors.Is(err, listing.ErrPlanetNotFound) {
				log.Printf("Planet with ID %s not found", ID)
				w.WriteHeader(http.StatusNotFound)
				return
			}

			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(planet)
	}
}
