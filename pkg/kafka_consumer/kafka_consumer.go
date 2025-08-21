package kafkaconsumer

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/zenmaster911/L0/internal/config"
)

type KafkaConsumer struct {
	Reader *kafka.Reader
}

func NewKafkaConsumer(config *config.KafkaConfig) *KafkaConsumer {
	return &KafkaConsumer{
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     []string{config.BrokerAddr},
			GroupID:     config.GroupID,
			Topic:       config.Topic,
			StartOffset: kafka.LastOffset,
		})}
}

func (k *KafkaConsumer) StartReading(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shut down process")
			return
		default:
			m, err := k.Reader.FetchMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
					return
				}
				fmt.Printf("message reading error: %v\n", err)
				continue
			}
			log.Printf("message received: Topic=%s Partition=%v Offset=%v Key=%s Value=%s",
				m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		}
		time.Sleep(time.Second)
	}
}

func (k *KafkaConsumer) Close() {
	if err := k.Reader.Close(); err != nil {
		log.Printf("connection closing error: %s", err)
	}
}
