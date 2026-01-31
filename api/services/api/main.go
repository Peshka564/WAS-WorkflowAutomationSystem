package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
    r := chi.NewRouter()
    r.Use(middleware.Logger)
    r.Use(cors.Handler(cors.Options{
        AllowedOrigins: []string{"http://*"},
        AllowedMethods: []string{"GET"},
    }))
    r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		  w.Write([]byte("Hello, Dad!"))
    })
    http.ListenAndServe(":3000", r)
}