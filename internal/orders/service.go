package orders

import (
	"fmt"
	"online-store-go/internal/items"
	"online-store-go/pkg/date_utils"
	"online-store-go/pkg/error_utils"
	"online-store-go/pkg/logger_utils"
	"time"
)

type service struct {
	repo      Repo
	itemsRepo items.Repo
}

type Service interface {
	CreateOrder(map[string]interface{}) *error_utils.RestErr
}

func NewService(repo Repo) Service {
	return &service{
		repo:      repo,
		itemsRepo: items.NewRepo(),
	}
}

func (s *service) CreateOrder(data map[string]interface{}) *error_utils.RestErr {
	// TODO: create order with the order items
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
	err := s.repo.CreateOrder(&request)
	if err != nil {
		return err
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

		item, err := s.itemsRepo.FindItemByID(int(orderItem.ItemID))
		if err != nil {
			return err
		}

		// validate if the stock is sufficient before order
		if item.Stock < orderItem.Quantity {
			logger_utils.Error("insufficient stock of items", fmt.Errorf("insufficient stock of items (%s)", item.Name))
			return error_utils.NewBadRequestError(fmt.Sprintf("insufficient stock of items (%s)", item.Name))
		}
	}

	err = s.repo.CreateOrderItems(&request, valueStrings, valueArgs)
	if err != nil {
		return err
	}
	return nil
}
