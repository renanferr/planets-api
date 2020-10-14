package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/renanferr/swapi-golang-rest-api/pkg/adding"
	"github.com/renanferr/swapi-golang-rest-api/pkg/fetching"
	"github.com/renanferr/swapi-golang-rest-api/pkg/http/rest"
	"github.com/renanferr/swapi-golang-rest-api/pkg/listing"
	"github.com/renanferr/swapi-golang-rest-api/pkg/storage"
)

func main() {

	timeoutStr := os.Getenv("DB_TIMEOUT_MS")
	timeoutInt, err := strconv.Atoi(timeoutStr)
	if err != nil {
		log.Fatalf("couldn't use \"%s\" as database timeout", timeoutStr)
	}
	timeoutMS := time.Duration(timeoutInt) * time.Millisecond

	s := storage.NewStorage(os.Getenv("DB_CONNECTION_URI")).WithTimeout(timeoutMS)

	ctx := context.Background()
	s.Connect(ctx)
	defer s.Disconnect(ctx)

	fetcher, err := fetching.NewClient(os.Getenv("PLANETS_API_BASE_URL"))
	if err != nil {
		log.Fatalf("Error creating fetching client")
	}
	adder := adding.NewService(s, fetcher)
	lister := listing.NewService(s)

	router := rest.Handler(adder, lister)

	port := os.Getenv("APP_PORT")
	if !strings.HasPrefix(port, ":") {
		port = fmt.Sprintf(":%s", port)
	}
	log.Printf("Starting server at %s", port)
	log.Fatal(http.ListenAndServe(port, router))
}
