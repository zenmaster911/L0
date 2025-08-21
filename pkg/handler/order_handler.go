package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/zenmaster911/L0/pkg/model"
)

func (h *Handler) GetOrderByUid(w http.ResponseWriter, r *http.Request) {
	orderUid := getOrderUid(w, r)
	fmt.Println(orderUid)
	if orderUid == "" {
		http.Error(w, "Empty order UID", http.StatusBadRequest)
		return
	}

	orderReply, err := h.services.Order.GetOrderByUid(orderUid)
	if err != nil {
		log.Printf("extracting order data error: %s", err)
		http.Error(w, "extracting order data error", http.StatusInternalServerError)
		return
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
