package service

import (
	"github.com/adal4ik/crypto-service/internal/config"
	"github.com/adal4ik/crypto-service/internal/repository"
	"github.com/adal4ik/crypto-service/pkg/logger"
)

type Service struct {
	Currency       CurrencyServiceInterface
	PriceCollector *PriceCollector
	Price          PriceServiceInterface
}

func NewService(repo *repository.Repository, logger logger.Logger, cfg *config.Config) *Service {
	return &Service{
		Currency:       NewCurrencyService(repo.CurrencyRepository, logger),
		PriceCollector: NewPriceCollector(repo.CurrencyRepository, repo.Price, logger, cfg.Collector),
		Price:          NewPriceService(repo.Price, logger),
	}
}
