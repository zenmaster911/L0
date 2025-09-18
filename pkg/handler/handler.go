package handler

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/zenmaster911/L0/pkg/cache"
	"github.com/zenmaster911/L0/pkg/service"
)

type Handler struct {
	services *service.Service
	cache    *cache.RedisCache
}

func NewHandler(services *service.Service, cache *cache.RedisCache) *Handler {
	return &Handler{services: services, cache: cache}
}

func (h *Handler) InitRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	staticDir := os.Getenv("STATIC_DIR")
	fs := http.FileServer(http.Dir(staticDir))
	router.Handle("/*", fs)

	router.Route("/order", func(r chi.Router) {
		r.Post("/", h.CreateOrder)
		r.Get("/{order_uid}", h.GetOrderByUid)
	})
	return router
}
