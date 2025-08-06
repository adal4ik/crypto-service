package domain

import "github.com/google/uuid"

type Currency struct {
	ID     uuid.UUID
	Symbol string
}

type Price struct {
	Price     float64
	Timestamp int64
}
