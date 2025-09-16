package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/zenmaster911/L0/pkg/model"
)

var validate = validator.New()

func (h *Handler) GetOrderByUid(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	orderUid := getOrderUid(w, r)
	var orderReply model.Reply
	fmt.Println(orderUid)
	if orderUid == "" {
		http.Error(w, "Empty order UID", http.StatusBadRequest)
		return
	}
	if val, exist := h.cache.LastMessages[orderUid]; exist {
		orderReply = val
	} else {
		exist, err := h.services.Order.CheckOrderExists(orderUid)
		if err != nil {
			log.Printf("order existance check error: %v", err)
			http.Error(w, "extracting order data error", http.StatusInternalServerError)
			return
		}
		if !exist {
			http.Error(w, "order with this uid doexn't exist", http.StatusNotFound)
			return
		}
		val, err := h.services.Order.GetOrderByUid(orderUid)
		if err != nil {
			log.Printf("extracting order data error: %s", err)
			http.Error(w, "extracting order data error", http.StatusInternalServerError)
			return
		}
		orderReply = val
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orderReply)
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var input model.Reply

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusInternalServerError)
	}

	if err := validate.Struct(input); err != nil {
		sendValidationErrors(w, err)
		return
	}

	uid, err := h.services.Order.CreateOrder(&input)
	if err != nil {
		log.Printf("creating order error: %s", err)
		http.Error(w, "creating order error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"order_uid": uid,
	})
}
