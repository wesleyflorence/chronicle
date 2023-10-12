//Package routes sets up routes
package routes

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"github.com/jomei/notionapi"
	"github.com/wesleyflorence/chronicle/handlers"
)

// SetupRoutes defines the endpoints and functions that are called
func SetupRoutes(r *chi.Mux, tmpl *template.Template) {
	notionAPIKey := os.Getenv("NOTION_API_KEY")
	digestionDbID := os.Getenv("DIGESTION_DB")
	medicinePageID := os.Getenv("MEDICINE_PAGE")
	username := os.Getenv("CHRONICLE_USERNAME")
	password := os.Getenv("CHRONICLE_PW")
	jwtSecret := os.Getenv("JWT_SECRET")
	tokenAuth := jwtauth.New("HS256", []byte(jwtSecret), nil)
	users := map[string]string{
		username: password,
	}
	client := notionapi.NewClient(notionapi.Token(notionAPIKey))

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "public"))
	handlers.FileServer(r, "/", filesDir)

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			handlers.Home(w, r, tokenAuth, tmpl)
		})
		r.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
			handlers.Logout(w, r)
		})
		r.Post("/api/v1/login", func(w http.ResponseWriter, r *http.Request) {
			handlers.Login(w, r, users, tokenAuth)
		})
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/api/v1/dig", func(w http.ResponseWriter, r *http.Request) {
			handlers.DigestionEntry(w, r, client, digestionDbID)
		})
		r.Post("/api/v1/med", func(w http.ResponseWriter, r *http.Request) {
			handlers.MedicineEntry(w, r, client, medicinePageID)
		})
	})
}
