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

	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Println("Error loading config:", err)
		//TODO: что произойдет если не удастся загрузить конфиг?
	}

	store := storage.New(cfg.Database)

	newRabbit := rabbit.NewRabbit(exitContext, store)
	err = newRabbit.NewConnection(cfg.RabbitMQ)
	if err != nil {
		log.Println("Connection error:", err)
		//TODO: что произойдет если не удастся подключиться к RabbitMQ? А пользователь хочет сделать платеж?
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
