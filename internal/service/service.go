package service

import (
	"github.com/adal4ik/crypto-service/internal/config"
	"github.com/adal4ik/crypto-service/internal/repository"
	"github.com/adal4ik/crypto-service/pkg/logger"
)

// Service - это контейнер для всех наших сервисов.
type Service struct {
	Currency       CurrencyServiceInterface // Используем интерфейс для гибкости
	PriceCollector *PriceCollector
}

// NewService - главный конструктор для слоя сервисов.
func NewService(repo *repository.Repository, logger logger.Logger, cfg *config.Config) *Service {
	return &Service{
		// Создаем сервис для валют.
		// Передаем ему только репозиторий, т.к. логгер ему не нужен.
		Currency: NewCurrencyService(repo.CurrencyRepository, logger),

		// Создаем сервис-сборщик цен.
		// Ему нужны оба репозитория, логгер и его часть конфига.
		PriceCollector: NewPriceCollector(repo.CurrencyRepository, repo.Price, logger, cfg.Collector),
	}
}
