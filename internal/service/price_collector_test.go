package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/adal4ik/crypto-service/internal/config"
	"github.com/adal4ik/crypto-service/internal/repository/mocks"
	"github.com/adal4ik/crypto-service/pkg/logger"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPriceCollector_collectPrices(t *testing.T) {
	nopLogger := logger.NewNopLogger()
	ctx := context.Background()

	t.Run("success_collect_and_save", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.Query().Get("ids"), "bitcoin")
			assert.Contains(t, r.URL.Query().Get("ids"), "ethereum")
			assert.Equal(t, "usd", r.URL.Query().Get("vs_currencies"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]map[string]float64{
				"bitcoin":  {"usd": 65000.50},
				"ethereum": {"usd": 3500.75},
			})
		}))
		defer mockServer.Close()

		mockCurrencyRepo := mocks.NewCurrencyRepositoryInterface(t)
		mockPriceRepo := mocks.NewPriceRepositoryInterface(t)

		trackedSymbols := []string{"BTC", "ETH"}
		mockCurrencyRepo.On("GetAll", ctx).Return(trackedSymbols, nil)

		mockPriceRepo.On("Add", ctx, "BTC", decimal.NewFromFloat(65000.50), mock.AnythingOfType("time.Time")).Return(nil)
		mockPriceRepo.On("Add", ctx, "ETH", decimal.NewFromFloat(3500.75), mock.AnythingOfType("time.Time")).Return(nil)

		cfg := config.CollectorConfig{
			Interval:   1 * time.Minute,
			ApiBaseURL: mockServer.URL,
		}

		collector := NewPriceCollector(mockCurrencyRepo, mockPriceRepo, nopLogger, cfg)

		collector.collectPrices(ctx)

	})

	t.Run("success_no_currencies_to_track", func(t *testing.T) {
		mockCurrencyRepo := mocks.NewCurrencyRepositoryInterface(t)
		mockPriceRepo := mocks.NewPriceRepositoryInterface(t)

		mockCurrencyRepo.On("GetAll", ctx).Return([]string{}, nil)

		cfg := config.CollectorConfig{}
		collector := NewPriceCollector(mockCurrencyRepo, mockPriceRepo, nopLogger, cfg)

		collector.collectPrices(ctx)

		mockPriceRepo.AssertNotCalled(t, "Add", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})
}
