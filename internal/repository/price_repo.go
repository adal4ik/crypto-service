package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/adal4ik/crypto-service/pkg/apperrors"
	"github.com/adal4ik/crypto-service/pkg/logger"
	"go.uber.org/zap"
)

// PriceRepositoryInterface - интерфейс для работы с хранилищем цен.
type PriceRepositoryInterface interface {
	Add(ctx context.Context, symbol string, price float64, timestamp time.Time) *apperrors.AppError
}

type priceRepo struct {
	db     *sql.DB
	logger logger.Logger
}

func NewPriceRepository(db *sql.DB, logger logger.Logger) PriceRepositoryInterface {
	return &priceRepo{db: db, logger: logger}
}

// Add сохраняет цену в базу данных.
// Он должен сначала найти ID валюты по ее символу.
func (r *priceRepo) Add(ctx context.Context, symbol string, price float64, timestamp time.Time) *apperrors.AppError {
	l := r.logger.With(zap.String("symbol", symbol), zap.String("layer", "price_repo"))
	l.Info("Adding price to DB")

	// Находим ID валюты, чтобы связать с ней цену.
	var currencyID string
	err := r.db.QueryRowContext(ctx, "SELECT id FROM tracked_currencies WHERE symbol = $1", symbol).Scan(&currencyID)
	if err != nil {
		if err == sql.ErrNoRows {
			l.Warn("currency not found, skipping price add")
			return nil // Не ошибка, просто валюту уже удалили.
		}
		l.Error("failed to get currency id", zap.Error(err))
		return apperrors.NewInternalServerError("database error", err)
	}

	query := `INSERT INTO price_history (currency_id, price, timestamp) VALUES ($1, $2, $3);`
	_, err = r.db.ExecContext(ctx, query, currencyID, price, timestamp)
	if err != nil {
		l.Error("DB error on price add", zap.Error(err))
		return apperrors.NewInternalServerError("database error", err)
	}

	return nil
}
