package routes

import (
	"github.com/go-chi/chi/v5"
	"payment-service/api/handlers"
)

func (payHandler *handlers.PaymentHandler) ApiRoute(r chi.Router) {
	r.Post("/update", payHandler.UpdateBalance)
}
