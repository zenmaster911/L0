package repository

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/zenmaster911/L0/internal/config"
	"github.com/zenmaster911/L0/pkg/model"
)

type OrderPostgres struct {
	db         *sqlx.DB
	maxRetries int
	retryDelay int
}

func NewOrderPostgres(db *sqlx.DB, cfg *config.DBRetriesConfig) *OrderPostgres {
	return &OrderPostgres{
		db:         db,
		maxRetries: cfg.MaxRetries,
		retryDelay: cfg.RetryDelay,
	}
}

func (r *OrderPostgres) CreateOrder(input *model.Reply) (uid string, err error) {
	exists := false
	var itemId, deliveryId int
	itemIds := make([]int, 0)
	var customerUid string
	for i := 0; i < r.maxRetries; i++ {
		tx, err := r.db.Begin()
		if err != nil {
			return "", fmt.Errorf("DB starting transaction error: %s", err)
		}
		err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM orders WHERE order_uid=$1)`, input.OrderUid).Scan(&exists)
		if err != nil {
			tx.Rollback()
			if isTransietError(err) && i < r.maxRetries-1 {
				log.Printf("retrying %d time with error: %s", i+1, err)
				time.Sleep(time.Duration(r.retryDelay) * time.Millisecond)
				continue
			}

			return "", fmt.Errorf("order uid %s check failed with error: %w", input.OrderUid, err)
		}
		if exists {
			tx.Rollback()
			if isTransietError(err) && i < r.maxRetries-1 {
				log.Printf("retrying %d time with error: %s", i+1, err)
				time.Sleep(time.Duration(r.retryDelay) * time.Millisecond)
				continue
			}
			return "", fmt.Errorf("order %s already exists", input.OrderUid)
		}
		// так как на данный момент база данных пустая, исхожу из того что в запросе упомянуты только валидные товары и, при необходимости, добавляю их в бд.
		// в случае работы с реальной базой данных, в случае отсутствия id товара в базе транзакция была бы остановлена.
		// аналочгично в дальнейшем будут реализованы таблицы customers и deliveries
		for _, v := range input.Items {
			err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM items WHERE nm_id=$1)`, v.NmId).Scan(&exists)
			if err != nil {
				tx.Rollback()
				if isTransietError(err) && i < r.maxRetries-1 {
					log.Printf("retrying %d time with error: %s", i+1, err)
					time.Sleep(time.Duration(r.retryDelay) * time.Millisecond)
					continue
				}
				return "", fmt.Errorf("item %d check failed with error: %w", v.NmId, err)
			}
			if !exists {
				itemQuery := `INSERT INTO items (nm_id,name,size,brand)
			VALUES($1,$2,$3,$4)
			RETURNING item_id`
				if err := tx.QueryRow(itemQuery, v.NmId, v.Name, v.Size, v.Brand).Scan(&itemId); err != nil {
					tx.Rollback()
					if isTransietError(err) && i < r.maxRetries-1 {
						log.Printf("retrying %d time with error: %s", i+1, err)
						time.Sleep(time.Duration(r.retryDelay) * time.Millisecond)
						continue
					}
					return "", fmt.Errorf("error inserting item into database %w", err)
				}
			} else {
				err := tx.QueryRow(`Select item_id FROM items WHERE nm_id=$1`, v.NmId).Scan(&itemId)
				if err != nil {
					tx.Rollback()
					if isTransietError(err) && i < r.maxRetries-1 {
						log.Printf("retrying %d time with error: %s", i+1, err)
						time.Sleep(time.Duration(r.retryDelay) * time.Millisecond)
						continue
					}
					return "", fmt.Errorf("receiving existing item id from database error %w", err)
				}
			}
			itemIds = append(itemIds, itemId)
		}

		err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM customers WHERE phone=$1)`, input.Delivery.Phone).Scan(&exists)
		if err != nil {
			tx.Rollback()
			if isTransietError(err) && i < r.maxRetries-1 {
				log.Printf("retrying %d time with error: %s", i+1, err)
				time.Sleep(time.Duration(r.retryDelay) * time.Millisecond)
				continue
			}
			return "", fmt.Errorf("customer's phone %s check failed with error: %w", input.Delivery.Phone, err)
		}
		if !exists {
			customerUid = uuid.NewString()
			customerQuery := `INSERT INTO customers (customer_uid, name, surname, phone, email)
		VALUES ($1,$2,$3,$4,$5)`
			name := strings.Split(input.Delivery.Name, " ")
			if _, err = tx.Exec(customerQuery, customerUid, name[0], name[1], input.Delivery.Phone, input.Delivery.Email); err != nil {
				tx.Rollback()
				if isTransietError(err) && i < r.maxRetries-1 {
					log.Printf("retrying %d time with error: %s", i+1, err)
					time.Sleep(time.Duration(r.retryDelay) * time.Millisecond)
					continue
				}
				return "", fmt.Errorf("error inserting customer into database %w", err)
			}
		} else {
			err = tx.QueryRow(`SELECT customer_uid FROM customers WHERE phone=$1`, input.Delivery.Phone).Scan(&customerUid)
			if err != nil {
				tx.Rollback()
				if isTransietError(err) && i < r.maxRetries-1 {
					log.Printf("retrying %d time with error: %s", i+1, err)
					time.Sleep(time.Duration(r.retryDelay) * time.Millisecond)
					continue
				}
				return "", fmt.Errorf(" receiving customer's  uid failed with error: %w", err)
			}
		}

		err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM deliveries WHERE CONCAT(street,' ',house)=$1 AND customer_uid=$2)`, input.Delivery.Address, customerUid).Scan(&exists)
		if err != nil {
			tx.Rollback()
			if isTransietError(err) && i < r.maxRetries-1 {
				log.Printf("retrying %d time with error: %s", i+1, err)
				time.Sleep(time.Duration(r.retryDelay) * time.Millisecond)
				continue
			}
			return "", fmt.Errorf("delivery address %s check failed with error: %w", input.Delivery.Phone, err)
		}
		if !exists {
			deliveryQuery := `INSERT INTO deliveries (region, zip, city, street, house, customer_uid)
		VALUES ($1,$2,$3,$4,$5,$6)
		Returning id`
			address := strings.Split(input.Delivery.Address, " ")
			if err = tx.QueryRow(deliveryQuery, input.Delivery.Region, input.Delivery.Zip, input.Delivery.City, address[0], address[1], customerUid).Scan(&deliveryId); err != nil {
				tx.Rollback()
				if isTransietError(err) && i < r.maxRetries-1 {
					log.Printf("retrying %d time with error: %s", i+1, err)
					time.Sleep(time.Duration(r.retryDelay) * time.Millisecond)
					continue
				}
				return "", fmt.Errorf("error inserting delivery into database %w", err)
			}
		} else {
			err = tx.QueryRow(`SELECT delivery_id FROM deliveries WHERE CONCAT(street,' ',house)=$1 AND customer_uid=$2)`, input.Delivery.Address, customerUid).Scan(&deliveryId)
			if err != nil {
				tx.Rollback()
				if isTransietError(err) && i < r.maxRetries-1 {
					log.Printf("retrying %d time with error: %s", i+1, err)
					time.Sleep(time.Duration(r.retryDelay) * time.Millisecond)
					continue
				}
				return "", fmt.Errorf(" receiving delivery id failed with error: %w", err)
			}
		}

		paymentUid := uuid.NewString()
		paymentQuery := `INSERT INTO payments (payment_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
	 VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
		if _, err = tx.Exec(paymentQuery, paymentUid, input.Payment.Transaction, input.Payment.RequestId, input.Payment.Currency, input.Payment.Provider,
			input.Payment.Amount, input.Payment.PaymentDt, input.Payment.Bank, input.Payment.DeliveryCost, input.Payment.GoodsTotal, input.Payment.CustomFee); err != nil {
			tx.Rollback()
			if isTransietError(err) && i < r.maxRetries-1 {
				log.Printf("retrying %d time with error: %s", i+1, err)
				time.Sleep(time.Duration(r.retryDelay) * time.Millisecond)
				continue
			}
			return "", fmt.Errorf("error in inserting payment into database %w", err)
		}

		orderQuery := `INSERT INTO orders (order_uid, track_number, entry_code, internal_signature, shardkey, sm_id, date_created,oof_shard,locale,customer_id,delivery_service,delivery_id,payment_id)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`
		if _, err = tx.Exec(orderQuery, input.OrderUid, input.TrackNumber, input.Entry, input.InternalSignature, input.Shardkey, input.SmId, input.DateCreated, input.OofShard,
			input.Locale, customerUid, input.DeliveryService, deliveryId, paymentUid); err != nil {
			tx.Rollback()
			if isTransietError(err) && i < r.maxRetries-1 {
				log.Printf("retrying %d time with error: %s", i+1, err)
				time.Sleep(time.Duration(r.retryDelay) * time.Millisecond)
				continue
			}
			return "", fmt.Errorf("error in inserting order into database %w", err)
		}

		for i, v := range input.Items {
			orderItemsQuery := `INSERT INTO order_items (item_id, order_uid, chrt_id, price, rid, sale, total_price, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
			if _, err := tx.Exec(orderItemsQuery, itemIds[i], input.OrderUid, v.ChrtId, v.Price, v.Rid, v.Sale, v.TotalPrice, v.Status); err != nil {
				tx.Rollback()
				if isTransietError(err) && i < r.maxRetries-1 {
					log.Printf("retrying %d time with error: %s", i+1, err)
					time.Sleep(time.Duration(r.retryDelay) * time.Millisecond)
					continue
				}
				return "", fmt.Errorf("error in inserting order items into database %w", err)
			}
		}

		if err := tx.Commit(); err != nil {
			tx.Rollback()
			if isTransietError(err) && i < r.maxRetries-1 {
				log.Printf("retrying %d time with error: %s", i+1, err)
				time.Sleep(time.Duration(r.retryDelay) * time.Millisecond)
				continue
			}
			return input.OrderUid, fmt.Errorf("error in committing transaction %w", err)
		}

		return input.OrderUid, nil
	}
	return input.OrderUid, fmt.Errorf("failed to create order: max retries exceeded")
}

func (r *OrderPostgres) CheckOrderExists(uid string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM orders WHERE order_uid=$1)`, uid).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("existance check error: %v", err)
	}
	if exists {
		return true, nil
	}
	return false, nil
}

func (r *OrderPostgres) GetOrderByUid(uid string) (model.Reply, error) {
	var items []model.DeliveryItem
	var payment model.Payment
	var delivery model.Delivery
	var order model.Order
	var reply model.Reply

	paymentQuery := `SELECT p.* FROM payments p INNER JOIN orders o ON p.payment_uid=o.payment_id WHERE o.order_uid=$1`
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

func isTransietError(err error) bool {
	if strings.Contains(err.Error(), "deadlock") || strings.Contains(err.Error(), "serialization failure") {
		return true
	}
	if _, ok := err.(net.Error); ok {
		return true
	}
	return false
}
