package handler

import "github.com/zenmaster911/L0/pkg/service"

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}
