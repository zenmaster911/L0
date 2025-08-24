package service

import "github.com/zenmaster911/L0/pkg/repository"

type CacheService struct {
	repo repository.Cache
}

func NewCacheService(repo repository.Cache) *CacheService {
	return &CacheService{
		repo: repo,
	}
}

func (s *CacheService) CacheLoad() ([]string, error) {
	return s.repo.CacheLoad()
}
