// Package apperror
package apperror

import (
	"net/http"
)

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
	TypeTemporaryUnavailable                  // 503
	TypeMethodNotAllowed                      // 405
	TypeTooManyRequests                       // 429
)

var (
	defaultMessages = map[ErrorType]string{
		TypeNotValid:             "Ошибка валидации данных",
		TypeBadRequest:           "Некорректные данные запроса",
		TypeUnauthorized:         "Пользователь не авторизован",
		TypeForbidden:            "Недостаточно прав для выполнения операции",
		TypeNotFound:             "Ресурс не найден",
		TypeUnprocessableEntity:  "Невозможно обработать запрос",
		TypeInternal:             "Внутренняя ошибка сервера",
		TypeTemporaryUnavailable: "Сервис временно недоступен",
		TypeMethodNotAllowed:     "Метод не поддерживается",
		TypeUnknown:              "Неизвестная ошибка",
		TypeTooManyRequests:      "Превышено допустимое количество запросов",
	}
)

type AppError struct {
	errorType ErrorType
	message   string
	original  error
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
	case TypeTemporaryUnavailable:
		status = http.StatusServiceUnavailable
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

func (eType ErrorType) New(message string, original error) *AppError {
	if message == "" {
		message = defaultMessages[eType]
	}
	return &AppError{errorType: eType, message: message, original: original}
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

func NewUnknown(message string, original error) *AppError {
	return TypeUnknown.New(message, original)
}

func NewMethodNotAllowed(message string, original error) *AppError {
	return TypeMethodNotAllowed.New(message, original)
}

func NewTooManyRequests(message string, original error) *AppError {
	return TypeTooManyRequests.New(message, original)
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

func TemporaryUnavailable(original error) *AppError {
	return TypeTemporaryUnavailable.New("", original)
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
