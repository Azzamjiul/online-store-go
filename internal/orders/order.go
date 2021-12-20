package orders

// Order Entity
type Order struct {
	ID          int    `json:"id"`
	UserID      string `json:"user_id"`
	Status      int    `json:"status"`
	Total       string `json:"total"`
	ExpiredDate string `json:"expired_date"`
}
