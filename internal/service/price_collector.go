package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/adal4ik/crypto-service/internal/config"
	"github.com/adal4ik/crypto-service/internal/repository"
	"github.com/adal4ik/crypto-service/pkg/logger"
	"go.uber.org/zap"
)

type PriceCollector struct {
	currencyRepo repository.CurrencyRepositoryInterface
	priceRepo    repository.PriceRepositoryInterface
	logger       logger.Logger
	cfg          config.CollectorConfig // <-- ДОБАВЬТЕ ПОЛЕ
}

func NewPriceCollector(
	currencyRepo repository.CurrencyRepositoryInterface,
	priceRepo repository.PriceRepositoryInterface,
	logger logger.Logger,
	cfg config.CollectorConfig, // <-- ДОБАВЬТЕ АРГУМЕНТ
) *PriceCollector {
	return &PriceCollector{
		currencyRepo: currencyRepo,
		priceRepo:    priceRepo,
		logger:       logger,
		cfg:          cfg, // <-- СОХРАНИТЕ ЕГО
	}
}

// Start запускает бесконечный цикл сбора цен.
func (pc *PriceCollector) Start(ctx context.Context) {
	l := pc.logger.With(zap.String("service", "PriceCollector"))
	l.Info("Starting price collector...", zap.Duration("interval", pc.cfg.Interval))

	// Запускаем сборщик раз в минуту. Для теста можно поставить 10 секунд.
	ticker := time.NewTicker(pc.cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			l.Info("Collector tick: starting price collection job")
			pc.collectPrices(ctx)
		case <-ctx.Done():
			l.Info("Stopping price collector...")
			return
		}
	}
}

func (pc *PriceCollector) collectPrices(ctx context.Context) {
	l := pc.logger.With(zap.String("job", "collectPrices"))

	// 1. Получаем список валют из нашей БД.
	symbols, appErr := pc.currencyRepo.GetAll(ctx)
	if appErr != nil {
		l.Error("failed to get tracked currencies", zap.Error(appErr))
		return
	}
	if len(symbols) == 0 {
		l.Info("no currencies to track, skipping collection")
		return
	}
	l.Info("found currencies to track", zap.Strings("symbols", symbols))

	// 2. Идем во внешнее API (CoinGecko).
	// CoinGecko использует полные имена (bitcoin, ethereum), а мы храним символы (BTC, ETH).
	// Для простоты будем считать, что они совпадают (в реальном проекте нужен был бы маппинг).
	// В запросе они должны быть в нижнем регистре.
	var lowerCaseSymbols []string
	for _, s := range symbols {
		lowerCaseSymbols = append(lowerCaseSymbols, strings.ToLower(s))
	}

	ids := strings.Join(lowerCaseSymbols, ",")
	url := fmt.Sprintf("%s?ids=%s&vs_currencies=usd", pc.cfg.ApiBaseURL, ids)

	resp, err := http.Get(url)
	if err != nil {
		l.Error("failed to fetch prices from coingecko", zap.Error(err))
		return
	}
	defer resp.Body.Close()

	// 3. Парсим ответ.
	var prices map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&prices); err != nil {
		l.Error("failed to decode coingecko response", zap.Error(err))
		return
	}
	l.Info("successfully fetched prices", zap.Any("prices", prices))

	// 4. Сохраняем каждую цену в нашу БД.
	now := time.Now()
	for i, s := range symbols {
		// CoinGecko возвращает ключ в нижнем регистре
		id := lowerCaseSymbols[i]
		if priceData, ok := prices[id]; ok {
			if usdPrice, ok := priceData["usd"]; ok {
				pc.priceRepo.Add(ctx, s, usdPrice, now)
			}
		}
	}
}
