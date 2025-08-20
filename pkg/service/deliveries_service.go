package service

import (
	"github.com/zenmaster911/L0/pkg/model"
	"github.com/zenmaster911/L0/pkg/repository"
)

type DeliveriesService struct {
	repo repository.Delivery
}

func NewDeliveryService(repo repository.Delivery) *DeliveriesService {
	return &DeliveriesService{repo: repo}
}

func (s *DeliveriesService) GetCustomerDeliveryByAddress(address, customerUid string) (model.Delivery, error) {
	return s.repo.GetCustomerDeliveryByAddress(address, customerUid)
}
