package service

import (
	"context"
	"errors"
	"strings"

	"github.com/adal4ik/crypto-service/internal/repository"
	"github.com/adal4ik/crypto-service/pkg/logger"
)

type CurrencyServiceInterface interface {
	AddCurrency(ctx context.Context, symbol string) error
}

type CurrencyService struct {
	repo   repository.CurrencyRepositoryInterface
	logger logger.Logger
}

func NewCurrencyService(repo repository.CurrencyRepositoryInterface, logger logger.Logger) *CurrencyService {
	return &CurrencyService{
		repo:   repo,
		logger: logger,
	}
}
func (s *CurrencyService) AddCurrency(ctx context.Context, symbol string) error {
	cleanSymbol := strings.TrimSpace(symbol)
	if cleanSymbol == "" {
		return errors.New("currency symbol cannot be empty")
	}

	return s.repo.Add(ctx, cleanSymbol)
}
