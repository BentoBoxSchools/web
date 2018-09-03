package web

// School represents individual school needing donations
type School struct {
	ID          int
	Name        string
	Description string
	Link        string
	Data        []DonationDetail
}

// DonationDetail represents individual donation detail (account, balance)
type DonationDetail struct {
	ID          int
	School      string // Is this needed?
	Grade       string
	AccountName string
	Balance     float64
}

// SchoolDAO represents common business behavior to retrieve schools
type SchoolDAO interface {
	GetSchools() ([]*School, error)
	GetSchool(id int64) (*School, error)
	Create(s School) (int64, error)
	Edit(id int64, s School) error
}
