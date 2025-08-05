package handler

import (
	"github.com/adal4ik/crypto-service/internal/service"
	"github.com/adal4ik/crypto-service/pkg/logger"
)

type Handlers struct {
}

func NewHandlers(service *service.Service, logger logger.Logger) *Handlers {
	return &Handlers{
		// SubscriptionHandler: NewSubscriptionHandler(service.SubscriptionService, logger),
	}
}
