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

func (s *SchoolDAOImpl) GetSchools() ([]*web.School, error) {
	rows, err := s.db.Query(`
		SELECT s.id, s.name, s.description, s.link, dd.id, dd.grade, dd.account_name, dd.balance
		FROM school AS s
		LEFT JOIN donation_detail AS dd
		on s.id = dd.school_id
	`)
	if err != nil {
		return []*web.School{}, err
	}

	m := make(map[int]*web.School)
	for rows.Next() {
		var (
			sId           int
			sName         string
			sDescription  string
			sLink         string
			ddId          int
			ddGrade       string
			ddAccountName string
			ddBalance     float64
		)

		if err := rows.Scan(&sId, &sName, &sDescription, &sLink, &ddId, &ddGrade, &ddAccountName, &ddBalance); err != nil {
			return []*web.School{}, err
		}

		school := m[sId]
		if school == nil {
			school = &web.School{
				ID:          sId,
				Name:        sName,
				Description: sDescription,
				Link:        sLink,
			}
			m[sId] = school
		}

		dd := web.DonationDetail{
			ID:          ddId,
			Grade:       ddGrade,
			AccountName: ddAccountName,
			Balance:     ddBalance,
		}

		school.Data = append(school.Data, dd)
	}

	out := []*web.School{}
	for _, v := range m {
		out = append(out, v)
	}

	return out, nil
}
func (s *SchoolDAOImpl) GetSchool(id int) (*web.School, error) {
	rows, err := s.db.Query(`
		SELECT s.id, s.name, s.description, s.link, dd.id, dd.grade, dd.account_name, dd.balance
		FROM school AS s
		LEFT JOIN donation_detail AS dd
		on s.id = dd.school_id
		WHERE s.id = ?
	`, id)
	if err != nil {
		return &web.School{}, err
	}

	var school *web.School
	for rows.Next() {
		var (
			sId           int
			sName         string
			sDescription  string
			sLink         string
			ddId          int
			ddGrade       string
			ddAccountName string
			ddBalance     float64
		)

		if err := rows.Scan(&sId, &sName, &sDescription, &sLink, &ddId, &ddGrade, &ddAccountName, &ddBalance); err != nil {
			return &web.School{}, err
		}

		if school == nil {
			school = &web.School{
				ID:          sId,
				Name:        sName,
				Description: sDescription,
				Link:        sLink,
			}
		}

		dd := web.DonationDetail{
			ID:          ddId,
			Grade:       ddGrade,
			AccountName: ddAccountName,
			Balance:     ddBalance,
		}

		school.Data = append(school.Data, dd)
	}

	return school, nil
}
func (s *SchoolDAOImpl) Create(school web.School) error {
	return nil
}
func (s *SchoolDAOImpl) Edit(id int, school web.School) error {
	return nil
}
