package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/adal4ik/crypto-service/pkg/apperrors"
	"github.com/adal4ik/crypto-service/pkg/logger"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type PriceRepositoryInterface interface {
	Add(ctx context.Context, symbol string, price decimal.Decimal, timestamp time.Time) *apperrors.AppError
	GetNearest(ctx context.Context, symbol string, timestamp time.Time) (decimal.Decimal, time.Time, *apperrors.AppError)
}

type priceRepo struct {
	db     *sql.DB
	logger logger.Logger
}

func NewPriceRepository(db *sql.DB, logger logger.Logger) PriceRepositoryInterface {
	return &priceRepo{db: db, logger: logger}
}

func (r *priceRepo) Add(ctx context.Context, symbol string, price decimal.Decimal, timestamp time.Time) *apperrors.AppError {
	l := r.logger.With(zap.String("symbol", symbol), zap.String("layer", "price_repo"))
	l.Info("Adding price to DB")

	var currencyID string
	err := r.db.QueryRowContext(ctx, "SELECT id FROM tracked_currencies WHERE symbol = $1", symbol).Scan(&currencyID)
	if err != nil {
		if err == sql.ErrNoRows {
			l.Warn("currency not found, skipping price add")
			return nil
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

func (r *priceRepo) GetNearest(ctx context.Context, symbol string, timestamp time.Time) (decimal.Decimal, time.Time, *apperrors.AppError) {
	l := r.logger.With(zap.String("symbol", symbol), zap.Time("timestamp", timestamp), zap.String("layer", "price_repo"))
	l.Info("Getting nearest price from DB")

	query := `
		SELECT p.price, p.timestamp
		FROM price_history p
		JOIN tracked_currencies c ON p.currency_id = c.id
		WHERE c.symbol = $1
		ORDER BY abs(extract(epoch from p.timestamp) - extract(epoch from $2::timestamptz))
		LIMIT 1;
	`
	var price decimal.Decimal
	var foundTimestamp time.Time

	err := r.db.QueryRowContext(ctx, query, symbol, timestamp).Scan(&price, &foundTimestamp)
	if err != nil {
		if err == sql.ErrNoRows {
			l.Warn("no price history found for symbol")
			return decimal.Zero, time.Time{}, apperrors.NewNotFound("no price history found for this currency", err)
		}
		l.Error("DB error on get nearest price", zap.Error(err))
		return decimal.Zero, time.Time{}, apperrors.NewInternalServerError("database error", err)
	}

	return price, foundTimestamp, nil
}
