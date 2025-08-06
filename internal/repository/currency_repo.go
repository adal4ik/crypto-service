package repository

import (
	"context"
	"database/sql"
	"strings"

	"github.com/adal4ik/crypto-service/pkg/logger"
)

type CurrencyRepositoryInterface interface {
	Add(ctx context.Context, symbol string) error
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

func (r *CurrencyRepository) Add(ctx context.Context, symbol string) error {

	query := `INSERT INTO tracked_currencies (symbol) VALUES ($1) ON CONFLICT (symbol) DO NOTHING;`

	normalizedSymbol := strings.ToUpper(symbol)

	_, err := r.db.ExecContext(ctx, query, normalizedSymbol)
	return err
}
