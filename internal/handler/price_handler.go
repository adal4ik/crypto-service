package handler

import (
	"encoding/json"
	"net/http"

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
	var req dto.GetPriceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, r, apperrors.NewBadRequest("invalid request body", err))
		return
	}

	if req.Coin == "" {
		h.handleError(w, r, apperrors.NewBadRequest("field 'coin' is required", nil))
		return
	}
	if req.Timestamp == 0 {
		h.handleError(w, r, apperrors.NewBadRequest("field 'timestamp' is required", nil))
		return
	}

	price, foundTime, appErr := h.service.GetNearestPrice(r.Context(), req.Coin, req.Timestamp)
	if appErr != nil {
		h.handleError(w, r, appErr)
		return
	}

	respDTO := dto.PriceResponse{
		Symbol:    req.Coin,
		Price:     price,
		Timestamp: foundTime.Unix(),
	}

	response.New(http.StatusOK, "success", respDTO).Send(w)
}
