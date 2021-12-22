package orders

import (
	"fmt"
	"online-store-go/pkg/error_utils"
)

const (
// queryCreateOrder = "INSERT INTO orders (user_id, status, )"
)

type repo struct {
}

type Repo interface {
	CreateOrder(map[string]interface{}) *error_utils.RestErr
}

func NewRepo() Repo {
	return &repo{}
}

func (r *repo) CreateOrder(data map[string]interface{}) *error_utils.RestErr {
	request := CreateOrderRequest{}

	// Order properties
	if data["user_id"] == nil {
		return error_utils.NewBadRequestError("user_id is missing")
	}
	request.UserID = uint64(data["user_id"].(float64))

	if data["total"] == nil {
		return error_utils.NewBadRequestError("total is missing")
	}
	request.Total = uint64(data["total"].(float64))

	// Order Items
	for _, item := range data["items"].([]interface{}) {
		orderItem := OrderItem{
			ItemID:   uint64((item.(map[string]interface{})["item_id"]).(float64)),
			Quantity: int((item.(map[string]interface{})["quantity"]).(float64)),
			Price:    int((item.(map[string]interface{})["price"]).(float64)),
		}

		request.Items = append(request.Items, orderItem)
	}

	fmt.Println(request)
	return nil
}
