package response

import (
	"encoding/json"
	"net/http"
)

// SuccessResponse - универсальный ответ для успешных операций.
// Поле Data теперь имеет тип interface{}, что позволяет передавать любую структуру.
type SuccessResponse struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

// New - конструктор для SuccessResponse.
func New(code int, status string, data interface{}) *SuccessResponse {
	return &SuccessResponse{
		Code:   code,
		Status: status,
		Data:   data,
	}
}

// Send отправляет SuccessResponse клиенту.
func (r *SuccessResponse) Send(w http.ResponseWriter) {
	j, err := json.MarshalIndent(r, "", "\t")
	if err != nil {
		APIError{
			Code:    http.StatusInternalServerError,
			Message: "failed to marshal success response",
		}.Send(w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Code)
	w.Write(j)
}

// --- Код для APIError остается без изменений ---

type APIError struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Resource string `json:"resource"`
}

func (e APIError) Send(w http.ResponseWriter) {
	j, err := json.MarshalIndent(e, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Code)
	w.Write(j)
}
