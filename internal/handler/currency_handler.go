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
		h.logger.Warn("failed to decode request body", zap.Error(err))
		// Для клиента отправляем более общее сообщение
		http.Error(w, `{"status": "error", "message": "invalid request body"}`, http.StatusBadRequest)
		return
	}

	// 2. Валидируем DTO
	if req.Symbol == "" {
		http.Error(w, `{"status": "error", "message": "symbol is required"}`, http.StatusBadRequest)
		return
	}

	// 3. Вызываем сервис с данными из DTO
	err := h.service.AddCurrency(r.Context(), req.Symbol)
	if err != nil {
		h.logger.Error("failed to add currency", zap.Error(err), zap.String("symbol", req.Symbol))
		http.Error(w, `{"status": "error", "message": "could not add currency"}`, http.StatusInternalServerError)
		return
	}

	// 4. Формируем и отправляем успешный ответ с помощью DTO
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created - правильный статус для успешного создания ресурса
	json.NewEncoder(w).Encode(dto.GenericResponse{
		Status:  "success",
		Message: "Currency added to tracking list",
	})
}
