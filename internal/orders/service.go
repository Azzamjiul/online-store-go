package orders

import (
	"online-store-go/pkg/error_utils"
)

type service struct {
	repo Repo
}

type Service interface {
	CreateOrder(map[string]interface{}) *error_utils.RestErr
}

func NewService(repo Repo) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateOrder(data map[string]interface{}) *error_utils.RestErr {
	// TODO: create order with the order items
	err := s.repo.CreateOrder(data)
	if err != nil {
		return err
	}
	return nil
}
