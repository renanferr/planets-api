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
	r.Get("/:id", getPlanet(l))

	return r
}

type badRequestResponse struct {
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields"`
}

// addPlanet returns a handler for POST /planets requests
func addPlanet(s adding.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var newPlanet adding.Planet
		err := decoder.Decode(&newPlanet)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := s.AddPlanet(r.Context(), newPlanet)

		var e *adding.ValidationError
		if errors.As(err, &e) {
			response, err := json.Marshal(&badRequestResponse{"error adding planet", e.Fields})
			if err != nil {
				log.Panicf("error marshalling bad request response: %s", err.Error())
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
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
		ID := chi.URLParam(r, "id")

		planet, err := s.GetPlanet(r.Context(), ID)
		if err != nil {

			if errors.Is(err, listing.ErrPlanetNotFound) {
				http.Error(w, "", http.StatusNotFound)
				return
			}

			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(planet)
	}
}
