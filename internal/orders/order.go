package orders

// Order Entity
type Order struct {
	ID          uint64 `json:"id"`
	UserID      uint64 `json:"user_id"`
	Status      int    `json:"status"`
	Total       uint64 `json:"total"`
	ExpiredDate string `json:"expired_date"`
}

// Order Item Element Request
type OrderItem struct {
	ItemID   uint64 `json:"item_id"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
}

// Create Order Request
type CreateOrderRequest struct {
	Order
	Items     []OrderItem
	CreatedAt string `json:"created_at"`
}
