package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	kafkaconsumer "github.com/zenmaster911/L0/pkg/kafka_consumer"
	"github.com/zenmaster911/L0/pkg/model"
	"github.com/zenmaster911/L0/pkg/repository"
	"github.com/zenmaster911/L0/pkg/service"
)

type Worker struct {
	services *service.Service
	consumer *kafkaconsumer.KafkaConsumer
	db       *repository.Repository
}

func NewWorker(services *service.Service, consumer *kafkaconsumer.KafkaConsumer, db *repository.Repository) *Worker {
	return &Worker{
		services: services,
		consumer: consumer,
		db:       db,
	}
}

func (w *Worker) StartWorker(ctx context.Context) error {
	//w.consumer.StartReading(ctx)
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

			uid, err := w.services.CreateOrder(&reply)
			if err != nil {
				log.Printf("creating order error: %v\n", err)
				continue
			}
			log.Printf("order with uid %s created", uid)
			// 	log.Printf("message received: Topic=%s Partition=%v Offset=%v Key=%s Value=%s",
			// 		m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
			// }
			time.Sleep(time.Second)
		}
	}
}
