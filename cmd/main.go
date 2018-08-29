package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BentoBoxSchool/web/dao"
	"github.com/BentoBoxSchool/web/handlers"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var (
	port       = os.Getenv("PORT")
	dbUsername = os.Getenv("DB_USERNAME")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbHost     = os.Getenv("DB_HOST")
)

func main() {
	r := mux.NewRouter()
	sqlDB := getDB()

	schoolDAO := dao.New(sqlDB)

	// Ops end points
	r.HandleFunc("/hello", handlers.Hello()).Methods("GET")
	r.HandleFunc("/health", handlers.CheckHealth(sqlDB)).Methods("GET")

	// renders javascript and css under "/assets"
	r.PathPrefix("/assets").Handler(handlers.Assets())

	// schools
	r.HandleFunc("/", handlers.RenderHomepage(schoolDAO)).Methods("GET")
	/* pre-registered pages
	r.HandleFunc("/api/csv/parse", handlers.RenderHomepage(schoolDAO)).Methods("POST")
	r.HandleFunc("/login", handlers.RenderHomepage(schoolDAO)).Methods("GET")
	r.HandleFunc("/login", handlers.RenderHomepage(schoolDAO)).Methods("POST")
	r.HandleFunc("/logout", handlers.RenderHomepage(schoolDAO)).Methods("POST")
	r.HandleFunc("/schools", handlers.RenderHomepage(schoolDAO)).Methods("GET")
	r.HandleFunc("/schools/create", handlers.RenderHomepage(schoolDAO)).Methods("GET")
	r.HandleFunc("/schools/create", handlers.RenderHomepage(schoolDAO)).Methods("POST")
	r.HandleFunc("/schools/edit/{id}", handlers.RenderHomepage(schoolDAO)).Methods("GET")
	r.HandleFunc("/schools/edit/{id}", handlers.RenderHomepage(schoolDAO)).Methods("POST")
	*/

	log.Printf("Running web server at port %s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}

func getDB() *sql.DB {
	defaultProtocol := "tcp"
	defaultPort := "3306"

	sqlDSN := fmt.Sprintf(
		"%s:%s@%s(%s:%s)/",
		dbUsername,
		dbPassword,
		defaultProtocol,
		dbHost,
		defaultPort,
	)

	db, err := sql.Open("mysql", sqlDSN)
	if err != nil {
		panic(err)
	}

	return db
}
