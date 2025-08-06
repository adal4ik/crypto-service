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
	GetAll(ctx context.Context) ([]string, *apperrors.AppError)
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

func (r *CurrencyRepository) GetAll(ctx context.Context) ([]string, *apperrors.AppError) {
	l := r.logger.With(zap.String("layer", "repo"))
	l.Info("Getting all tracked currencies from DB")

	query := `SELECT symbol FROM tracked_currencies;`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		l.Error("DB error on get all", zap.Error(err))
		return nil, apperrors.NewInternalServerError("database error", err)
	}
	defer rows.Close()

	var symbols []string
	for rows.Next() {
		var symbol string
		if err := rows.Scan(&symbol); err != nil {
			l.Error("DB error on scan symbol", zap.Error(err))
			return nil, apperrors.NewInternalServerError("database error", err)
		}
		symbols = append(symbols, symbol)
	}

	return symbols, nil
}
