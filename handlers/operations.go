package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

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
