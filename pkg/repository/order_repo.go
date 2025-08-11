package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/zenmaster911/L0/pkg/model"
)

type OrderPostgres struct {
	db *sqlx.DB
}

func NewOrderPostgres(db *sqlx.DB) *OrderPostgres {
	return &OrderPostgres{db: db}
}

func (r *OrderPostgres) GetOrderByUid(uid string) (model.Reply, error) {
	var items []model.Item
	var payment model.Payment
	var delivery model.Delivery
	var order model.Order
	var reply model.Reply

	paymentQuery := `SELECT * FROM payments p INNER JOIN orders o ON p.payment_uid=o.payment_id where o.order_uid=$1`
	err := r.db.Get(&payment, paymentQuery, uid)
	if err != nil {
		return reply, fmt.Errorf("error in  reading payment information from db: %s", err)
	}

	itemsQuery := `SELECT i.nm_id, i.name, i.size, i.brand, oi.chrt_id, oi.price, oi.rid, oi.sale, oi.total_price, oi.status FROM order_items oi 
	INNER JOIN items i ON oi.item_id=i.item_id WHERE oi.order_uid=$1`
	err = r.db.Select(&items, itemsQuery, uid)
	if err != nil {
		return reply, fmt.Errorf("error in reading item information form db: %s", err)
	}

	deliveryQuery := `SELECT c.phone, c.email, CONCAT(c.name, ' ', c.surname) AS name,
	d.region, d.zip, d.city, CONCAT(d.street,' ', d.house) AS address
	FROM customers c
	INNER JOIN orders o ON o.customer_id=customer_uid
	INNER JOIN deliveries d ON d.id=o.delivery_id
	WHERE o.order_uid=$1 `
	err = r.db.Get(&delivery, deliveryQuery, uid)
	if err != nil {
		return reply, fmt.Errorf("error in reading delivery information form db: %s", err)
	}

	orderQuery := `SELECT order_uid, track_number, entry_code, internal_signature, shardkey, sm_id, date_created, oof_shard, locale, customer_id FROM orders WHERE order_uid=$1`
	err = r.db.Get(&order, orderQuery, uid)
	if err != nil {
		return reply, fmt.Errorf("error in reading order information form db: %s", err)
	}

	reply = model.Reply{
		OrderUid:          order.OrderUid,
		TrackNumber:       order.TrackNumber,
		Entry:             order.Entry,
		Delivery:          delivery,
		Payment:           payment,
		Items:             items,
		Locale:            order.Locale,
		InternalSignature: order.InternalSignature,
		CustomerId:        order.CustomerId,
		DeliveryService:   order.DeliveryService,
		Shardkey:          order.Shardkey,
		SmId:              order.SmId,
		DateCreated:       order.DateCreated,
		OofShard:          order.OofShard,
	}
	return reply, nil
}
