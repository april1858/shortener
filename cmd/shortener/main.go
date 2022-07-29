package main

import (
	"log"
	"net/http"
	"os"

	"github.com/april1858/shortener/internal/app"
	"github.com/go-chi/chi/v5"
)

func main() {
	addr := os.Getenv("SERVER_ADDRESS")
	if addr == "" {
		addr = ":8080"
	}
	bURL := os.Getenv("BASE_URL")
	if bURL != "" {
		bURL = ""
	}
	if fsp := os.Getenv("FILE_STORAGE_PATH"); fsp != "" {
		app.WriteInDB(fsp)
	}
	r := chi.NewRouter()
	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	r.Post("/", app.CreateShort)
	r.Get("/{id}", app.ReturnLong)
	r.Post("/api/shorten", app.APIShorten)

	log.Fatal(server.ListenAndServe())
}
