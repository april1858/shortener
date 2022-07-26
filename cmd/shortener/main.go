package main

import (
	"log"
	"net/http"

	"github.com/april1858/shortener/internal/app"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	r.Post("/", app.CreateShort)
	r.Get("/{id}", app.ReturnLong)

	log.Fatal(server.ListenAndServe())
}
