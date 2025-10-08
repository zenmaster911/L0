package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func getOrderUid(w http.ResponseWriter, r *http.Request) string {
	orderuid := chi.URLParam(r, "order_uid")
	return orderuid
}

func sendValidationErrors(w http.ResponseWriter, err error) {
	var errors []string
	for _, err := range err.(validator.ValidationErrors) {
		errors = append(errors, fmt.Sprintf("%s is %s", err.Field(), err.Tag()))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode((map[string]interface{}{
		"error":  "Validation failed",
		"fields": errors,
	})); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
	}
}
