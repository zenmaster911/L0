package kafkaconsumer

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/zenmaster911/L0/internal/config"
)

//go:generate minimock -i github.com/zenmaster911/L0/pkg/kafka_consumer.KafkaReader -o ./kafka_mocks -s _mock.go

type KafkaReader interface {
	FetchMessage(ctx context.Context) (kafka.Message, error)
	ReadMessage(ctx context.Context) (kafka.Message, error)
	Close() error
}

type KafkaWriter interface {
	WriteMessages(ctx context.Context, messages ...kafka.Message) error
	Close() error
}

type KafkaConsumer struct {
	Reader  KafkaReader
	Writer  KafkaWriter
	Retries int
}

func NewKafkaConsumer(config *config.KafkaConfig) *KafkaConsumer {
	return &KafkaConsumer{
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     []string{config.BrokerAddr},
			GroupID:     config.GroupID,
			Topic:       config.Topic,
			StartOffset: kafka.LastOffset,
		}),
		Writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:  []string{config.BrokerAddr},
			Topic:    config.DLQTopic,
			Balancer: &kafka.LeastBytes{},
		}),
		Retries: config.MaxRetries,
	}

}

func (k *KafkaConsumer) SendToDLQ(ctx context.Context, m kafka.Message, err error) {
	dlqMessage := kafka.Message{
		Value: m.Value,
		Headers: []kafka.Header{
			{Key: "original_topic", Value: []byte(m.Topic)},
			{Key: "original_partition", Value: []byte(strconv.Itoa(m.Partition))},
			{Key: "original_offset", Value: []byte(strconv.Itoa(int(m.Offset)))},
			{Key: "error", Value: []byte(err.Error())},
		},
	}
	if err := k.Writer.WriteMessages(ctx, dlqMessage); err != nil {
		log.Printf("message sending to DLQ error: %v\n", err)
	}
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
		log.Printf("connection reader closing error: %s", err)
	}
	if err := k.Writer.Close(); err != nil {

		log.Printf("connection DLQ writer closing error: %s", err)

	}
}
