package handlers

import (
	"html/template"
	"net/http"
)

// RenderLogin renders the login page
func RenderLogin() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("./templates/login.html", "./templates/base.html")
		if err != nil {
			panic(err)
		}
		t.ExecuteTemplate(w, "base", nil)
	})
}

// HandleLogin accept Google auth login information
func HandleLogin() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: handle login and store user info into session?
	})
}

// HandleLogout invalidates the login session
func HandleLogout() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: handle login and store user info into session?
	})
}
