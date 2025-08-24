package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"maps"
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
	var messagesReceived int
	var queue int
	messagesReceived = len(w.Cache.LastMessages)
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
			if err = json.Unmarshal(m.Value, &reply); err != nil {
				log.Printf("unmarshaling message from kafka error: %v\n", err)
				continue
			}

			messagesReceived++
			queue++

			w.Cache.MessagesList[messagesReceived] = reply.OrderUid
			w.Cache.LastMessages[reply.OrderUid] = reply
			w.Cache.UnreadMessages[queue] = reply

			if len(w.Cache.MessagesList) > 20 {
				keys := slices.Sorted(maps.Keys(w.Cache.MessagesList))
				uid := w.Cache.MessagesList[keys[0]]
				maps.DeleteFunc(w.Cache.LastMessages, func(k string, v model.Reply) bool {
					return k == uid
				})
				maps.DeleteFunc(w.Cache.MessagesList, func(k int, v string) bool {
					return k == keys[0]
				})
			}

			if err := w.db.StatusCheck.DBConnectionCheck(); err != nil {
				log.Printf("db connection error: %v, message will be saved to cache", err)
				continue
			}
			for i := range queue {
				reply := w.Cache.UnreadMessages[i+1]
				uid, err := w.services.CreateOrder(&reply)
				if err != nil {
					maps.DeleteFunc(w.Cache.LastMessages, func(k string, v model.Reply) bool {
						return k == uid
					})
					maps.DeleteFunc(w.Cache.MessagesList, func(k int, v string) bool {
						return v == uid
					})
					log.Printf("creating order error: %v\nMessage with %s uid deleted from the cache", err, uid)
					continue
				}
				log.Printf("order with uid %s created", uid)
			}
			queue = 0
			for _, val := range w.Cache.LastMessages {
				fmt.Println(val, "\n")
			}

			time.Sleep(time.Second)
		}
	}
}
