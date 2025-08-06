package handler

import (
	"github.com/adal4ik/crypto-service/internal/service"
	"github.com/adal4ik/crypto-service/pkg/logger"
)

type Handlers struct {
	CurrencyHandler *CurrencyHandler
}

func NewHandlers(service *service.Service, logger logger.Logger) *Handlers {
	return &Handlers{
		CurrencyHandler: NewCurrencyHandler(service.CurrencyService, logger),
	}
}
