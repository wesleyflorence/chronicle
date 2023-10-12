// Serve Chronicle Fiber App
package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/wesleyflorence/chronicle/routes"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	tmpl, err := template.ParseGlob(filepath.Join("views", "*.html"))
	if err != nil {
		log.Fatal(err)
	}

	routes.SetupRoutes(r, tmpl)
	log.Fatal(http.ListenAndServe(":8080", r))
}
