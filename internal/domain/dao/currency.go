package dao

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// TrackedCurrencyDAO - это модель, соответствующая таблице tracked_currencies.
// Она содержит все поля таблицы, включая служебные, такие как created_at.
type TrackedCurrencyDAO struct {
	ID        uuid.UUID `db:"id"`
	Symbol    string    `db:"symbol"`
	CreatedAt time.Time `db:"created_at"`
}

// PriceHistoryDAO - это модель, соответствующая таблице price_history.
type PriceHistoryDAO struct {
	CurrencyID uuid.UUID       `db:"currency_id"`
	Price      decimal.Decimal `db:"price"` // В реальных фин. приложениях лучше использовать github.com/shopspring/decimal
	Timestamp  time.Time       `db:"timestamp"`
}
