package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func getOrderUid(w http.ResponseWriter, r *http.Request) string {
	orderuid := chi.URLParam(r, "order_uid")
	return orderuid
}
