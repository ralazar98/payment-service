package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"os"
	"payment-service/internal/rabbit"
	"payment-service/internal/storage"
)

func main() {
	r := chi.NewRouter()
	store := storage.New()
	newRabbit := rabbit.NewRabbit(store)
	newRabbit.NewConnection()

	go newRabbit.Updater()

	address := ":" + os.Getenv("PORT")
	http.ListenAndServe(address, r)
}
