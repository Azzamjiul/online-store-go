package orders

import (
	"database/sql"
	"fmt"
	"online-store-go/app/postgresql"
	"online-store-go/pkg/date_utils"
	"online-store-go/pkg/error_utils"
	"online-store-go/pkg/logger_utils"
	"strings"
	"time"
)

const (
	queryCreateOrder      = "INSERT INTO orders(user_id, status, total, expired_date) VALUES ($1, $2, $3, $4) RETURNING id"
	queryInsertOrderItems = "INSERT INTO order_items(order_id, item_id, quantity, price) VALUES "
)

type repo struct {
	Client *sql.DB
}

type Repo interface {
	CreateOrder(map[string]interface{}) *error_utils.RestErr
}

func NewRepo() Repo {
	return &repo{
		Client: postgresql.Client,
	}
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

	tx, txErr := r.Client.Begin()
	if txErr != nil {
		logger_utils.Error("error when begin transaction", txErr)
		return error_utils.NewInternalServerError("database error")
	}

	// Save Order to database
	stmt, err := r.Client.Prepare(queryCreateOrder)
	if err != nil {
		tx.Rollback()
		logger_utils.Error("error when trying to prepare create order statement, ", err)
		return error_utils.NewInternalServerError("database error")
	}
	defer stmt.Close()

	err = stmt.QueryRow(request.Order.UserID, 0, request.Order.Total, request.ExpiredDate).Scan(&request.Order.ID)
	if err != nil {
		tx.Rollback()
		logger_utils.Error("error when trying to create order", err)
		return error_utils.NewInternalServerError("database error")
	}

	// Order Items
	valueStrings := []string{}
	valueArgs := []interface{}{}
	for k, item := range data["items"].([]interface{}) {
		orderItem := OrderItem{
			ItemID:   uint64((item.(map[string]interface{})["item_id"]).(float64)),
			Quantity: int((item.(map[string]interface{})["quantity"]).(float64)),
			Price:    int((item.(map[string]interface{})["price"]).(float64)),
		}

		if data["total"] == nil {
			return error_utils.NewBadRequestError("total is missing")
		}

		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)", (k*4)+1, (k*4)+2, (k*4)+3, (k*4)+4))
		valueArgs = append(valueArgs, request.Order.ID)
		valueArgs = append(valueArgs, orderItem.ItemID)
		valueArgs = append(valueArgs, orderItem.Quantity)
		valueArgs = append(valueArgs, orderItem.Price)
	}

	// Save Order Items
	stmt, err = r.Client.Prepare(fmt.Sprintf("%s %s", queryInsertOrderItems, strings.Join(valueStrings, ",")))

	if err != nil {
		tx.Rollback()
		logger_utils.Error("error when trying to prepare insert order item statement, ", err)
		return error_utils.NewInternalServerError("database error")
	}
	defer stmt.Close()

	_, err = stmt.Exec(valueArgs...)
	if err != nil {
		tx.Rollback()
		logger_utils.Error("error when trying to insert order item", err)
		return error_utils.NewInternalServerError("database error")
	}

	txErr = tx.Commit()
	if txErr != nil {
		logger_utils.Error("error when commit transaction", txErr)
		return error_utils.NewInternalServerError("database error")
	}

	return nil
}
