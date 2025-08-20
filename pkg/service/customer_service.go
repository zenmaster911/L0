package service

import (
	"github.com/zenmaster911/L0/pkg/model"
	"github.com/zenmaster911/L0/pkg/repository"
)

type CustomerService struct {
	repo repository.Customer
}

func NewCustomerservice(repo repository.Customer) *CustomerService {
	return &CustomerService{repo: repo}
}

func (s *CustomerService) GetCustomerByPhone(phone string) (model.Customer, error) {
	return s.repo.GetCustomerByPhone(phone)
}
