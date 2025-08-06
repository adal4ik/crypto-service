package handler

import (
	"net/http"
	"strconv"

	"github.com/adal4ik/crypto-service/internal/domain/dto"
	"github.com/adal4ik/crypto-service/internal/service"
	"github.com/adal4ik/crypto-service/pkg/apperrors"
	"github.com/adal4ik/crypto-service/pkg/logger"
	"github.com/adal4ik/crypto-service/pkg/response"
)

type PriceHandler struct {
	service     service.PriceServiceInterface
	logger      logger.Logger
	handleError func(w http.ResponseWriter, r *http.Request, err error)
}

func NewPriceHandler(
	s service.PriceServiceInterface,
	l logger.Logger,
	errorHandler func(w http.ResponseWriter, r *http.Request, err error),
) *PriceHandler {
	return &PriceHandler{
		service:     s,
		logger:      l,
		handleError: errorHandler,
	}
}

func (h *PriceHandler) GetPrice(w http.ResponseWriter, r *http.Request) {
	coin := r.URL.Query().Get("coin")
	timestampStr := r.URL.Query().Get("timestamp")

	if coin == "" {
		h.handleError(w, r, apperrors.NewBadRequest("query parameter 'coin' is required", nil))
		return
	}
	if timestampStr == "" {
		h.handleError(w, r, apperrors.NewBadRequest("query parameter 'timestamp' is required", nil))
		return
	}
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		h.handleError(w, r, apperrors.NewBadRequest("query parameter 'timestamp' must be a valid integer", err))
		return
	}

	price, foundTime, appErr := h.service.GetNearestPrice(r.Context(), coin, timestamp)
	if appErr != nil {
		h.handleError(w, r, appErr)
		return
	}

	respDTO := dto.PriceResponse{
		Symbol:    coin,
		Price:     price,
		Timestamp: foundTime.Unix(),
	}

	response.New(http.StatusOK, "success", respDTO).Send(w)
}
