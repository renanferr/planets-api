package planets

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/renanferr/swapi-golang-rest-api/pkg/adding"
	"github.com/renanferr/swapi-golang-rest-api/pkg/listing"
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

// addPlanet returns a handler for POST /planets requests
func addPlanet(s adding.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var newPlanet adding.Planet
		err := decoder.Decode(&newPlanet)
		if err != nil {
			log.Printf("error decoding planet: %s", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-type", "application/json")
			if err = json.NewEncoder(w).Encode(&errorResponse{"decoding error"}); err != nil {
				log.Panicf("error encoding error response: %s", err)
			}
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

// getPlanets returns a handler for GET /planets requests
func getPlanets(s listing.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		list := s.GetPlanets(r.Context())
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
