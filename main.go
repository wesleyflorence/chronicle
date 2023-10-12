// Serve Chronicle Fiber App
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/jomei/notionapi"
	"github.com/wesleyflorence/chronicle/notion"
)

const admin = "wesley"

var tokenAuth *jwtauth.JWTAuth
var tmpl *template.Template

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

	var err error
	tmpl, err = parseTemplates("views")
	if err != nil {
		log.Fatal(err)
	}

	setupRoutes(r)
	log.Fatal(http.ListenAndServe(":8080", r))
}

func parseTemplates(folder string) (*template.Template, error) {
	return template.ParseGlob(filepath.Join(folder, "*.html"))
}

func setupUsers() map[string]string {
	username := os.Getenv("CHRONICLE_USERNAME")
	password := os.Getenv("CHRONICLE_PW")
	users := map[string]string{
		username: password,
	}
	return users
}

func setupRoutes(r *chi.Mux) {
	notionAPIKey := os.Getenv("NOTION_API_KEY")
	digestionDbID := os.Getenv("DIGESTION_DB")
	medicinePageID := os.Getenv("MEDICINE_PAGE")
	jwtSecret := os.Getenv("JWT_SECRET")
	tokenAuth = jwtauth.New("HS256", []byte(jwtSecret), nil)
	users := setupUsers()
	client := notionapi.NewClient(notionapi.Token(notionAPIKey))

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "public"))
	fileServer(r, "/", filesDir)

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/", handleHome)
		r.Get("/logout", handleLogout)
		r.Post("/api/v1/login", func(w http.ResponseWriter, r *http.Request) {
			handleLogin(w, r, users)})
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/api/v1/dig", func(w http.ResponseWriter, r *http.Request) {
			handleDigestionEntry(w, r, client, digestionDbID)
		})
		r.Post("/api/v1/med", func(w http.ResponseWriter, r *http.Request) {
			handleMedicineEntry(w, r, client, medicinePageID)
		})
	})
}

// Serve Static Files
func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	isLoggedIn := false
	cookie, _ := r.Cookie("jwt")
	if cookie != nil {
		token, err := tokenAuth.Decode(cookie.Value)
		if err == nil {
			name, _ := token.Get("name")
			if name == admin {
				isLoggedIn = true
			}
		}
	}

	data := map[string]interface{}{
		"Title":        "Chronicle",
		"IsLoggedIn":   isLoggedIn,
		"BristolSlice": []int{1, 2, 3, 4, 5, 6, 7},
	}

	err := tmpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		log.Printf("Error rendering template: %+v", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func handleLogout(w http.ResponseWriter, _ *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})
	w.Header().Set("HX-Redirect", "/")
	w.Write([]byte("Logged out"))
}

func handleLogin(w http.ResponseWriter, r *http.Request, users map[string]string) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	formPassword := r.Form.Get("Password")
	password, ok := users[admin]
	if !ok || password != formPassword {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `<div class="text-red-500">Login failed. Please try again.</div>`)
		return
	}

	_, tokenString, err := tokenAuth.Encode(jwt.MapClaims{
		"name":  admin,
		"admin": true,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})
	if err != nil {
		log.Printf("Error encoding token: %v", err)
		http.Error(w, "Error encoding token", http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  time.Now().Add(72 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)

	w.Header().Set("HX-Redirect", "/")
	fmt.Fprint(w, "Logged in!")
}

func handleMedicineEntry(w http.ResponseWriter, r *http.Request, client *notionapi.Client, medicinePageID string) {
	type Payload struct {
		Medicine string
		Note     string
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing medicine entry: %v", err)
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}
	payload := Payload{
		Medicine: r.Form.Get("medicine"),
		Note:     r.Form.Get("note"),
	}

	page, err := notion.AppendMedicineEntry(client, medicinePageID, payload.Medicine, payload.Note)
	if err != nil {
		log.Printf("Error appending medicine entry: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	loc, _ := time.LoadLocation("Local")
	created := page.CreatedTime.In(loc).Format("2006-01-02 03:04PM")
	doseProp, ok := page.Properties["Dose"].(*notionapi.TitleProperty)
	if !ok {
		log.Printf("Error unwrapping Dose returned from page: %v", err)
		http.Error(w, "Error unwrapping Dose returned from page", http.StatusInternalServerError)
		return
	}
	dose := doseProp.Title[0].Text.Content
	response := fmt.Sprintf(`<div id="med-response-target" class="text-xs text-stone-600" hx-ext="remove-me"><div remove-me="5s">%s dose %s :: %s</div></div>`, payload.Medicine, dose, created)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(response))
}

func handleDigestionEntry(w http.ResponseWriter, r *http.Request, client *notionapi.Client, digestionDbID string) {
	type Payload struct {
		Bristol int
		Size    string
		Note    string
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing digestion entry: %v", err)
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}
	bristol, err := strconv.Atoi(r.Form.Get("bristol"))
	if err != nil {
		log.Printf("Error parsing bristol value: %v", err)
		http.Error(w, "Invalid value for Bristol", http.StatusBadRequest)
		return
	}
	payload := Payload{
		Bristol: bristol,
		Size:    r.Form.Get("size"),
		Note:    r.Form.Get("note"),
	}

	page, err := notion.AppendDigestionEntry(client, digestionDbID, payload.Bristol, payload.Size, payload.Note)
	if err != nil {
		log.Printf("Error appending digestion entry: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	loc, _ := time.LoadLocation("Local")
	created := page.CreatedTime.In(loc).Format("2006-01-02 03:04PM")
	response := fmt.Sprintf(`<div id="dig-response-target" class="text-xs text-stone-600" hx-ext="remove-me"><div remove-me="5s">New Entry :: %s</div></div>`, created)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(response))
}
