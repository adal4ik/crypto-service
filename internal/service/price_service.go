package service

import (
	"context"
	"time"

	"github.com/adal4ik/crypto-service/internal/repository"
	"github.com/adal4ik/crypto-service/pkg/apperrors"
	"github.com/adal4ik/crypto-service/pkg/logger"
	"go.uber.org/zap"
)

type PriceServiceInterface interface {
	GetNearestPrice(ctx context.Context, symbol string, unixTimestamp int64) (float64, time.Time, *apperrors.AppError)
}

type priceService struct {
	repo   repository.PriceRepositoryInterface
	logger logger.Logger
}

func NewPriceService(repo repository.PriceRepositoryInterface, logger logger.Logger) PriceServiceInterface {
	return &priceService{repo: repo, logger: logger}
}

func (s *priceService) GetNearestPrice(ctx context.Context, symbol string, unixTimestamp int64) (float64, time.Time, *apperrors.AppError) {
	l := s.logger.With(zap.String("symbol", symbol), zap.Int64("timestamp", unixTimestamp), zap.String("layer", "price_service"))
	l.Info("Getting nearest price")

	targetTime := time.Unix(unixTimestamp, 0)

	return s.repo.GetNearest(ctx, symbol, targetTime)
}
