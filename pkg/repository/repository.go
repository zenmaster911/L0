package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/zenmaster911/L0/pkg/model"
)

type Order interface {
	GetOrderByUid(uid string) (model.Reply, error)
}

type Item interface {
	GetItemByArticle(nmId int) (model.Item, int, error)
}

type Delivery interface {
	GetCustomerDeliveryByAddress(address, customerUid string) (model.Delivery, int, error)
}

type Customer interface {
	GetCustomerByPhone(phone string) (model.Customer, string, error)
}

type Repository struct {
	Order
	Item
	Customer
	Delivery
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Order:    NewOrderPostgres(db),
		Item:     NewItemsPostgres(db),
		Customer: NewCustomerPostgres(db),
		Delivery: NewDeliveryPostgrtes(db),
	}
}
