package planets

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/asaskevich/govalidator"
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
		if err != nil {
			errPayload, marshallErr := json.Marshal(&badRequestResponse{"error adding planet", govalidator.ErrorsByField(err)})
			if marshallErr != nil {
				log.Panicf("error marshalling bad request response: %s", marshallErr.Error())
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(errPayload)
			return
		}
		w.Header().Set("Content-Type", "application/json")
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

			if err == listing.ErrNotFound {
				http.Error(w, "The planet you requested does not exist.", http.StatusNotFound)
				return
			}

			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(planet)
	}
}
