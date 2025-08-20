package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/zenmaster911/L0/pkg/model"
)

type CustomerPostgres struct {
	db *sqlx.DB
}

func NewCustomerPostgres(db *sqlx.DB) *CustomerPostgres {
	return &CustomerPostgres{db: db}
}

func (r *CustomerPostgres) GetCustomerByPhone(phone string) (model.Customer, error) {
	var customer model.Customer
	query := `SELECT * FROM customers WHERE phone=$1`
	if err := r.db.Get(&customer, query, phone); err != nil {
		return customer, fmt.Errorf("error in geting customer by it's phone: %s", err)
	}
	return customer, nil
}
