package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"payment-service/configs"
	"payment-service/internal/rabbit"
	"payment-service/internal/storage"
	"syscall"
	"time"
)

func main() {

	exitContext, cancel := context.WithCancel(context.Background())

	cfg, err := configs.LoadConfig("configs")
	if err != nil {
		log.Println("Error loading config:", err)
		//TODO: что произойдет если не удастся загрузить конфиг?
	}

	store, err := storage.New(cfg.Database)
	if err != nil {
		log.Println("Error connecting to database:", err)
	}

	_, err = rabbit.NewRabbit(exitContext, cfg.RabbitMQ, store)
	if err != nil {
		log.Println("Error creating new rabbit:", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
		log.Println("Canceling context")
		cancel()
	}
	time.Sleep(1 * time.Second)
	log.Println("Shutting down server...")

}
