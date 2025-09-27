package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/zenmaster911/L0/pkg/model"
)

//go:generate minimock -i github.com/zenmaster911/L0/pkg/repository.* -o ./repo_mocks -s _mock.go

type Cache interface {
	CacheLoad(limit int) ([]string, error)
}

type Order interface {
	GetOrderByUid(uid string) (model.Reply, error)
	CreateOrder(input *model.Reply) (uid string, err error)
	CheckOrderExists(uid string) (bool, error)
}

type Item interface {
	GetItemByArticle(nmId int) (model.Item, error)
}

type Delivery interface {
	GetCustomerDeliveryByAddress(address, customerUid string) (model.Delivery, error)
}

type Customer interface {
	GetCustomerByPhone(phone string) (model.Customer, error)
}

type StatusCheck interface {
	DBConnectionCheck() error
}

type Repository struct {
	Order
	Item
	Customer
	Delivery
	StatusCheck
	Cache
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Order:       NewOrderPostgres(db),
		Item:        NewItemsPostgres(db),
		Customer:    NewCustomerPostgres(db),
		Delivery:    NewDeliveryPostgrtes(db),
		StatusCheck: NewStatusCheck(db),
		Cache:       NewCachePostgres(db),
	}
}

func (r *Repository) ConnectionCheck(db *sqlx.DB) error {
	return db.Ping()
}
