package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/adal4ik/crypto-service/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCurrencyRepository_Add(t *testing.T) {
	nopLogger := logger.NewNopLogger()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewCurrencyRepository(db, nopLogger)
		symbol := "BTC"
		query := regexp.QuoteMeta(`INSERT INTO tracked_currencies (symbol) VALUES ($1) ON CONFLICT (symbol) DO NOTHING;`)

		mock.ExpectExec(query).WithArgs(symbol).WillReturnResult(sqlmock.NewResult(1, 1))

		appErr := repo.Add(ctx, symbol)

		assert.Nil(t, appErr)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failure_db_error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewCurrencyRepository(db, nopLogger)
		symbol := "ETH"
		query := regexp.QuoteMeta(`INSERT INTO tracked_currencies (symbol) VALUES ($1) ON CONFLICT (symbol) DO NOTHING;`)
		dbError := errors.New("db is down")

		mock.ExpectExec(query).WithArgs(symbol).WillReturnError(dbError)

		appErr := repo.Add(ctx, symbol)

		require.Error(t, appErr)
		assert.Equal(t, "database error", appErr.Message)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCurrencyRepository_Remove(t *testing.T) {
	nopLogger := logger.NewNopLogger()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewCurrencyRepository(db, nopLogger)
		symbol := "BTC"
		query := regexp.QuoteMeta(`DELETE FROM tracked_currencies WHERE symbol = $1;`)

		mock.ExpectExec(query).WithArgs(symbol).WillReturnResult(sqlmock.NewResult(0, 1))

		appErr := repo.Remove(ctx, symbol)

		assert.Nil(t, appErr)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCurrencyRepository_GetAll(t *testing.T) {
	nopLogger := logger.NewNopLogger()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewCurrencyRepository(db, nopLogger)
		query := regexp.QuoteMeta(`SELECT symbol FROM tracked_currencies;`)

		expectedSymbols := []string{"BTC", "ETH"}
		rows := sqlmock.NewRows([]string{"symbol"}).
			AddRow(expectedSymbols[0]).
			AddRow(expectedSymbols[1])

		mock.ExpectQuery(query).WillReturnRows(rows)

		symbols, appErr := repo.GetAll(ctx)

		assert.Nil(t, appErr)
		assert.Equal(t, expectedSymbols, symbols)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failure_db_error_on_query", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewCurrencyRepository(db, nopLogger)
		query := regexp.QuoteMeta(`SELECT symbol FROM tracked_currencies;`)
		dbError := errors.New("query failed")

		mock.ExpectQuery(query).WillReturnError(dbError)

		_, appErr := repo.GetAll(ctx)

		require.Error(t, appErr)
		assert.Equal(t, "database error", appErr.Message)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
