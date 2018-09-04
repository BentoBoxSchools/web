package main

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/BentoBoxSchool/web"
	"github.com/BentoBoxSchool/web/dao"
	"github.com/BentoBoxSchool/web/handlers"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var (
	port                    = os.Getenv("PORT")
	dbUsername              = os.Getenv("DB_USERNAME")
	dbPassword              = os.Getenv("DB_PASSWORD")
	dbHost                  = os.Getenv("DB_HOST")
	dbName                  = os.Getenv("DB_NAME")
	googleClientID          = os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret      = os.Getenv("GOOGLE_CLIENT_SECRET")
	googleRedirectURI       = os.Getenv("GOOGLE_REDIRECT_URI")
	googleWhitelistedEmails = strings.Split(os.Getenv("GOOGLE_WHITE_LIST_EMAILS"), ",")
	sessionSecret           = os.Getenv("SESSION_SECRET")
)

var store = sessions.NewCookieStore([]byte(sessionSecret))

func init() {
	gob.Register(&web.User{})
}

func main() {
	r := mux.NewRouter()
	sqlDB := getDB()

	schoolDAO := dao.New(sqlDB)

	// Ops end points
	r.HandleFunc("/hello", handlers.Hello()).Methods("GET")
	r.HandleFunc("/health", handlers.CheckHealth(sqlDB)).Methods("GET")

	// renders javascript and css under "/assets"
	r.PathPrefix("/assets").Handler(handlers.Assets())

	// authentication & authorization
	r.HandleFunc("/login", handlers.RedirectGoogleLogin(googleClientID, googleRedirectURI)).Methods("GET")
	r.HandleFunc("/google/callback", handlers.HandleGoogleCallback(store, googleClientID, googleClientSecret, googleRedirectURI, googleWhitelistedEmails)).Methods("GET")
	r.HandleFunc("/logout", handlers.HandleLogout(store)).Methods("GET")

	// schools
	r.HandleFunc("/", handlers.RenderHomepage(store, schoolDAO)).Methods("GET")
	r.HandleFunc("/schools", handlers.RenderSchools(store, schoolDAO)).Methods("GET")
	r.HandleFunc("/schools/create", handlers.RenderCreateSchool(store)).Methods("GET")
	r.HandleFunc("/schools/{id}", handlers.RenderSchool(store, schoolDAO)).Methods("GET")
	r.HandleFunc("/schools/edit/{id}", handlers.RenderEditSchool(store, schoolDAO)).Methods("GET")
	r.HandleFunc("/schools/create", handlers.HandleCreateSchool(store, schoolDAO)).Methods("POST")
	r.HandleFunc("/schools/edit/{id}", handlers.HandleEditSchool(store, schoolDAO)).Methods("POST")
	r.HandleFunc("/api/csv/donation", handlers.HandleCSVUpload(store)).Methods("POST")

	log.Printf("Running web server at port %s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}

func getDB() *sql.DB {
	defaultProtocol := "tcp"
	defaultPort := "3306"

	sqlDSN := fmt.Sprintf(
		"%s:%s@%s(%s:%s)/%s",
		dbUsername,
		dbPassword,
		defaultProtocol,
		dbHost,
		defaultPort,
		dbName,
	)

	db, err := sql.Open("mysql", sqlDSN)
	if err != nil {
		panic(err)
	}

	return db
}
