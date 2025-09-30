package main

import (
	"context"
	"errors"
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
	"github.com/zenmaster911/L0/pkg/cache"
	"github.com/zenmaster911/L0/pkg/handler"
	kafkaconsumer "github.com/zenmaster911/L0/pkg/kafka_consumer"
	"github.com/zenmaster911/L0/pkg/repository"
	"github.com/zenmaster911/L0/pkg/service"
	"github.com/zenmaster911/L0/pkg/worker"
)

var wg sync.WaitGroup

func main() {
	cfg := config.MustLoad()
	dbConn, err := db.NewPostgresDB(cfg.DB)
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}
	defer dbConn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	Repos := repository.NewRepository(dbConn, cfg.Retries)
	Services := service.NewService(Repos)
	Cache := cache.NewCache(cfg.Redis, Services)
	Handlers := handler.NewHandler(Services, Cache)
	KafkaReader := kafkaconsumer.NewKafkaConsumer(cfg.Kafka)
	Worker := worker.NewWorker(Services, KafkaReader, Repos, Cache)

	defer KafkaReader.Close()

	srv := new(server.Server)

	Cache.CacheLoad(ctx, cfg.Cache.CacheStartUpLimit)

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

	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := srv.Run(cfg.App.Port, Handlers.InitRouter()); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("fatal error in server startup: %s", err)
		}
	}()

	log.Printf("server is running on port %s\n", cfg.App.Port)

	go func() {
		defer wg.Done()
		if err := Worker.StartWorker(ctx); err != nil {
			log.Fatalf("kafka worker error: %v", err)
		}

	}()

	wg.Wait()

	log.Print("Server finished it's work")

}
