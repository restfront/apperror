// Package apperror предоставляет типизированные ошибки приложения
// с поддержкой кодов ошибок, метаданных и HTTP-статусов.
package apperror

import (
	"maps"
	"net/http"
	"time"
)

// ErrorType представляет категорию ошибки.
type ErrorType uint

const (
	TypeUnknown              ErrorType = iota // 500
	TypeNotValid                              // 422
	TypeBadRequest                            // 400
	TypeUnauthorized                          // 401
	TypeForbidden                             // 403
	TypeNotFound                              // 404
	TypeUnprocessableEntity                   // 422
	TypeInternal                              // 500
	TypeNotImplemented                        // 501
	TypeBadGateway                            // 502
	TypeTemporaryUnavailable                  // 503
	TypeGatewayTimeout                        // 504
	TypeMethodNotAllowed                      // 405
	TypeTooManyRequests                       // 429
)

var errorTypeNames = map[ErrorType]string{
	TypeUnknown:              "Unknown",
	TypeNotValid:             "NotValid",
	TypeBadRequest:           "BadRequest",
	TypeUnauthorized:         "Unauthorized",
	TypeForbidden:            "Forbidden",
	TypeNotFound:             "NotFound",
	TypeUnprocessableEntity:  "UnprocessableEntity",
	TypeInternal:             "Internal",
	TypeNotImplemented:       "NotImplemented",
	TypeBadGateway:           "BadGateway",
	TypeTemporaryUnavailable: "TemporaryUnavailable",
	TypeGatewayTimeout:       "GatewayTimeout",
	TypeMethodNotAllowed:     "MethodNotAllowed",
	TypeTooManyRequests:      "TooManyRequests",
}

// String возвращает строковое представление типа ошибки.
func (eType ErrorType) String() string {
	if name, ok := errorTypeNames[eType]; ok {
		return name
	}
	return "Unknown"
}

var (
	defaultMessages = map[ErrorType]string{
		TypeNotValid:             "Ошибка валидации данных",
		TypeBadRequest:           "Некорректные данные запроса",
		TypeUnauthorized:         "Пользователь не авторизован",
		TypeForbidden:            "Недостаточно прав для выполнения операции",
		TypeNotFound:             "Ресурс не найден",
		TypeUnprocessableEntity:  "Невозможно обработать запрос",
		TypeInternal:             "Внутренняя ошибка сервера",
		TypeNotImplemented:       "Метод не реализован",
		TypeBadGateway:           "Сервис недоступен",
		TypeTemporaryUnavailable: "Сервис временно недоступен",
		TypeGatewayTimeout:       "Таймаут соединения",
		TypeMethodNotAllowed:     "Метод не поддерживается",
		TypeUnknown:              "Неизвестная ошибка",
		TypeTooManyRequests:      "Превышено допустимое количество запросов",
	}
)

// AppError представляет ошибку приложения с типом, кодом и метаданными.
type AppError struct {
	errorType ErrorType
	code      ErrorCode
	message   string
	original  error
	details   map[string]any
	timestamp time.Time
	requestID string
	fields    []FieldError
	stack     []uintptr
}

func (e *AppError) Error() string {
	if e.original == nil {
		return e.message
	}
	return e.original.Error()
}

func (e *AppError) Unwrap() error {
	return e.original
}

func (e *AppError) Type() ErrorType {
	return e.errorType
}

func (e *AppError) Message() string {
	if e.message != "" {
		return e.message
	}
	return defaultMessages[e.errorType]
}

func (e *AppError) HTTPStatusCode() int {
	status := http.StatusInternalServerError
	switch e.errorType {
	case TypeNotValid:
		status = http.StatusUnprocessableEntity
	case TypeBadRequest:
		status = http.StatusBadRequest
	case TypeUnauthorized:
		status = http.StatusUnauthorized
	case TypeForbidden:
		status = http.StatusForbidden
	case TypeNotFound:
		status = http.StatusNotFound
	case TypeUnprocessableEntity:
		status = http.StatusUnprocessableEntity
	case TypeInternal:
		status = http.StatusInternalServerError
	case TypeNotImplemented:
		status = http.StatusNotImplemented
	case TypeBadGateway:
		status = http.StatusBadGateway
	case TypeTemporaryUnavailable:
		status = http.StatusServiceUnavailable
	case TypeGatewayTimeout:
		status = http.StatusGatewayTimeout
	case TypeMethodNotAllowed:
		status = http.StatusMethodNotAllowed
	case TypeTooManyRequests:
		status = http.StatusTooManyRequests
	}

	return status
}

func (e *AppError) WithMessage(message string) *AppError {
	e.message = message
	return e
}

// Code возвращает машиночитаемый код ошибки.
func (e *AppError) Code() ErrorCode {
	return e.code
}

// WithCode устанавливает машиночитаемый код ошибки.
func (e *AppError) WithCode(code ErrorCode) *AppError {
	e.code = code
	return e
}

// Details возвращает дополнительные метаданные ошибки.
func (e *AppError) Details() map[string]any {
	return e.details
}

// Detail возвращает значение метаданных по ключу.
func (e *AppError) Detail(key string) (any, bool) {
	if e.details == nil {
		return nil, false
	}
	v, ok := e.details[key]
	return v, ok
}

// WithDetail добавляет одно поле метаданных.
func (e *AppError) WithDetail(key string, value any) *AppError {
	if e.details == nil {
		e.details = make(map[string]any)
	}
	e.details[key] = value
	return e
}

// WithDetails добавляет несколько полей метаданных.
func (e *AppError) WithDetails(details map[string]any) *AppError {
	if e.details == nil {
		e.details = make(map[string]any)
	}
	maps.Copy(e.details, details)
	return e
}

// Timestamp возвращает время создания ошибки.
func (e *AppError) Timestamp() time.Time {
	return e.timestamp
}

// RequestID возвращает идентификатор запроса.
func (e *AppError) RequestID() string {
	return e.requestID
}

// WithRequestID устанавливает идентификатор запроса.
func (e *AppError) WithRequestID(requestID string) *AppError {
	e.requestID = requestID
	return e
}

func (eType ErrorType) New(message string, original error) *AppError {
	if message == "" {
		message = defaultMessages[eType]
	}
	return &AppError{
		errorType: eType,
		message:   message,
		original:  original,
		timestamp: time.Now(),
	}
}

func NewError(eType ErrorType, message string, original error) *AppError {
	return eType.New(message, original)
}

func NewValidation(message string, original error) *AppError {
	return TypeNotValid.New(message, original)
}

func NewBadRequest(message string, original error) *AppError {
	return TypeBadRequest.New(message, original)
}

func NewUnauthorized(message string, original error) *AppError {
	return TypeUnauthorized.New(message, original)
}

func NewForbidden(message string, original error) *AppError {
	return TypeForbidden.New(message, original)
}

func NewNotFound(message string, original error) *AppError {
	return TypeNotFound.New(message, original)
}

func NewUnprocessableEntity(message string, original error) *AppError {
	return TypeUnprocessableEntity.New(message, original)
}

func NewInternal(message string, original error) *AppError {
	return TypeInternal.New(message, original)
}

func NewTemporaryUnavailable(message string, original error) *AppError {
	return TypeTemporaryUnavailable.New(message, original)
}

func NewGatewayTimeout(message string, original error) *AppError {
	return TypeGatewayTimeout.New(message, original)
}

func NewUnknown(message string, original error) *AppError {
	return TypeUnknown.New(message, original)
}

func NewMethodNotAllowed(message string, original error) *AppError {
	return TypeMethodNotAllowed.New(message, original)
}

func NewTooManyRequests(message string, original error) *AppError {
	return TypeTooManyRequests.New(message, original)
}

func NewNotImplemented(message string, original error) *AppError {
	return TypeNotImplemented.New(message, original)
}

func NewBadGateway(message string, original error) *AppError {
	return TypeBadGateway.New(message, original)
}

func Validation(original error) *AppError {
	return TypeNotValid.New("", original)
}

func BadRequest(original error) *AppError {
	return TypeBadRequest.New("", original)
}

func Unauthorized(original error) *AppError {
	return TypeUnauthorized.New("", original)
}

func Forbidden(original error) *AppError {
	return TypeForbidden.New("", original)
}

func NotFound(original error) *AppError {
	return TypeNotFound.New("", original)
}

func UnprocessableEntity(original error) *AppError {
	return TypeUnprocessableEntity.New("", original)
}

func Internal(original error) *AppError {
	return TypeInternal.New("", original)
}

func NotImplemented(original error) *AppError {
	return TypeNotImplemented.New("", original)
}

func BadGateway(original error) *AppError {
	return TypeBadGateway.New("", original)
}

func TemporaryUnavailable(original error) *AppError {
	return TypeTemporaryUnavailable.New("", original)
}

func GatewayTimeout(original error) *AppError {
	return TypeGatewayTimeout.New("", original)
}

func Unknown(original error) *AppError {
	return TypeUnknown.New("", original)
}

func MethodNotAllowed(original error) *AppError {
	return TypeMethodNotAllowed.New("", original)
}

func TooManyRequests(original error) *AppError {
	return TypeTooManyRequests.New("", original)
}
