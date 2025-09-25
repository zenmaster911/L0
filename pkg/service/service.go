package service

import (
	"github.com/zenmaster911/L0/pkg/model"
	"github.com/zenmaster911/L0/pkg/repository"
)

//go:generate minimock -i github.com/zenmaster911/L0/pkg/service.* -o ./mocks -s _mock.go
type Cache interface {
	CacheLoad(limit int) ([]string, error)
}

type Order interface {
	GetOrderByUid(uid string) (model.Reply, error)
	CreateOrder(input *model.Reply) (uid string, err error)
	CheckOrderExists(uid string) (bool, error)
}

type Customer interface {
	GetCustomerByPhone(phone string) (model.Customer, error)
}

type Item interface {
	GetItemByArticle(nmId int) (model.Item, error)
}

type Delivery interface {
	GetCustomerDeliveryByAddress(address, customerUid string) (model.Delivery, error)
}

type Service struct {
	Order
	Customer
	Delivery
	Item
	Cache
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Order:    NewOrderService(repo.Order),
		Customer: NewCustomerservice(repo.Customer),
		Delivery: NewDeliveryService(repo.Delivery),
		Item:     NewItemService(repo.Item),
		Cache:    NewCacheService(repo.Cache),
	}
}
