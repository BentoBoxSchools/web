package dao

import (
	"database/sql"

	"github.com/BentoBoxSchool/web"
)

// SchoolDAOImpl implements dao to interact with MySQL
type SchoolDAOImpl struct {
	db *sql.DB
}

// New constructs new SchoolDAOImpl object for interacting with MySQL
func New(db *sql.DB) *SchoolDAOImpl {
	return &SchoolDAOImpl{db}
}

func (s *SchoolDAOImpl) GetSchools() (error, []web.School) {
	return nil, nil
}
func (s *SchoolDAOImpl) GetSchool(id int) (error, web.School) {
	return nil, web.School{}
}
func (s *SchoolDAOImpl) Create(school web.School) error {
	return nil
}
func (s *SchoolDAOImpl) Edit(id int, school web.School) error {
	return nil
}
