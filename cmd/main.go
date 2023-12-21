package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	ai "hackathon-2023-vendor-reviews"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Route("/api/v1/genai", func(router chi.Router) {
		router.Get("/summary", ai.GetReviewsSummary)
		router.Get("/reviews", ai.GetReviews)
	})
	err := http.ListenAndServe(":3030", r)
	if err != nil {
		log.Fatal("error in starting")
	}
}
