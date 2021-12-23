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

		request.Items = append(request.Items, orderItem)

		if data["total"] == nil {
			s.repo.DeleteOrderByID(int(request.Order.ID))
			return error_utils.NewBadRequestError("total is missing")
		}

		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)", (k*4)+1, (k*4)+2, (k*4)+3, (k*4)+4))
		valueArgs = append(valueArgs, request.Order.ID)
		valueArgs = append(valueArgs, orderItem.ItemID)
		valueArgs = append(valueArgs, orderItem.Quantity)
		valueArgs = append(valueArgs, orderItem.Price)

		item, err := s.itemsRepo.FindItemByID(int(orderItem.ItemID))
		if err != nil {
			s.repo.DeleteOrderByID(int(request.Order.ID))
			return err
		}

		// validate if the stock is sufficient before order
		if item.Stock < orderItem.Quantity {
			s.repo.DeleteOrderByID(int(request.Order.ID))
			logger_utils.Error("insufficient stock of items", fmt.Errorf("insufficient stock of items (%s)", item.Name))
			return error_utils.NewBadRequestError(fmt.Sprintf("insufficient stock of items (%s)", item.Name))
		}
	}

	err = s.repo.CreateOrderItems(&request, valueStrings, valueArgs)
	if err != nil {
		s.repo.DeleteOrderByID(int(request.Order.ID))
		s.repo.DeleteAllOrderItemsByOrderID(int(request.Order.ID))
		return err
	}

	// Update items stock after all order items recorded
	valueStrings = []string{}
	valueArgs = []interface{}{}
	for k, v := range request.Items {
		item, err := s.itemsRepo.FindItemByID(int(v.ItemID))
		if err != nil {
			s.repo.DeleteAllOrderItemsByOrderID(int(request.Order.ID))
			s.repo.DeleteOrderByID(int(request.Order.ID))
			return err
		}

		totalColumn := 5
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", (k*totalColumn)+1, (k*totalColumn)+2, (k*totalColumn)+3, (k*totalColumn)+4, (k*totalColumn)+5))
		valueArgs = append(valueArgs, item.ID)
		valueArgs = append(valueArgs, item.SKU)
		valueArgs = append(valueArgs, item.Name)
		valueArgs = append(valueArgs, (item.Stock - v.Quantity))
		valueArgs = append(valueArgs, item.Price)
	}

	err = s.itemsRepo.UpsertItemsStock(valueStrings, valueArgs)
	if err != nil {
		s.repo.DeleteAllOrderItemsByOrderID(int(request.Order.ID))
		s.repo.DeleteOrderByID(int(request.Order.ID))
		return err
	}

	return nil
}
