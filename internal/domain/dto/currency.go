package dto

// AddCurrencyRequest - DTO для запроса на добавление валюты.
// POST /currency/add
type AddCurrencyRequest struct {
	Symbol string `json:"symbol"`
}

// RemoveCurrencyRequest - DTO для запроса на удаление валюты.
// POST /currency/remove
type RemoveCurrencyRequest struct {
	Symbol string `json:"symbol"`
}

// GetPriceRequest - DTO для запроса цены.
// GET /currency/price
// Используем теги, чтобы связать поля с параметрами запроса или телом JSON.
type GetPriceRequest struct {
	Coin      string `json:"coin"`
	Timestamp int64  `json:"timestamp"`
}

// PriceResponse - DTO для ответа с ценой.
type PriceResponse struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
}

// GenericResponse - универсальный ответ для простых операций.
type GenericResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}
