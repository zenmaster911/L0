package cache

import (
	"fmt"

	"github.com/zenmaster911/L0/pkg/model"
	"github.com/zenmaster911/L0/pkg/service"
)

type Cache struct {
	service        *service.Service
	LastMessages   map[string]model.Reply
	UnreadMessages map[int]model.Reply
	MessagesList   map[int]string
}

func NewCache(service *service.Service) *Cache {
	return &Cache{
		service:        service,
		LastMessages:   make(map[string]model.Reply),
		UnreadMessages: make(map[int]model.Reply),
		MessagesList:   make(map[int]string),
	}
}

func (c *Cache) CacheLoad() error {
	uids, err := c.service.Cache.CacheLoad()
	if err != nil {
		return fmt.Errorf("loading order uids to cache error: %v", err)
	}

	for _, v := range uids {
		c.LastMessages[v], err = c.service.GetOrderByUid(v)
		if err != nil {
			return fmt.Errorf("loading message to cache error: %v", err)
		}
	}
	return nil
}
