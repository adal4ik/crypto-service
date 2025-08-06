package repository

import (
	"context"
	"database/sql"
	"strings"

	"github.com/adal4ik/crypto-service/pkg/apperrors"
	"github.com/adal4ik/crypto-service/pkg/logger"
	"go.uber.org/zap"
)

type CurrencyRepositoryInterface interface {
	Add(ctx context.Context, symbol string) *apperrors.AppError
	Remove(ctx context.Context, symbol string) *apperrors.AppError
}

type CurrencyRepository struct {
	db     *sql.DB
	logger logger.Logger
}

func NewCurrencyRepository(db *sql.DB, logger logger.Logger) *CurrencyRepository {
	return &CurrencyRepository{
		db:     db,
		logger: logger,
	}
}
func (r *CurrencyRepository) Add(ctx context.Context, symbol string) *apperrors.AppError {
	l := r.logger.With(zap.String("symbol", symbol), zap.String("layer", "repo"))
	l.Info("Adding currency to DB")

	query := `INSERT INTO tracked_currencies (symbol) VALUES ($1) ON CONFLICT (symbol) DO NOTHING;`
	normalizedSymbol := strings.ToUpper(symbol)

	_, err := r.db.ExecContext(ctx, query, normalizedSymbol)
	if err != nil {
		l.Error("DB error on add", zap.Error(err))
		return apperrors.NewInternalServerError("database error", err)
	}
	return nil
}

func (r *CurrencyRepository) Remove(ctx context.Context, symbol string) *apperrors.AppError {
	l := r.logger.With(zap.String("symbol", symbol), zap.String("layer", "repo"))
	l.Info("Removing currency from DB")

	query := `DELETE FROM tracked_currencies WHERE symbol = $1;`
	normalizedSymbol := strings.ToUpper(symbol)

	_, err := r.db.ExecContext(ctx, query, normalizedSymbol)
	if err != nil {
		l.Error("DB error on remove", zap.Error(err))
		return apperrors.NewInternalServerError("database error", err)
	}
	return nil
}
