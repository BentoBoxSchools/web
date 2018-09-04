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
	"github.com/pkg/errors"
)

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
		fmt.Println(schools)
		t.ExecuteTemplate(w, "base", schoolListingDTO{
			Schools: schools,
		})
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

// HandleCSVUpload parses CSV and return JSON of the student loan detail
func HandleCSVUpload() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
func HandleCreateSchool(dao web.SchoolDAO) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
}

// RenderSchool renders individual school detail
func RenderSchool(dao web.SchoolDAO) http.HandlerFunc {
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
func RenderEditSchool(dao web.SchoolDAO) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("./templates/edit.html", "./temnplates/base.html")
		if err != nil {
			panic(err)
		}
		t.ExecuteTemplate(w, "base", nil)
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
