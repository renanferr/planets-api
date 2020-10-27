package rest

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/renanferr/planets-api/pkg/adding"
	"github.com/renanferr/planets-api/pkg/http/rest/planets"
	"github.com/renanferr/planets-api/pkg/listing"
)

const ApiPrefix = "/api"

func Handler(a adding.Service, l listing.Service) http.Handler {
	apiRouter := chi.NewRouter()

	apiRouter.Use(middleware.Recoverer, middleware.RequestID, middleware.Logger)
	apiRouter.Mount("/planets", planets.Handler(a, l))

	router := chi.NewRouter()

	router.Get("/healthcheck", handleHealthcheck)

	router.Mount(ApiPrefix, apiRouter)

	return router
}

func handleHealthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"alive": true})
}
