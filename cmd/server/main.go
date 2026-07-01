package main

import (
	"cinema_booking/db"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	db.Init()
	defer db.DB.Close()

	r := chi.NewRouter()

	r.Get("/api/v1")

	http.ListenAndServe(":8000", r)

}
