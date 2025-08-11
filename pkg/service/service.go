package service

import (
	"github.com/zenmaster911/L0/pkg/model"
	"github.com/zenmaster911/L0/pkg/repository"
)

type Order interface {
	GetOrderByUid(uid string) (model.Reply, error)
}

type Service struct {
	Order
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Order: NewOrderService(repo.Order),
	}
}
