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
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

var httpClient = &http.Client{Timeout: 15 * time.Second}

var symbolToIDMap = map[string]string{
	"BTC":   "bitcoin",
	"ETH":   "ethereum",
	"LTC":   "litecoin",
	"XRP":   "ripple",
	"BCH":   "bitcoin-cash",
	"DOT":   "polkadot",
	"LINK":  "chainlink",
	"ADA":   "cardano",
	"XLM":   "stellar",
	"UNI":   "uniswap",
	"AVAX":  "avalanche-2",
	"SOL":   "solana",
	"MATIC": "matic-network",
	"TRX":   "tron",
	"ALGO":  "algorand",
	"ATOM":  "cosmos",
}

type PriceCollector struct {
	currencyRepo repository.CurrencyRepositoryInterface
	priceRepo    repository.PriceRepositoryInterface
	logger       logger.Logger
	cfg          config.CollectorConfig
}

func NewPriceCollector(
	currencyRepo repository.CurrencyRepositoryInterface,
	priceRepo repository.PriceRepositoryInterface,
	logger logger.Logger,
	cfg config.CollectorConfig,
) *PriceCollector {
	return &PriceCollector{
		currencyRepo: currencyRepo,
		priceRepo:    priceRepo,
		logger:       logger,
		cfg:          cfg,
	}
}

func (pc *PriceCollector) Start(ctx context.Context) {
	l := pc.logger.With(zap.String("service", "PriceCollector"))
	l.Info("Starting price collector...", zap.Duration("interval", pc.cfg.Interval))

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

	var coingeckoIDs []string
	for _, s := range symbols {
		if id, ok := symbolToIDMap[strings.ToUpper(s)]; ok {
			coingeckoIDs = append(coingeckoIDs, id)
		} else {
			l.Warn("no coingecko mapping for symbol", zap.String("symbol", s))
		}
	}

	if len(coingeckoIDs) == 0 {
		l.Info("no valid currencies to query from coingecko")
		return
	}

	idsString := strings.Join(coingeckoIDs, ",")
	url := fmt.Sprintf("%s?ids=%s&vs_currencies=usd", pc.cfg.ApiBaseURL, idsString)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		l.Error("failed to create http request", zap.Error(err))
		return
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		l.Error("failed to fetch prices from coingecko", zap.Error(err), zap.String("url", url))
		return
	}
	defer resp.Body.Close()

	var prices map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&prices); err != nil {
		l.Error("failed to decode coingecko response", zap.Error(err))
		return
	}
	l.Info("successfully fetched prices", zap.Any("prices", prices))

	now := time.Now()
	for _, symbol := range symbols {
		coingeckoID, ok := symbolToIDMap[strings.ToUpper(symbol)]
		if !ok {
			continue
		}

		if priceData, ok := prices[coingeckoID]; ok {
			if usdPriceFloat, ok := priceData["usd"]; ok {
				priceDecimal := decimal.NewFromFloat(usdPriceFloat)

				if addErr := pc.priceRepo.Add(ctx, symbol, priceDecimal, now); addErr != nil {
					l.Error("failed to save price to db", zap.Error(addErr), zap.String("symbol", symbol))
				}
			}
		}
	}
}
