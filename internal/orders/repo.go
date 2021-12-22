package orders

import (
	"online-store-go/app/postgresql"
	"online-store-go/pkg/date_utils"
	"online-store-go/pkg/error_utils"
	"online-store-go/pkg/logger_utils"
	"time"
)

const (
	queryCreateOrder = "INSERT INTO orders(user_id, status, total, expired_date) VALUES ($1, $2, $3, $4) RETURNING id"
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

	request.ExpiredDate = date_utils.GetDBFormat(time.Now().UTC().Add(10 * time.Hour))

	// Order Items
	for _, item := range data["items"].([]interface{}) {
		orderItem := OrderItem{
			ItemID:   uint64((item.(map[string]interface{})["item_id"]).(float64)),
			Quantity: int((item.(map[string]interface{})["quantity"]).(float64)),
			Price:    int((item.(map[string]interface{})["price"]).(float64)),
		}

		request.Items = append(request.Items, orderItem)
	}

	// Save Order to database
	stmt, err := postgresql.Client.Prepare(queryCreateOrder)
	if err != nil {
		logger_utils.Error("error when trying to prepare create order statement, ", err)
		return error_utils.NewInternalServerError("database error")
	}
	defer stmt.Close()

	err = stmt.QueryRow(request.Order.UserID, 0, request.Order.Total, request.ExpiredDate).Scan(&request.Order.ID)
	if err != nil {
		logger_utils.Error("error when trying to create order", err)
		return error_utils.NewInternalServerError("database error")
	}

	// Save Order Items to database

	return nil
}
