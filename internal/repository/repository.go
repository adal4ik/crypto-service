package repository

import (
	"database/sql"

	"github.com/adal4ik/crypto-service/pkg/logger"
)

type Repository struct {
	CurrencyRepository CurrencyRepositoryInterface
	Price              PriceRepositoryInterface
}

func NewRepository(db *sql.DB, logger logger.Logger) *Repository {
	return &Repository{
		CurrencyRepository: NewCurrencyRepository(db, logger),
		Price:              NewPriceRepository(db, logger),
	}
}
