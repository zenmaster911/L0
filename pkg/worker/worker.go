package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"maps"
	"time"

	kafkaconsumer "github.com/zenmaster911/L0/pkg/kafka_consumer"
	"github.com/zenmaster911/L0/pkg/middleware"
	"github.com/zenmaster911/L0/pkg/model"
	"github.com/zenmaster911/L0/pkg/repository"
	"github.com/zenmaster911/L0/pkg/service"
)

type Worker struct {
	services *service.Service
	consumer *kafkaconsumer.KafkaConsumer
	db       *repository.Repository
	Cache    *middleware.Cache
}

func NewWorker(services *service.Service, consumer *kafkaconsumer.KafkaConsumer, db *repository.Repository, cache *middleware.Cache) *Worker {
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

			if messagesReceived > 100 {
				uid := w.Cache.MessagesList[messagesReceived]
				maps.DeleteFunc(w.Cache.LastMessages, func(k string, v model.Reply) bool {
					return k == uid
				})
				maps.DeleteFunc(w.Cache.MessagesList, func(k int, v string) bool {
					return k == messagesReceived-100
				})
			}
			w.Cache.MessagesList[messagesReceived] = reply.OrderUid
			w.Cache.LastMessages[reply.OrderUid] = reply
			w.Cache.UnreadMessages[reply.OrderUid] = reply

			for range queue {
				uid, err := w.services.CreateOrder(&reply)
				if err != nil {
					log.Printf("creating order error: %v\n", err)
					continue
				}
				log.Printf("order with uid %s created", uid)
			}

			time.Sleep(time.Second)
		}
	}
}
