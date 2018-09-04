package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/BentoBoxSchool/web"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
)

// Assets serves the static assets (js & css)
func Assets() http.Handler {
	return http.StripPrefix("/assets", http.FileServer(http.Dir("./assets")))
}

type homepageDTO struct {
	Schools []*web.School
	User    *web.User
}

// FIXME: probably should create a UserDAO
func getUser(store sessions.Store, r *http.Request) *web.User {
	user := &web.User{}
	session, err := store.Get(r, "user")
	if err != nil {
		log.Println("failed to get session from request", err)
	} else {
		userValue := session.Values["user"]
		if userValue != nil {
			user = session.Values["user"].(*web.User)
		}
	}
	return user
}

// RenderHomepage renders the homepage for bentobbox with a list of schools
func RenderHomepage(store sessions.Store, dao web.SchoolDAO) http.HandlerFunc {
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
		user := getUser(store, r)
		t.ExecuteTemplate(w, "base", homepageDTO{
			Schools: schools,
			User:    user,
		})
	})
}

type schoolListingDTO struct {
	Schools []*web.School
	User    *web.User
}

// RenderSchools renders a list of schools to give user general idea on what
// schools needs donations
func RenderSchools(store sessions.Store, dao web.SchoolDAO) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		schools, err := dao.GetSchools()
		user := getUser(store, r)
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
			User:    user,
		})
	})
}

type createSchoolDTO struct {
	User *web.User
}

// RenderCreateSchool renders the login page
func RenderCreateSchool(store sessions.Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := getUser(store, r)
		if user.Email == "" {
			http.Error(
				w,
				"Unauthorized",
				http.StatusUnauthorized,
			)
			return
		}
		t, err := template.ParseFiles("./templates/create.html", "./templates/base.html")
		if err != nil {
			panic(err)
		}
		t.ExecuteTemplate(w, "base", createSchoolDTO{
			User: user,
		})
	})
}

// HandleCSVUpload parses CSV and return JSON of the student loan detail
func HandleCSVUpload(store sessions.Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := getUser(store, r)
		if user.Email == "" {
			http.Error(
				w,
				"Unauthorized",
				http.StatusUnauthorized,
			)
			return
		}
		f, _, err := r.FormFile("file")
		if err != nil {
			http.Error(
				w,
				errors.Wrap(err, "Failed to parse file uploaded").Error(),
				http.StatusInternalServerError,
			)
			return
		}
		defer f.Close()

		reader := csv.NewReader(f)
		record, err := reader.ReadAll()
		if err != nil {
			http.Error(
				w,
				"Failed to read CSV file uploaded",
				http.StatusInternalServerError,
			)
			fmt.Println("Error handling CSV file upload:", err)
			return
		}

		donationDetails := []web.DonationDetail{}
		for i, line := range record {
			if i == 0 {
				// skip header
				continue
			}
			if len(line) < 5 {
				// skip any corrupted data row
				// TODO: maybe handle error better?
				fmt.Println("Row contains less than 4 columns", line)
				continue
			}

			// TODO: Switch back to a proper type for balances instead of string
			// lineWithoutSpecialCharacters := strings.Replace(line[4], "$", "", -1)
			// balance, err := strconv.ParseFloat(lineWithoutSpecialCharacters, 64)

			if err != nil {
				fmt.Println("Failed to parse balance on fourth column. Skipping", err)
				continue
			}
			donationDetails = append(donationDetails, web.DonationDetail{
				ID:          0, // bad magic number, I know!
				School:      line[1],
				Grade:       line[2],
				AccountName: line[3],
				Balance:     line[4],
			})
		}
		js, err := json.Marshal(donationDetails)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	})
}

// HandleCreateSchool creates school from the request form body
func HandleCreateSchool(store sessions.Store, dao web.SchoolDAO) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := getUser(store, r)
		if user.Email == "" {
			http.Error(
				w,
				"Unauthorized",
				http.StatusUnauthorized,
			)
			return
		}
		// parse JSON
		decoder := json.NewDecoder(r.Body)
		var s web.School
		err := decoder.Decode(&s)
		if err != nil {
			fmt.Println("Failed to parse JSON from request", err)
			http.Error(
				w,
				"Failed to parse JSON from request body",
				http.StatusBadRequest,
			)
			return
		}

		// call dao.CreateSchool
		_, err = dao.Create(s)
		if err != nil {
			fmt.Println("Failed to create new school", err)
			http.Error(
				w,
				"Failed to create new school",
				http.StatusInternalServerError,
			)
			return
		}
	})
}

type singleSchoolDTO struct {
	School *web.School
	User   *web.User
}

// RenderSchool renders individual school detail
func RenderSchool(store sessions.Store, dao web.SchoolDAO) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("./templates/school.html", "./templates/base.html")
		if err != nil {
			panic(err)
		}
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			http.Error(
				w,
				"Failed to grab id as int from route param",
				http.StatusBadRequest,
			)
			return
		}
		school, err := dao.GetSchool(id)
		if err != nil {
			http.Error(
				w,
				"Cannot retrieve school from database. Please try again later.",
				http.StatusInternalServerError,
			)
			return
		}
		t.ExecuteTemplate(w, "base", singleSchoolDTO{
			School: school,
		})
	})
}

// RenderEditSchool renders the login page
func RenderEditSchool(store sessions.Store, dao web.SchoolDAO) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := getUser(store, r)
		if user.Email == "" {
			http.Error(
				w,
				"Unauthorized",
				http.StatusUnauthorized,
			)
			return
		}
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			http.Error(
				w,
				"Failed to grab id as int from route param",
				http.StatusBadRequest,
			)
			return
		}
		school, err := dao.GetSchool(id)
		if err != nil {
			http.Error(
				w,
				"Cannot retrieve school from database. Please try again later.",
				http.StatusInternalServerError,
			)
			return
		}
		t, err := template.ParseFiles("./templates/edit.html", "./templates/base.html")
		if err != nil {
			panic(err)
		}
		t.ExecuteTemplate(w, "base", singleSchoolDTO{
			School: school,
			User:   user,
		})
	})
}

// HandleEditSchool edits the school by its id
func HandleEditSchool(store sessions.Store, dao web.SchoolDAO) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := getUser(store, r)
		if user.Email == "" {
			http.Error(
				w,
				"Unauthorized",
				http.StatusUnauthorized,
			)
			return
		}
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			http.Error(
				w,
				"Failed to grab id as int from route param",
				http.StatusBadRequest,
			)
			return
		}
		// parse JSON
		decoder := json.NewDecoder(r.Body)
		var s web.School
		err = decoder.Decode(&s)
		if err != nil {
			fmt.Println("Failed to parse JSON from request", err)
			http.Error(
				w,
				"Failed to parse JSON from request body",
				http.StatusBadRequest,
			)
			return
		}
		if err = dao.Edit(id, s); err != nil {
			fmt.Println("Failed to edit school.", err)
			http.Error(
				w,
				"Failed to edit school.",
				http.StatusInternalServerError,
			)
		}
		// return 200 or 500
	})
}
