package model

import (
	"github.com/google/uuid"
)

// type ModelConverter interface {
// 	MarshalBinary() ([]byte, error)
// 	UnmarshalBinary(data []byte) error
// }

type Payment struct {
	ID           uuid.UUID `json:"payment_uid" db:"payment_uid"`
	Transaction  string    `json:"transaction" validate:"required,min=8,max=24" db:"transaction"`
	RequestId    string    `json:"request_id" db:"request_id"`
	Currency     string    `json:"currency" validate:"required,min=2,max=4" db:"currency"`
	Provider     string    `json:"provider" validate:"required,min=3,max=15" db:"provider"`
	Amount       int       `json:"amount" validate:"required" db:"amount"`
	PaymentDt    int       `json:"payment_dt" validate:"required,min=8,max=15" db:"payment_dt"`
	Bank         string    `json:"bank" validate:"required,min=1,max=20"  db:"bank"`
	DeliveryCost int       `json:"delivery_cost" validate:"required" db:"delivery_cost"`
	GoodsTotal   int       `json:"goods_total"  validate:"required" db:"goods_total"`
	CustomFee    int       `json:"custom_fee"   db:"custom_fee"`
}

type Item struct {
	Id    int    `json:"id" db:"id"`
	NmId  int    `json:"nm_id" db:"nm_id"`
	Name  string `json:"name" db:"name"`
	Size  string `json:"size" db:"size"`
	Brand string `json:"brand" db:"brand"`
}

type DeliveryItem struct {
	ChrtId      int    `json:"chrt_id"  db:"chrt_id"`
	TrackNumber string `json:"track_number" validate:"required,min=10,max=20"  db:"track_number"`
	Price       int    `json:"price" validate:"required" db:"price"`
	Rid         string `json:"rid" validate:"required,min=8,max=30" db:"rid"`
	Name        string `json:"name" validate:"required,min=2,max=20" db:"name"`
	Sale        int    `json:"sale" db:"sale"`
	Size        string `json:"size" db:"size"`
	TotalPrice  int    `json:"total_price" db:"total_price"`
	NmId        int    `json:"nm_id" db:"nm_id"`
	Brand       string `json:"brand" db:"brand"`
	Status      int    `json:"status" db:"status"`
}

type Delivery struct {
	Id      int    `json:"id" db:"id"`
	Name    string `json:"name"  validate:"required,min=5,max=40" db:"name"`
	Phone   string `json:"phone"  validate:"required,min=5,max=14" db:"phone"`
	Zip     string `json:"zip" validate:"required,min=3,max=10" db:"zip"`
	City    string `json:"city" validate:"required,min=3,max=30" db:"city"`
	Address string `json:"address" validate:"required,min=3,max=30" db:"address"`
	Region  string `json:"region" validate:"required,min=3,max=30" db:"region"`
	Email   string `json:"email" validate:"required,min=6,max=50" db:"email"`
}

type Customer struct {
	Uid     uuid.UUID `json:"customer_uid" validate:"required" db:"customer_uid"`
	Name    string    `json:"name" validate:"required,min=2,max=20" db:"name"`
	Surname string    `json:"surname" validate:"required,min=2,max=20" db:"surname"`
	Phone   string    `json:"phone" validate:"required,min=5,max=14" db:"phone"`
	Email   string    `json:"email" validate:"required,min=6,max=50" db:"email"`
}

type Order struct {
	OrderUid          string `json:"order_uid" validate:"required,min=15,max=20" db:"order_uid"`
	TrackNumber       string `json:"track_number" validate:"required,min=10,max=20" db:"track_number"`
	Entry             string `json:"entry" validate:"required,min=4,max=10" db:"entry_code"`
	Locale            string `json:"locale" validate:"required,min=2,max=4" db:"locale"`
	InternalSignature string `json:"internal_signature" db:"internal_signature"`
	CustomerId        string `json:"customer_id" db:"customer_id"`
	DeliveryService   string `json:"delivery_service" db:"delivery_service"`
	Shardkey          string `json:"shardkey" db:"shardkey"`
	SmId              int    `json:"sm_id" db:"sm_id"`
	DateCreated       string `json:"date_created" db:"date_created"`
	OofShard          string `json:"oof_shard" db:"oof_shard"`
}

type Reply struct {
	OrderUid          string         `json:"order_uid" validate:"required" db:"order_uid"`
	TrackNumber       string         `json:"track_number" validate:"required,min=10,max=20" db:"track_number"`
	Entry             string         `json:"entry" validate:"required,min=3,max=10" db:"entry_code"`
	Delivery          Delivery       `json:"delivery" validate:"required"  db:"delivery"`
	Payment           Payment        `json:"payment" validate:"required" db:"payment"`
	Items             []DeliveryItem `json:"items" validate:"required,min=1" db:"items"`
	Locale            string         `json:"locale" db:"locale"`
	InternalSignature string         `json:"internal_signature"  db:"internal_signature"`
	CustomerId        string         `json:"customer_id" validate:"required" db:"customer_id"`
	DeliveryService   string         `json:"delivery_service" validate:"required" db:"delivery_service"`
	Shardkey          string         `json:"shardkey" validate:"required" db:"shardkey"`
	SmId              int            `json:"sm_id" validate:"required" db:"sm_id"`
	DateCreated       string         `json:"date_created" validate:"required" db:"date_created"`
	OofShard          string         `json:"oof_shard" validate:"required" db:"oof_shard"`
}
