package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/BentoBoxSchool/web"
	"github.com/pkg/errors"
)

// Hello says hello
func Hello() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, BentoBox Server!")
	})
}

// CheckHealth returns the healthcheck for the critical resources
func CheckHealth(sqlDB *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := sqlDB.Ping(); err != nil {
			http.Error(
				w,
				errors.Wrap(err, "MySQL failed").Error(),
				http.StatusInternalServerError,
			)
			return
		}

		fmt.Fprintln(w, "Healthy")
	})
}

// Assets serves the static assets (js & css)
func Assets() http.Handler {
	return http.StripPrefix("/assets", http.FileServer(http.Dir("./assets")))
}

type homepageDTO struct {
	Schools []web.School
}

// RenderHomepage renders the homepage for bentobbox with a list of schools
func RenderHomepage(dao web.SchoolDAO) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err, schools := dao.GetSchools()
		if err != nil {
			log.Println("Failed to grab schools", err)
			schools = []web.School{}
		}
		t, err := template.ParseFiles("./templates/index.html")
		if err != nil {
			panic(err)
		}
		t.Execute(w, homepageDTO{
			Schools: schools,
		})
	})
}
