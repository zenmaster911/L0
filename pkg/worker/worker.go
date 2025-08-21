package worker

import (
	kafkaconsumer "github.com/zenmaster911/L0/pkg/kafka_consumer"
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
