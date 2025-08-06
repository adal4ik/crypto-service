package domain

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Currency struct {
	ID     uuid.UUID
	Symbol string
}

type Price struct {
	Price     decimal.Decimal
	Timestamp int64
}
