package cache

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/zenmaster911/L0/internal/config"
	"github.com/zenmaster911/L0/pkg/model"
	"github.com/zenmaster911/L0/pkg/service"
)

type RedisCacheInterface interface {
	CacheLoad(ctx context.Context, limit int) error
	AddToCache(ctx context.Context, order model.Reply) error
	ReadFromCache(ctx context.Context, uid string) (reply model.Reply, err error)
}

type RedisCache struct {
	service *service.Service
	Client  *redis.Client
}

func NewRedisCache(cfg *config.RedisConfig, service *service.Service) *RedisCache {
	return &RedisCache{
		service: service,
		Client: redis.NewClient(&redis.Options{
			Addr:       cfg.Addr,
			Password:   cfg.Password,
			DB:         cfg.DB,
			Username:   cfg.User,
			MaxRetries: cfg.MaxRetries,
			// DialTimeout:  cfg.DialTimeout * time.Second,
			// ReadTimeout:  cfg.Timeout * time.Second,
			// WriteTimeout: cfg.Timeout * time.Second,
		}),
	}
}

func (rc *RedisCache) CacheLoad(ctx context.Context, limit int) error {
	uids, err := rc.service.Cache.CacheLoad(limit)
	if err != nil {
		fmt.Printf("errr, %v", err)
		return fmt.Errorf("loading order uids to cache error: %v", err)
	}
	for _, v := range uids {
		order, err := rc.service.GetOrderByUid(v)
		if err != nil {
			log.Printf("%s failed to upload to cache due to: %v", v, err)
		}
		fmt.Println(v)
		err = rc.Client.Set(ctx, v, order, 0).Err()
		if err != nil {
			log.Printf("%s failed to upload to cache due to: %v", v, err)
		}
	}
	return nil
}

func (rc *RedisCache) AddToCache(ctx context.Context, order model.Reply) error {
	err := rc.Client.Set(ctx, order.OrderUid, order, 0).Err()
	if err != nil {
		log.Printf("%s failed to upload to cache due to: %v", order.OrderUid, err)
	}
	return err
}

func (rc *RedisCache) ReadFromCache(ctx context.Context, uid string) (reply model.Reply, err error) {
	err = rc.Client.Get(context.Background(), uid).Scan(&reply)
	if err != nil {
		{
			return reply, fmt.Errorf("read from cache error: %v", err)
		}
	}
	fmt.Println(reply.OrderUid)
	return reply, nil
}
