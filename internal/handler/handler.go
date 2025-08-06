package handler

import (
	"github.com/adal4ik/crypto-service/internal/service"
	"github.com/adal4ik/crypto-service/pkg/logger"
)

type Handlers struct {
	Currency *CurrencyHandler
	Price    *PriceHandler
}

func NewHandlers(s *service.Service, logger logger.Logger) *Handlers {
	currencyHandler := NewCurrencyHandler(s.Currency, logger)

	return &Handlers{
		Currency: currencyHandler,
		Price:    NewPriceHandler(s.Price, logger, currencyHandler.handleError),
	}
}
