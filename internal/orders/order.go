package orders

// Order Entity
type Order struct {
	ID          uint64 `json:"id"`
	UserID      uint64 `json:"user_id"`
	Status      int    `json:"status"`
	Total       string `json:"total"`
	ExpiredDate string `json:"expired_date"`
}

// Create Order Request
type CreateOrderRequest struct {
	Order

	CreatedAt string `json:"created_at"`
}
