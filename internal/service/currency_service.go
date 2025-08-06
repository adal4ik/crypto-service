package service

import (
	"context"
	"strings"

	"github.com/adal4ik/crypto-service/internal/repository"
	"github.com/adal4ik/crypto-service/pkg/apperrors"
	"github.com/adal4ik/crypto-service/pkg/logger"
	"go.uber.org/zap"
)

type CurrencyServiceInterface interface {
	AddCurrency(ctx context.Context, symbol string) *apperrors.AppError
	RemoveCurrency(ctx context.Context, symbol string) *apperrors.AppError
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
func (s *CurrencyService) AddCurrency(ctx context.Context, symbol string) *apperrors.AppError {
	l := s.logger.With(zap.String("symbol", symbol), zap.String("layer", "service"))
	l.Info("Adding currency")

	cleanSymbol := strings.TrimSpace(symbol)
	if cleanSymbol == "" {
		return apperrors.NewBadRequest("currency symbol cannot be empty", nil)
	}

	return s.repo.Add(ctx, cleanSymbol)
}

func (s *CurrencyService) RemoveCurrency(ctx context.Context, symbol string) *apperrors.AppError {
	l := s.logger.With(zap.String("symbol", symbol), zap.String("layer", "service"))
	l.Info("Removing currency")

	cleanSymbol := strings.TrimSpace(symbol)
	if cleanSymbol == "" {
		return apperrors.NewBadRequest("currency symbol cannot be empty", nil)
	}
	return s.repo.Remove(ctx, cleanSymbol)
}
