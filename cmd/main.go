package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/zenmaster911/L0/internal/config"
	"github.com/zenmaster911/L0/internal/db"
	"github.com/zenmaster911/L0/internal/server"
	"github.com/zenmaster911/L0/pkg/handler"
	kafkaconsumer "github.com/zenmaster911/L0/pkg/kafka_consumer"
	"github.com/zenmaster911/L0/pkg/repository"
	"github.com/zenmaster911/L0/pkg/service"
	"github.com/zenmaster911/L0/pkg/worker"
)

var wg sync.WaitGroup

func main() {
	cfg, kafkacfg := config.MustLoad()
	dbConn, err := db.NewPostgresDB(cfg.DB)
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}
	defer dbConn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// fmt.Println(os.Getenv("KAFKA_BROKER_ADDR"))
	// kafkacfg := &kafkaconsumer.Config{
	// 	BrokerAddr: os.Getenv("KAFKA_BROKER_ADDR"),
	// 	GroupID:    os.Getenv("GROUP_ID"),
	// 	Topic:      os.Getenv("TOPIC_NAME"),
	// }

	// kafkaReader := kafkaconsumer.NewKafkaConsumer(kafkacfg)
	// if err != nil {
	// 	log.Fatalf("error in creating Kafka consumer: %v", err)
	// }
	// defer kafkaReader.Close()

	Repos := repository.NewRepository(dbConn)
	Services := service.NewService(Repos)
	Handlers := handler.NewHandler(Services)
	KafkaReader := kafkaconsumer.NewKafkaConsumer(kafkacfg)
	Worker := worker.NewWorker(Services, KafkaReader, Repos)

	defer KafkaReader.Close()

	srv := new(server.Server)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	go func() {
		<-sig
		log.Println("Received termination signal")
		cancel()
		if err := srv.HttpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
	}()

	log.Print("server started sucessfully, what do you wish,Master?")

	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := srv.Run(cfg.App.Port, Handlers.InitRouter()); err != nil && err != http.ErrServerClosed {
			log.Fatalf("fatal error in server startup: %s", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := Worker.StartWorker(ctx); err != nil {
			log.Fatalf("kafka worker error: %v", err)
		}
		fmt.Printf("kafkaReader: %v\n", Worker)
	}()

	wg.Wait()

	log.Print("Server finished it's work")

}
