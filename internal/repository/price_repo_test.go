package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/adal4ik/crypto-service/pkg/logger"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPriceRepository_Add(t *testing.T) {
	nopLogger := logger.NewNopLogger()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewPriceRepository(db, nopLogger)

		symbol := "BTC"
		currencyID := uuid.New()
		price := decimal.NewFromFloat(65000.50)
		timestamp := time.Now()

		selectQuery := regexp.QuoteMeta(`SELECT id FROM tracked_currencies WHERE symbol = $1`)
		rows := sqlmock.NewRows([]string{"id"}).AddRow(currencyID.String())
		mock.ExpectQuery(selectQuery).WithArgs(symbol).WillReturnRows(rows)

		insertQuery := regexp.QuoteMeta(`INSERT INTO price_history (currency_id, price, timestamp) VALUES ($1, $2, $3);`)
		mock.ExpectExec(insertQuery).WithArgs(currencyID.String(), price, timestamp).WillReturnResult(sqlmock.NewResult(1, 1))

		appErr := repo.Add(ctx, symbol, price, timestamp)

		assert.Nil(t, appErr)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success_currency_not_found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewPriceRepository(db, nopLogger)
		symbol := "UNKNOWN"
		selectQuery := regexp.QuoteMeta(`SELECT id FROM tracked_currencies WHERE symbol = $1`)

		mock.ExpectQuery(selectQuery).WithArgs(symbol).WillReturnError(sql.ErrNoRows)

		appErr := repo.Add(ctx, symbol, decimal.Zero, time.Now())

		assert.Nil(t, appErr)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPriceRepository_GetNearest(t *testing.T) {
	nopLogger := logger.NewNopLogger()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewPriceRepository(db, nopLogger)

		symbol := "BTC"
		timestamp := time.Now()
		expectedPrice := decimal.NewFromFloat(65123.45)
		expectedTimestamp := timestamp.Add(-5 * time.Second)

		query := `^SELECT p.price, p.timestamp FROM price_history p JOIN tracked_currencies c ON p.currency_id = c.id WHERE c.symbol = \$1`

		rows := sqlmock.NewRows([]string{"price", "timestamp"}).AddRow(expectedPrice, expectedTimestamp)
		mock.ExpectQuery(query).WithArgs(symbol, timestamp).WillReturnRows(rows)

		price, foundTime, appErr := repo.GetNearest(ctx, symbol, timestamp)

		assert.Nil(t, appErr)
		assert.Equal(t, expectedPrice, price)
		assert.Equal(t, expectedTimestamp, foundTime)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failure_not_found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewPriceRepository(db, nopLogger)

		symbol := "NONEXISTENT"
		timestamp := time.Now()
		query := `^SELECT p.price, p.timestamp FROM price_history p JOIN tracked_currencies c ON p.currency_id = c.id WHERE c.symbol = \$1`

		mock.ExpectQuery(query).WithArgs(symbol, timestamp).WillReturnError(sql.ErrNoRows)

		_, _, appErr := repo.GetNearest(ctx, symbol, timestamp)

		require.Error(t, appErr)
		assert.Equal(t, "no price history found for this currency", appErr.Message)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
