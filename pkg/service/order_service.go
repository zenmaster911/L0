package service

import (
	"github.com/google/uuid"
	"github.com/zenmaster911/L0/pkg/model"
	"github.com/zenmaster911/L0/pkg/repository"
)

type OrderService struct {
	repo repository.Order
}

func NewOrderService(repo repository.Order) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) GetOrderByUid(uid string) (model.Reply, error) {
	return s.repo.GetOrderByUid(uid)
}

func (s *OrderService) CreateOrder(input *model.Reply) (uid uuid.UUID, err error) {
	return s.repo.CreateOrder(input)
}
