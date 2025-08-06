package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Router(h *Handlers) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/currency", func(r chi.Router) {
		r.Post("/add", h.Currency.CreateCurrency)
		r.Post("/remove", h.Currency.RemoveCurrency)
		r.Get("/price", h.Price.GetPrice)
	})

	return r
}
