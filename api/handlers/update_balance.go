package handlers

import (
	"github.com/go-chi/render"
	"net/http"
)

type operation string

const (
	AddOperation  operation = "add"
	TakeOperation operation = "take"
)

type UpdateBalance struct {
	UserID            int    `json:"user_id"`
	Operation         string `json:"operation"`
	ChangingInBalance int    `json:"changing_in_balance"`
}

func (payHand *PaymentHandler) UpdateBalance(w http.ResponseWriter, r *http.Request) {
	var request UpdateBalance
	if err := render.DecodeJSON(r.Body, request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	payHand.payServ.UpdateBalance()
}
