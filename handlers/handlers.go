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
	Schools []*web.School
}

// RenderHomepage renders the homepage for bentobbox with a list of schools
func RenderHomepage(dao web.SchoolDAO) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		schools, err := dao.GetSchools()
		if err != nil {
			log.Println("Failed to grab schools", err)
			schools = []*web.School{}
		}
		t, err := template.ParseFiles("./templates/index.html", "./templates/base.html")
		if err != nil {
			panic(err)
		}
		t.ExecuteTemplate(w, "base", homepageDTO{
			Schools: schools,
		})
	})
}

type schoolListingDTO struct {
	Schools []*web.School
}

// RenderSchools renders a list of schools to give user general idea on what
// schools needs donations
func RenderSchools(dao web.SchoolDAO) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		schools, err := dao.GetSchools()
		if err != nil {
			log.Println("Failed to grab schools", err)
			schools = []*web.School{}
		}
		t, err := template.ParseFiles("./templates/schools.html", "./templates/base.html")
		if err != nil {
			panic(err)
		}
		t.ExecuteTemplate(w, "base", schoolListingDTO{
			Schools: schools,
		})
	})
}

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

// RenderCreateSchool renders the login page
func RenderCreateSchool() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("./templates/create.html", "./templates/base.html")
		if err != nil {
			panic(err)
		}
		t.ExecuteTemplate(w, "base", nil)
	})
}

// RenderSchool renders individual school detail
func RenderSchool(dao web.SchoolDAO) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("./templates/school.html", "./templates/base.html")
		if err != nil {
			panic(err)
		}
		t.ExecuteTemplate(w, "base", nil)
	})
}

// RenderEditSchool renders the login page
func RenderEditSchool(dao web.SchoolDAO) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("./templates/edit.html", "./temnplates/base.html")
		if err != nil {
			panic(err)
		}
		t.ExecuteTemplate(w, "base", nil)
	})
}

// HandleCSVUpload parses CSV and return JSON of the student loan detail
func HandleCSVUpload() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: parse CSV from body and return JSON representation of CSV
	})
}

// HandleCreateSchool creates school from the request form body
func HandleCreateSchool(dao web.SchoolDAO) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// parse form
		// call dao.CreateSchool
		// return 200 or 500
	})
}

// HandleEditSchool edits the school by its id
func HandleEditSchool(dao web.SchoolDAO) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// parse form and id
		// call dao.EditSchool
		// return 200 or 500
	})
}
