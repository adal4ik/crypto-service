package service

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/adal4ik/crypto-service/internal/repository/mocks"
	"github.com/adal4ik/crypto-service/pkg/apperrors"
	"github.com/adal4ik/crypto-service/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCurrencyService_AddCurrency(t *testing.T) {
	nopLogger := logger.NewNopLogger()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mockRepo := mocks.NewCurrencyRepositoryInterface(t)

		mockRepo.On("Add", ctx, "BTC").Return(nil)

		currencyService := NewCurrencyService(mockRepo, nopLogger)

		appErr := currencyService.AddCurrency(ctx, "  btc  ")

		// Assert
		assert.Nil(t, appErr)
	})

	t.Run("failure_empty_symbol", func(t *testing.T) {
		mockRepo := mocks.NewCurrencyRepositoryInterface(t)
		currencyService := NewCurrencyService(mockRepo, nopLogger)

		appErr := currencyService.AddCurrency(ctx, "   ")

		require.Error(t, appErr)
		assert.Equal(t, http.StatusBadRequest, appErr.Code)
		assert.Equal(t, "currency symbol cannot be empty", appErr.Message)
		mockRepo.AssertNotCalled(t, "Add", ctx, "")
	})

	t.Run("failure_repo_error", func(t *testing.T) {
		mockRepo := mocks.NewCurrencyRepositoryInterface(t)
		expectedError := apperrors.NewInternalServerError("db error", errors.New("something went wrong"))

		mockRepo.On("Add", ctx, "ETH").Return(expectedError)

		currencyService := NewCurrencyService(mockRepo, nopLogger)

		appErr := currencyService.AddCurrency(ctx, "ETH")

		require.Error(t, appErr)
		assert.Equal(t, expectedError, appErr)
	})
}

func TestCurrencyService_RemoveCurrency(t *testing.T) {
	nopLogger := logger.NewNopLogger()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := mocks.NewCurrencyRepositoryInterface(t)

		mockRepo.On("Remove", ctx, "XRP").Return(nil)

		currencyService := NewCurrencyService(mockRepo, nopLogger)

		appErr := currencyService.RemoveCurrency(ctx, " xrp ")

		assert.Nil(t, appErr)
	})

	t.Run("failure_empty_symbol", func(t *testing.T) {
		mockRepo := mocks.NewCurrencyRepositoryInterface(t)
		currencyService := NewCurrencyService(mockRepo, nopLogger)

		appErr := currencyService.RemoveCurrency(ctx, "")

		require.Error(t, appErr)
		assert.Equal(t, http.StatusBadRequest, appErr.Code)
		mockRepo.AssertNotCalled(t, "Remove", ctx, "")
	})
}
