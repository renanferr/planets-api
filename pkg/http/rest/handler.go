package rest

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/renanferr/swapi-golang-rest-api/pkg/adding"
	"github.com/renanferr/swapi-golang-rest-api/pkg/http/rest/planets"
	"github.com/renanferr/swapi-golang-rest-api/pkg/listing"
)

const docsDirRelativePath = "../../docs"

func Handler(a adding.Service, l listing.Service) http.Handler {
	apiRouter := chi.NewRouter()

	apiRouter.Use(middleware.Recoverer, middleware.RequestID, middleware.Logger)
	apiRouter.Mount("/planets", planets.Handler(a, l))

	router := chi.NewRouter()
	router.Mount("/api", apiRouter)

	return router
}
