// Package apperror
package apperror

import (
	"encoding/json"
	"net/http"
	"time"
)

// ErrorResponse представляет структуру JSON-ответа об ошибке.
type ErrorResponse struct {
	Code      string         `json:"code,omitempty"`
	Message   string         `json:"message"`
	Status    int            `json:"status"`
	RequestID string         `json:"request_id,omitempty"`
	Details   map[string]any `json:"details,omitempty"`
	Fields    []FieldError   `json:"fields,omitempty"`
	Timestamp string         `json:"timestamp,omitempty"`
}

// ToResponse преобразует AppError в ErrorResponse для JSON-сериализации.
func (e *AppError) ToResponse() ErrorResponse {
	resp := ErrorResponse{
		Code:      string(e.code),
		Message:   e.Message(),
		Status:    e.HTTPStatusCode(),
		RequestID: e.requestID,
		Details:   e.details,
	}

	if !e.timestamp.IsZero() {
		resp.Timestamp = e.timestamp.Format(time.RFC3339)
	}

	if e.fields != nil {
		resp.Fields = e.fields
	}

	return resp
}

// ToJSON сериализует ошибку в JSON.
func (e *AppError) ToJSON() ([]byte, error) {
	return json.Marshal(e.ToResponse())
}

// MarshalJSON реализует интерфейс json.Marshaler.
func (e *AppError) MarshalJSON() ([]byte, error) {
	return e.ToJSON()
}

// WriteJSON записывает JSON-ответ об ошибке в http.ResponseWriter.
func (e *AppError) WriteJSON(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(e.HTTPStatusCode())
	return json.NewEncoder(w).Encode(e.ToResponse())
}

// WriteJSONWithRequestID записывает JSON-ответ с указанным request ID.
func (e *AppError) WriteJSONWithRequestID(w http.ResponseWriter, requestID string) error {
	e.requestID = requestID
	return e.WriteJSON(w)
}

// HTTPError - хелпер для быстрого создания HTTP-ответа об ошибке.
// Автоматически оборачивает ошибку в AppError если необходимо.
func HTTPError(w http.ResponseWriter, err error) error {
	appErr := AsAppError(err)
	return appErr.WriteJSON(w)
}

// HTTPErrorWithRequestID - хелпер для HTTP-ответа с request ID.
func HTTPErrorWithRequestID(w http.ResponseWriter, err error, requestID string) error {
	appErr := AsAppError(err)
	return appErr.WriteJSONWithRequestID(w, requestID)
}

// AsAppError преобразует error в *AppError.
// Если ошибка уже является AppError, возвращает её.
// Иначе оборачивает в Internal error.
func AsAppError(err error) *AppError {
	if err == nil {
		return nil
	}

	if appErr, ok := err.(*AppError); ok {
		return appErr
	}

	return Internal(err)
}
