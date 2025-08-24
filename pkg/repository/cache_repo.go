package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type CachePostgres struct {
	db *sqlx.DB
}

func NewCachePostgres(db *sqlx.DB) *CachePostgres {
	return &CachePostgres{db: db}
}

func (r *CachePostgres) CacheLoad() ([]string, error) {
	var uids []string
	query := `SELECT order_uid FROM orders ORDER BY date_created ASC LIMIT 10;`
	if err := r.db.Select(&uids, query); err != nil {
		return nil, fmt.Errorf("uids receiving error: %v", err)
	}
	return uids, nil
}
