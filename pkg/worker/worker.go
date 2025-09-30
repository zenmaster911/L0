package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/zenmaster911/L0/pkg/cache"
	kafkaconsumer "github.com/zenmaster911/L0/pkg/kafka_consumer"
	"github.com/zenmaster911/L0/pkg/model"
	"github.com/zenmaster911/L0/pkg/repository"
	"github.com/zenmaster911/L0/pkg/service"
)

type Worker struct {
	services *service.Service
	consumer *kafkaconsumer.KafkaConsumer
	db       *repository.Repository
	Cache    *cache.Cache
}

func NewWorker(services *service.Service, consumer *kafkaconsumer.KafkaConsumer, db *repository.Repository, cache *cache.Cache) *Worker {
	return &Worker{
		services: services,
		consumer: consumer,
		db:       db,
		Cache:    cache,
	}
}

func (w *Worker) StartWorker(ctx context.Context) error {

	unsavedOrders := make([]string, 0)
OuterLoop:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shut down process")
			return nil

		default:
			m, err := w.consumer.Reader.ReadMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
					return nil
				}
				fmt.Printf("message reading error: %v\n", err)
				continue
			}
			var reply model.Reply
			for retries := 0; retries < w.consumer.Retries; retries++ {
				if err = json.Unmarshal(m.Value, &reply); err != nil {
					if retries == w.consumer.Retries-1 {
						log.Printf("unmarshaling message from kafka attempt %d failed with error: %v. Message will be sent to DLQ\n", retries+1, err)
						w.consumer.SendToDLQ(ctx, m, err)
						continue OuterLoop
					}
					log.Printf("unmarshaling message from kafka attempt %d failed with error: %v\n", retries+1, err)
					continue
				} else {
					break
				}

			}

			err = w.Cache.AddToCache(ctx, reply)
			if err != nil {
				log.Printf("failed to write %s to cache due to: %v\n", reply.OrderUid, err)
				continue
			}

			if err := w.db.StatusCheck.DBConnectionCheck(); err != nil {
				unsavedOrders = append(unsavedOrders, reply.OrderUid)
				log.Printf("db connection error: %v, message will be saved to cache", err)
				time.Sleep(time.Second)
				continue
			}

			uid, err := w.services.CreateOrder(&reply)
			if err != nil {
				log.Printf("attempt failed to save order %s to db due to: %v\n", reply.OrderUid, err)
				continue
			}
			log.Printf("order with uid %s created\n", uid)

			for _, v := range unsavedOrders {
				reply, err = w.Cache.ReadFromCache(ctx, v)
				if err != nil {
					log.Printf("failed to read %s from cache due to: %v\n", v, err)
					continue
				}
				log.Printf("order with uid %s created\n", v)
			}
			unsavedOrders = slices.Delete(unsavedOrders, 0, len(unsavedOrders))
			time.Sleep(time.Second)
		}
	}
}
