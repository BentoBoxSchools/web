package dao

import (
	"database/sql"
	"fmt"

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

	m := make(map[int64]*web.School)
	for rows.Next() {
		var (
			sId           int64
			sName         string
			sDescription  string
			sLink         string
			ddId          sql.NullInt64
			ddGrade       sql.NullString
			ddAccountName sql.NullString
			ddBalance     sql.NullFloat64
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
			ID:          ddId.Int64,
			Grade:       ddGrade.String,
			AccountName: ddAccountName.String,
			Balance:     ddBalance.Float64,
		}

		school.Data = append(school.Data, dd)
	}

	out := []*web.School{}
	for _, v := range m {
		out = append(out, v)
	}

	return out, nil
}
func (s *SchoolDAOImpl) GetSchool(id int64) (*web.School, error) {
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
			sId           int64
			sName         string
			sDescription  string
			sLink         string
			ddId          sql.NullInt64
			ddGrade       sql.NullString
			ddAccountName sql.NullString
			ddBalance     sql.NullFloat64
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
			ID:          ddId.Int64,
			Grade:       ddGrade.String,
			AccountName: ddAccountName.String,
			Balance:     ddBalance.Float64,
		}

		school.Data = append(school.Data, dd)
	}

	return school, nil
}
func (s *SchoolDAOImpl) Create(school web.School) (int64, error) {
	var err error

	if school.ID != 0 {
		return 0, fmt.Errorf("received school value with id=%d. cannot create entry with non-zero id value", school.ID)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	stmt, err := s.db.Prepare(`
		INSERT INTO school(name, description, link)
		VALUES(?, ?, ?)
	`)
	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(school.Name, school.Description, school.Link)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	stmt, err = s.db.Prepare(`
		INSERT INTO donation_detail (school_id, grade, account_name, balance)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return 0, err
	}

	for _, dd := range school.Data {
		_, err = stmt.Exec(id, dd.Grade, dd.AccountName, dd.Balance)
		if err != nil {
			return 0, err
		}
	}

	return id, nil
}
func (s *SchoolDAOImpl) Update(school web.School) error {
	var err error

	if school.ID == 0 {
		return fmt.Errorf("received school value with id=%d. cannot create entry with zero id value", school.ID)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	stmt, err := s.db.Prepare(`
		UPDATE school(name, description, link)
		VALUES(?, ?, ?)
		WHERE id=?
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(school.Name, school.Description, school.Link, school.ID)
	if err != nil {
		return err
	}

	for _, dd := range school.Data {
		stmt, err := s.db.Prepare(`
			UPDATE donation_detail(school_id, grade, account_name, balance)
			VALUES(?, ?, ?, ?)
			WHERE id=?
		`)
		if err != nil {
			return err
		}

		_, err = stmt.Exec(school.ID, dd.Grade, dd.AccountName, dd.Balance, dd.ID)
	}

	return nil
}

func (s *SchoolDAOImpl) Edit(id int64, school web.School) error {
	return nil
}
