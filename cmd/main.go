package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/BentoBoxSchools/web"
	"github.com/BentoBoxSchools/web/dao"
	"github.com/BentoBoxSchools/web/handlers"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/acme/autocert"
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

func makeRouter(db *sql.DB, dao *dao.SchoolDAOImpl) *mux.Router {
	r := mux.NewRouter()

	// Ops end points
	r.HandleFunc("/hello", handlers.Hello()).Methods("GET")
	r.HandleFunc("/health", handlers.CheckHealth(db)).Methods("GET")

	// renders javascript and css under "/assets"
	r.PathPrefix("/assets").Handler(handlers.Assets())

	// authentication & authorization
	r.HandleFunc("/login", handlers.RedirectGoogleLogin(googleClientID, googleRedirectURI)).Methods("GET")
	r.HandleFunc("/google/callback", handlers.HandleGoogleCallback(store, googleClientID, googleClientSecret, googleRedirectURI, googleWhitelistedEmails)).Methods("GET")
	r.HandleFunc("/logout", handlers.HandleLogout(store)).Methods("GET")

	// schools
	r.HandleFunc("/", handlers.RenderHomepage(store, dao)).Methods("GET")
	r.HandleFunc("/schools", handlers.RenderSchools(store, dao)).Methods("GET")
	r.HandleFunc("/schools/create", handlers.RenderCreateSchool(store)).Methods("GET")
	r.HandleFunc("/schools/{id}", handlers.RenderSchool(store, dao)).Methods("GET")
	r.HandleFunc("/schools/edit/{id}", handlers.RenderEditSchool(store, dao)).Methods("GET")
	r.HandleFunc("/schools/create", handlers.HandleCreateSchool(store, dao)).Methods("POST")
	r.HandleFunc("/schools/edit/{id}", handlers.HandleEditSchool(store, dao)).Methods("POST")
	r.HandleFunc("/api/csv/donation", handlers.HandleCSVUpload(store)).Methods("POST")

	return r
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

func makeAutocertManager() *autocert.Manager {
	return &autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache:  autocert.DirCache("."),
		HostPolicy: func(ctx context.Context, host string) error {
			if host == "bentoboxschools.com" {
				return nil
			}
			return fmt.Errorf("host policy '%s' violated; only 'bentoboxschools.com' is allowed", host)
		},
	}
}

func runHTTPSServer(r *mux.Router, m *autocert.Manager) {
	plain := &http.Server{
		Addr:    ":80",
		Handler: m.HTTPHandler(nil),
	}
	go func() {
		log.Println("started plain traffic listener on port 80")
		if err := plain.ListenAndServe(); err != nil {
			log.Println("an error occurred with the plain traffic listener on port 80", err)
		}
	}()

	secure := &http.Server{
		Addr:    ":443",
		Handler: r,
		TLSConfig: &tls.Config{
			GetCertificate: m.GetCertificate,
		},
	}

	log.Println("started secure traffic listener on port 443")
	if err := secure.ListenAndServeTLS("", ""); err != nil {
		log.Println("an error occurred with the secure traffic listener of port 443", err)
	}
}

func runHTTPServer(r *mux.Router) {
	s := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	log.Printf("started plain traffic listener on port %s\n", port)
	if err := s.ListenAndServe(); err != nil {
		log.Printf("an error occurred with the plain traffic listener on port %d\n%s\n", port, err)
	}
}

func main() {
	db := getDB()
	schoolDAO := dao.New(db)
	r := makeRouter(db, schoolDAO)

	production := *flag.Bool("production", false, "Enables HTTPS traffic on port 443 (HTTP requests on port 80 are redirected in this mode)")
	if production {
		m := makeAutocertManager()
		runHTTPSServer(r, m)
	} else {
		runHTTPServer(r)
	}
}
