// Package handlers defines route handlers
package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/golang-jwt/jwt/v5"
)

const admin = "wesley"

// Home handles login status for rendering templates
func Home(w http.ResponseWriter, r *http.Request, tokenAuth *jwtauth.JWTAuth, tmpl *template.Template) {
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

// Logout clears the cookie
func Logout(w http.ResponseWriter, _ *http.Request) {
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

// Login does JWT auth
func Login(w http.ResponseWriter, r *http.Request, users map[string]string, tokenAuth *jwtauth.JWTAuth) {
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
