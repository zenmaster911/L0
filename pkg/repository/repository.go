package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/zenmaster911/L0/pkg/model"
)

type Order interface {
	GetOrderByUid(uid string) (model.Reply, error)
}

type Repository struct {
	Order
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Order: NewOrderPostgres(db),
	}
}
