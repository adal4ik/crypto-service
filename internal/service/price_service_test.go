package service

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/adal4ik/crypto-service/internal/repository/mocks"
	"github.com/adal4ik/crypto-service/pkg/apperrors"
	"github.com/adal4ik/crypto-service/pkg/logger"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPriceService_GetNearestPrice(t *testing.T) {
	nopLogger := logger.NewNopLogger()
	ctx := context.Background()

	t.Run("success_get_price", func(t *testing.T) {
		mockRepo := mocks.NewPriceRepositoryInterface(t)

		symbol := "BTC"
		unixTimestamp := int64(1672531200)
		expectedTime := time.Unix(unixTimestamp, 0)

		expectedPrice := decimal.NewFromFloat(65000.0)
		expectedFoundTime := expectedTime.Add(-10 * time.Second)

		mockRepo.On("GetNearest", ctx, symbol, expectedTime).
			Return(expectedPrice, expectedFoundTime, nil)

		priceService := NewPriceService(mockRepo, nopLogger)

		price, foundTime, appErr := priceService.GetNearestPrice(ctx, symbol, unixTimestamp)

		assert.Nil(t, appErr)
		assert.Equal(t, expectedPrice, price)
		assert.Equal(t, expectedFoundTime, foundTime)
	})

	t.Run("failure_repo_returns_not_found", func(t *testing.T) {
		mockRepo := mocks.NewPriceRepositoryInterface(t)

		symbol := "UNKNOWN"
		unixTimestamp := int64(1672531200)
		expectedTime := time.Unix(unixTimestamp, 0)

		expectedError := apperrors.NewNotFound("price not found", sql.ErrNoRows)
		mockRepo.On("GetNearest", ctx, symbol, expectedTime).
			Return(decimal.Zero, time.Time{}, expectedError)

		priceService := NewPriceService(mockRepo, nopLogger)

		_, _, appErr := priceService.GetNearestPrice(ctx, symbol, unixTimestamp)

		require.Error(t, appErr)
		assert.Equal(t, http.StatusNotFound, appErr.Code)
		assert.Equal(t, expectedError, appErr)
	})

	t.Run("failure_repo_returns_internal_error", func(t *testing.T) {
		mockRepo := mocks.NewPriceRepositoryInterface(t)

		symbol := "BTC"
		unixTimestamp := int64(1672531200)
		expectedTime := time.Unix(unixTimestamp, 0)

		expectedError := apperrors.NewInternalServerError("db error", errors.New("connection failed"))
		mockRepo.On("GetNearest", ctx, symbol, expectedTime).
			Return(decimal.Zero, time.Time{}, expectedError)

		priceService := NewPriceService(mockRepo, nopLogger)

		_, _, appErr := priceService.GetNearestPrice(ctx, symbol, unixTimestamp)

		require.Error(t, appErr)
		assert.Equal(t, http.StatusInternalServerError, appErr.Code)
	})
}
