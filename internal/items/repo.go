package items

import (
	"database/sql"
	"fmt"
	"online-store-go/app/postgresql"
	"online-store-go/pkg/error_utils"
	"online-store-go/pkg/logger_utils"
	"strings"
)

const (
	queryFindItemByID        = "SELECT sku, name, stock, price FROM items WHERE id = $1"
	queryUpdateItemStock     = "UPDATE items SET stock = $1 WHERE id = $2"
	queryUpsertItemsStock    = "INSERT INTO items(id, sku, name, stock, price) VALUES "
	queryUpsertItemsStockEnd = "ON CONFLICT (sku) DO UPDATE SET stock = EXCLUDED.stock"
)

type repo struct {
	Client *sql.DB
}

type Repo interface {
	FindItemByID(int) (*Item, *error_utils.RestErr)
	UpdateItemStockByID(int, int) (*Item, *error_utils.RestErr)
	UpsertItemsStock([]string, []interface{}) *error_utils.RestErr
}

func NewRepo() Repo {
	return &repo{
		Client: postgresql.Client,
	}
}

func (r *repo) FindItemByID(id int) (*Item, *error_utils.RestErr) {
	stmt, err := r.Client.Prepare(queryFindItemByID)
	if err != nil {
		logger_utils.Error("error when prepare find item by id statement", err)
		return nil, error_utils.NewInternalServerError("database error")
	}
	defer stmt.Close()

	var item Item
	err = stmt.QueryRow(id).Scan(&item.SKU, &item.Name, &item.Stock, &item.Price)
	if err != nil {
		logger_utils.Error("error when execute find item by id statement", err)
		return nil, error_utils.NewInternalServerError("database error")
	}

	return &item, nil
}

func (r *repo) UpdateItemStockByID(id int, newStock int) (*Item, *error_utils.RestErr) {
	stmt, err := r.Client.Prepare(queryUpdateItemStock)
	if err != nil {
		logger_utils.Error("error when prepare update item stock by id statement", err)
		return nil, error_utils.NewInternalServerError("database error")
	}
	defer stmt.Close()

	var item Item
	err = stmt.QueryRow(newStock, id).Scan(&item.SKU, &item.Name, &item.Stock, &item.Price)
	if err != nil {
		logger_utils.Error(fmt.Sprintf("error when execute update item stock by id statement for item with id = %v", id), err)
		return nil, error_utils.NewInternalServerError("database error")
	}

	logger_utils.Info(fmt.Sprintf("Item id %v. Updated Stock %v", item.ID, item.Stock))

	return &item, nil
}

func (r *repo) UpsertItemsStock(valueStrings []string, valueArgs []interface{}) *error_utils.RestErr {
	tx, txErr := r.Client.Begin()
	if txErr != nil {
		logger_utils.Error("error when begin transaction", txErr)
		return error_utils.NewInternalServerError("database error")
	}

	// Save Order Items
	statement := fmt.Sprintf("%s %s %s", queryUpsertItemsStock, strings.Join(valueStrings, ","), queryUpsertItemsStockEnd)

	stmt, err := r.Client.Prepare(statement)
	if err != nil {
		tx.Rollback()
		logger_utils.Error("error when trying to prepare upsert items stock statement, ", err)
		return error_utils.NewInternalServerError("database error")
	}
	defer stmt.Close()

	_, err = stmt.Exec(valueArgs...)
	if err != nil {
		tx.Rollback()
		logger_utils.Error("error when trying to execute upsert items stock", err)
		return error_utils.NewInternalServerError("database error")
	}

	txErr = tx.Commit()
	if txErr != nil {
		logger_utils.Error("error when commit transaction", txErr)
		return error_utils.NewInternalServerError("database error")
	}

	return nil
}
