package orders

import (
	"fmt"
	"online-store-go/pkg/error_utils"
)

type service struct {
	repo Repo
}

type Service interface {
	CreateOrder(interface{}) *error_utils.RestErr
}

func NewService(repo Repo) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateOrder(data interface{}) *error_utils.RestErr {
	fmt.Println(data)
	return nil
}
