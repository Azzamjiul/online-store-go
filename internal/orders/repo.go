package orders

import "online-store-go/pkg/error_utils"

type repo struct {
}

type Repo interface {
	CreateOrder(interface{}) *error_utils.RestErr
}

func NewRepo() Repo {
	return &repo{}
}

func (r *repo) CreateOrder(data interface{}) *error_utils.RestErr {
	return nil
}
