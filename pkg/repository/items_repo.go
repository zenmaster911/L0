package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/zenmaster911/L0/pkg/model"
)

type ItemPostgres struct {
	db *sqlx.DB
}

func NewItemsPostgres(db *sqlx.DB) *ItemPostgres {
	return &ItemPostgres{db: db}
}

func (r *ItemPostgres) GetItemByArticle(nmId int) (model.Item, error) {
	var item model.Item
	query := `SELECT * FROM items WHERE nm_id=$1`
	if err := r.db.Get(&item, query, nmId); err != nil {
		return item, fmt.Errorf("error in geting item by article: %s", err)
	}
	return item, nil
}
