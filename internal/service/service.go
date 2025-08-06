package service

import (
	"github.com/adal4ik/crypto-service/internal/repository"
	"github.com/adal4ik/crypto-service/pkg/logger"
)

type Service struct {
	CurrencyService *CurrencyService
}

func NewService(repo *repository.Repository, logger logger.Logger) *Service {
	return &Service{
		CurrencyService: NewCurrencyService(repo.CurrencyRepository, logger),
	}
}
