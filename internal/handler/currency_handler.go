package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/adal4ik/crypto-service/internal/domain/dto"
	"github.com/adal4ik/crypto-service/internal/service"
	"github.com/adal4ik/crypto-service/pkg/apperrors"
	"github.com/adal4ik/crypto-service/pkg/logger"
	"github.com/adal4ik/crypto-service/pkg/response"
	"go.uber.org/zap"
)

type CurrencyHandler struct {
	service service.CurrencyServiceInterface
	logger  logger.Logger
}

func NewCurrencyHandler(service service.CurrencyServiceInterface, logger logger.Logger) *CurrencyHandler {
	return &CurrencyHandler{
		service: service,
		logger:  logger,
	}
}

func (h *CurrencyHandler) handleError(w http.ResponseWriter, r *http.Request, err error) {
	var appErr *apperrors.AppError

	isAppError := errors.As(err, &appErr)

	if isAppError && appErr.Code >= 400 && appErr.Code < 500 {
		h.logger.Warn("Client Error",
			zap.Int("status_code", appErr.Code),
			zap.String("message", appErr.Message),
			zap.Error(err),
			zap.String("url", r.URL.Path),
		)
	} else {
		h.logger.Error("Server Error",
			zap.Error(err),
			zap.String("url", r.URL.Path),
		)
	}

	if isAppError {
		jsonErr := response.APIError{
			Code:     appErr.Code,
			Message:  appErr.Message,
			Resource: r.URL.Path,
		}
		jsonErr.Send(w)
		return
	}

	jsonErr := response.APIError{
		Code:     http.StatusInternalServerError,
		Message:  "Internal Server Error",
		Resource: r.URL.Path,
	}
	jsonErr.Send(w)
}
func (h *CurrencyHandler) CreateCurrency(w http.ResponseWriter, r *http.Request) {
	var req dto.AddCurrencyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, r, apperrors.NewBadRequest("invalid request body", err))
		return
	}

	if err := h.service.AddCurrency(r.Context(), req.Symbol); err != nil {
		h.handleError(w, r, err)
		return
	}

	response.New(http.StatusCreated, "success", "Currency added to tracking list").Send(w)
}

func (h *CurrencyHandler) RemoveCurrency(w http.ResponseWriter, r *http.Request) {
	var req dto.RemoveCurrencyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, r, apperrors.NewBadRequest("invalid request body", err))
		return
	}

	if err := h.service.RemoveCurrency(r.Context(), req.Symbol); err != nil {
		h.handleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
