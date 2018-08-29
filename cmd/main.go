package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

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

	// Ops end points
	r.HandleFunc("/hello", handlers.Hello()).Methods("GET")
	r.HandleFunc("/health", handlers.CheckHealth(sqlDB)).Methods("GET")

	// schools
	r.HandleFunc("/", handlers.Homepage(sqlDB)).Methods("GET")

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
