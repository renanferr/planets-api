package rest

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/renanferr/swapi-golang-rest-api/pkg/adding"
	"github.com/renanferr/swapi-golang-rest-api/pkg/http/rest/middleware"
	"github.com/renanferr/swapi-golang-rest-api/pkg/http/rest/planets"
	"github.com/renanferr/swapi-golang-rest-api/pkg/listing"
)

func Handler(a adding.Service, l listing.Service) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer, middleware.RequestID, middleware.Logger)

	router.Mount("/planets", planets.Handler(a, l))

	return router
}
