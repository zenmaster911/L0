package service

import (
	"github.com/zenmaster911/L0/pkg/model"
	"github.com/zenmaster911/L0/pkg/repository"
)

type ItemsService struct {
	repo repository.Item
}

func NewItemService(repo repository.Item) *ItemsService {
	return &ItemsService{repo: repo}
}

func (s *ItemsService) GetItemByArticle(nmId int) (model.Item, error) {
	return s.repo.GetItemByArticle(nmId)
}
