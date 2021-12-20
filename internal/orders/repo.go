package orders

import (
	"online-store-go/pkg/error_utils"
)

const (
	queryCreateOrder = "INSERT INTO orders (user_id, status, )"
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
	// request := CreateOrderRequest{}

	if data["user_id"] == nil {
		return error_utils.NewBadRequestError("user_id is missing")
	}

	return nil
}
