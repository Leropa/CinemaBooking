package main

import (
	"cinema_booking/db"
	"cinema_booking/internal/handlers"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	db.Init()
	defer db.DB.Close()

	r := chi.NewRouter()

	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir("public"))))
	r.Get("/api/v1/movies", handlers.GetMovies)
	r.Get("/api/v1/shows", handlers.GetShows)
	r.Get("/api/v1/seats", handlers.GetSeats)

	log.Fatal(http.ListenAndServe(":8000", r))

}
