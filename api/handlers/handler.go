package handlers

import "github.com/go-chi/chi/v5"

type PaymentHandler struct {
	payServ PaymentServiceInterface
}

type PaymentServiceInterface interface {
	UpdateBalance()
}

func NewPaymentHandler(payServ PaymentServiceInterface) *PaymentHandler {
	return &PaymentHandler{payServ}
}

func (payHandler *PaymentHandler) ApiRoute(r chi.Router) {
	r.Post("/update", payHandler.UpdateBalance)

}
