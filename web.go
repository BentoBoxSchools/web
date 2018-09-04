package web

// School represents individual school needing donations
type School struct {
	ID          int64            `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Link        string           `json:"link"`
	Data        []DonationDetail `json:"donationDetail"`
}

// DonationDetail represents individual donation detail (account, balance)
type DonationDetail struct {
	ID          int64  `json:"id"`
	School      string `json:"school"` // Is this needed? Eric: Yes, we need to differeniate between high school and middle school under district
	Grade       string `json:"grade"`
	AccountName string `json:"accountName"`
	Balance     string `json:"balance"`
	// Balance     float64 `json:"balance"`	// Temporarily using string
}

type User struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

// SchoolDAO represents common business behavior to retrieve schools
type SchoolDAO interface {
	GetSchools() ([]*School, error)
	GetSchool(id int64) (*School, error)
	Create(s School) (int64, error)
	Edit(id int64, s School) error
}
