package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func getOrderUid(w http.ResponseWriter, r *http.Request) (orderUid string) {
	orderUid = chi.URLParam(r, "order_uid")
	return
}
