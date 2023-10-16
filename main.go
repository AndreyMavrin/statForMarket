package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"statForMarket/internal/api"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	ctx := context.Background()

	a := api.Application{}

	a.Run(ctx)

	r := mux.NewRouter()

	r.HandleFunc("/api/event", a.TestEvents).Methods(http.MethodPost)
	r.HandleFunc("/api/events", a.Events).Methods(http.MethodGet)

	r.HandleFunc("/api/event", a.CreateEvent).Methods(http.MethodPost)

	port := os.Getenv("PORT")
	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Starting server at: %s", port)
	log.Fatal(srv.ListenAndServe())
}
