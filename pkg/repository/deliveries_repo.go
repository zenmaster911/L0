package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/zenmaster911/L0/pkg/model"
)

type DeliveryPostgres struct {
	db *sqlx.DB
}

func NewDeliveryPostgrtes(db *sqlx.DB) *DeliveryPostgres {
	return &DeliveryPostgres{db: db}
}

func (r *DeliveryPostgres) GetCustomerDeliveryByAddress(address, customerUid string) (model.Delivery, int, error) {
	var delivery model.Delivery
	query := `SELECT * FROM deliveries WHERE CONCAT(street,' ',house)=$1 AND customer_uid=$2 RETURNING id`
	if err := r.db.Get(delivery, query, address, customerUid); err != nil {
		return delivery, 0, fmt.Errorf("error in geting delivery by address and customerUid")
	}
	return delivery, delivery.Id, nil
}
