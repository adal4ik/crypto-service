package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Router(handlers Handlers) http.Handler {
	r := chi.NewRouter()
	r.Post("/currency/add", handlers.CurrencyHandler.CreateCurrency)
	r.Post("/currency/remove", handlers.CurrencyHandler.RemoveCurrency)
	return r
}
