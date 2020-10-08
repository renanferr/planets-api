package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/renanferr/swapi-golang-rest-api/pkg/adding"
	"github.com/renanferr/swapi-golang-rest-api/pkg/http/rest"
	"github.com/renanferr/swapi-golang-rest-api/pkg/listing"
	"github.com/renanferr/swapi-golang-rest-api/pkg/storage"
)

func main() {

	s, err := storage.NewStorage()
	if err != nil {
		log.Fatalf("Error creating storage: %v", err)
	}

	adder := adding.NewService(s)
	lister := listing.NewService(s)

	// set up the HTTP server
	router := rest.Handler(adder, lister)

	fmt.Println("The planet server is on tap now: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
