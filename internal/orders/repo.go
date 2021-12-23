package orders

import (
	"database/sql"
	"fmt"
	"online-store-go/app/postgresql"
	"online-store-go/pkg/error_utils"
	"online-store-go/pkg/logger_utils"
	"strings"
)

const (
	queryCreateOrder      = "INSERT INTO orders(user_id, status, total, expired_date) VALUES ($1, $2, $3, $4) RETURNING id"
	queryInsertOrderItems = "INSERT INTO order_items(order_id, item_id, quantity, price) VALUES "
	queryFindItemByID     = "SELECT sku, name, stock, price FROM items WHERE id = $1"
	queryUpdateItemStock  = "UPDATE items SET stock = $1 WHERE id = $2"
)

type repo struct {
	Client *sql.DB
}

type Repo interface {
	CreateOrder(*CreateOrderRequest) *error_utils.RestErr
	CreateOrderItems(*CreateOrderRequest, []string, []interface{}) *error_utils.RestErr
}

func NewRepo() Repo {
	return &repo{
		Client: postgresql.Client,
	}
}

func (r *repo) CreateOrder(request *CreateOrderRequest) *error_utils.RestErr {
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

	txErr = tx.Commit()
	if txErr != nil {
		logger_utils.Error("error when commit transaction", txErr)
		return error_utils.NewInternalServerError("database error")
	}

	return nil
}

func (r *repo) CreateOrderItems(request *CreateOrderRequest, valueStrings []string, valueArgs []interface{}) *error_utils.RestErr {
	tx, txErr := r.Client.Begin()
	if txErr != nil {
		logger_utils.Error("error when begin transaction", txErr)
		return error_utils.NewInternalServerError("database error")
	}

	// Save Order Items
	stmt, err := r.Client.Prepare(fmt.Sprintf("%s %s", queryInsertOrderItems, strings.Join(valueStrings, ",")))
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
