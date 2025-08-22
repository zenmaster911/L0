package cache

import "github.com/zenmaster911/L0/pkg/model"

type Cache struct {
	LastMessages   map[string]model.Reply
	UnreadMessages map[string]model.Reply
	MessagesList   map[int]string
}

func NewCache() *Cache {
	return &Cache{
		LastMessages:   make(map[string]model.Reply),
		UnreadMessages: make(map[string]model.Reply),
		MessagesList:   make(map[int]string),
	}
}
