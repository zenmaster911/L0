package repository

import "github.com/jmoiron/sqlx"

type DBStatusCheck struct {
	db *sqlx.DB
}

func NewStatusCheck(db *sqlx.DB) *DBStatusCheck {
	return &DBStatusCheck{db: db}
}

func (r *DBStatusCheck) DBConnectionCheck() error {
	return r.db.Ping()
}
