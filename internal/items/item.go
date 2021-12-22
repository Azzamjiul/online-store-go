package items

type Item struct {
	ID    uint64 `json:"id"`
	SKU   string `json:"sku"`
	Name  string `json:"name"`
	Stock int    `json:"stock"`
	Price int    `json:"price"`
}
