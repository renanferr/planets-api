package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/renanferr/planets-api/pkg/adding"
	"github.com/renanferr/planets-api/pkg/http/client"
	"github.com/renanferr/planets-api/pkg/http/rest"
	"github.com/renanferr/planets-api/pkg/listing"
	"github.com/renanferr/planets-api/pkg/storage/mongo"
)

func main() {

	timeoutStr := os.Getenv("DB_TIMEOUT_MS")
	timeoutInt, err := strconv.Atoi(timeoutStr)
	if err != nil {
		log.Fatalf("couldn't use \"%s\" as database timeout: %s", timeoutStr, err.Error())
	}
	timeoutMS := time.Duration(timeoutInt) * time.Millisecond

	s := mongo.NewStorage(os.Getenv("DB_CONNECTION_URI")).WithTimeout(timeoutMS)

	ctx := context.Background()
	s.Connect(ctx)
	defer s.Disconnect(ctx)

	planetsClient, err := client.NewClient(os.Getenv("PLANETS_API_BASE_URL"))
	if err != nil {
		log.Fatalf("error creating fetching client: %s", err.Error())
	}
	adder := adding.NewService(s, planetsClient)
	lister := listing.NewService(s)

	router := rest.Handler(adder, lister)

	port := os.Getenv("PORT")
	if !strings.HasPrefix(port, ":") {
		port = fmt.Sprintf(":%s", port)
	}

	go func() {
		log.Printf("starting server at %s", port)
		log.Fatal(http.ListenAndServe(port, router))
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch
}
