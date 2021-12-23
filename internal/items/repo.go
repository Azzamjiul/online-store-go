package items

import (
	"database/sql"
	"online-store-go/app/postgresql"
	"online-store-go/pkg/error_utils"
	"online-store-go/pkg/logger_utils"
)

const (
	queryFindItemByID    = "SELECT sku, name, stock, price FROM items WHERE id = $1"
	queryUpdateItemStock = "UPDATE items SET stock = $1 WHERE id = $2"
)

type repo struct {
	Client *sql.DB
}

type Repo interface {
	FindItemByID(int) (*Item, *error_utils.RestErr)
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
