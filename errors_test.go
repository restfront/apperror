package apperror

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errOriginal = errors.New("original error")

func TestAppErrorMethods(t *testing.T) {
	tests := []struct {
		name        string
		errorType   ErrorType
		message     string
		originalErr error
		newMessage  string
		expectedMsg string
	}{
		{
			name:        "Error() with original error",
			errorType:   TypeBadRequest,
			message:     "Custom message",
			originalErr: errOriginal,
		},
		{
			name:        "Error() without original error",
			errorType:   TypeBadRequest,
			message:     "Custom message",
			originalErr: nil,
		},
		{
			name:        "Message() with custom message",
			errorType:   TypeBadRequest,
			message:     "Custom message",
			originalErr: errOriginal,
			expectedMsg: "Custom message",
		},
		{
			name:        "Message() with default message",
			errorType:   TypeBadRequest,
			message:     "",
			originalErr: errOriginal,
			expectedMsg: "Некорректные данные запроса",
		},
		{
			name:        "WithMessage() updates message",
			errorType:   TypeBadRequest,
			message:     "Initial message",
			originalErr: errOriginal,
			newMessage:  "Updated message",
			expectedMsg: "Updated message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewError(tt.errorType, tt.message, tt.originalErr)

			if tt.originalErr != nil {
				assert.Equal(t, tt.originalErr.Error(), err.Error(), "unexpected error message from Error()")
			} else {
				assert.Equal(t, tt.message, err.Error(), "unexpected error message from Error()")
			}

			if tt.newMessage != "" {
				err.WithMessage(tt.newMessage)
				assert.Equal(t, tt.newMessage, err.Message(), "unexpected message after WithMessage()")
			}

			if tt.expectedMsg != "" {
				assert.Equal(t, tt.expectedMsg, err.Message(), "unexpected error message from Message()")
			}
		})
	}
}

func TestNewErrorConstructors(t *testing.T) {
	tests := []struct {
		name           string
		fn             func(string, error) *AppError
		errorType      ErrorType
		message        string
		expectedMsg    string
		expectedStatus int
	}{
		{
			name:           "NewValidation with custom message",
			fn:             NewValidation,
			errorType:      TypeNotValid,
			message:        "Custom validation error",
			expectedMsg:    "Custom validation error",
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "NewValidation with default message",
			fn:             NewValidation,
			errorType:      TypeNotValid,
			message:        "",
			expectedMsg:    "Ошибка валидации данных",
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "NewBadRequest with custom message",
			fn:             NewBadRequest,
			errorType:      TypeBadRequest,
			message:        "Custom bad request error",
			expectedMsg:    "Custom bad request error",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "NewBadRequest with default message",
			fn:             NewBadRequest,
			errorType:      TypeBadRequest,
			message:        "",
			expectedMsg:    "Некорректные данные запроса",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "NewUnauthorized with custom message",
			fn:             NewUnauthorized,
			errorType:      TypeUnauthorized,
			message:        "Custom unauthorized error",
			expectedMsg:    "Custom unauthorized error",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "NewUnauthorized with default message",
			fn:             NewUnauthorized,
			errorType:      TypeUnauthorized,
			message:        "",
			expectedMsg:    "Пользователь не авторизован",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "NewForbidden with custom message",
			fn:             NewForbidden,
			errorType:      TypeForbidden,
			message:        "Custom forbidden error",
			expectedMsg:    "Custom forbidden error",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "NewForbidden with default message",
			fn:             NewForbidden,
			errorType:      TypeForbidden,
			message:        "",
			expectedMsg:    "Недостаточно прав для выполнения операции",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "NewNotFound with custom message",
			fn:             NewNotFound,
			errorType:      TypeNotFound,
			message:        "Custom not found error",
			expectedMsg:    "Custom not found error",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "NewNotFound with default message",
			fn:             NewNotFound,
			errorType:      TypeNotFound,
			message:        "",
			expectedMsg:    "Ресурс не найден",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "NewUnprocessableEntity with custom message",
			fn:             NewUnprocessableEntity,
			errorType:      TypeUnprocessableEntity,
			message:        "Custom unprocessable entity error",
			expectedMsg:    "Custom unprocessable entity error",
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "NewUnprocessableEntity with default message",
			fn:             NewUnprocessableEntity,
			errorType:      TypeUnprocessableEntity,
			message:        "",
			expectedMsg:    "Невозможно обработать запрос",
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "NewInternal with custom message",
			fn:             NewInternal,
			errorType:      TypeInternal,
			message:        "Custom unprocessable entity error",
			expectedMsg:    "Custom unprocessable entity error",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "NewInternal with default message",
			fn:             NewInternal,
			errorType:      TypeInternal,
			message:        "",
			expectedMsg:    "Внутренняя ошибка сервера",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "NewTemporaryUnavailable with custom message",
			fn:             NewTemporaryUnavailable,
			errorType:      TypeTemporaryUnavailable,
			message:        "Custom temporary unavailable error",
			expectedMsg:    "Custom temporary unavailable error",
			expectedStatus: http.StatusServiceUnavailable,
		},
		{
			name:           "NewTemporaryUnavailable with default message",
			fn:             NewTemporaryUnavailable,
			errorType:      TypeTemporaryUnavailable,
			message:        "",
			expectedMsg:    "Сервис временно недоступен",
			expectedStatus: http.StatusServiceUnavailable,
		},
		{
			name:           "NewMethodNotAllowed with custom message",
			fn:             NewMethodNotAllowed,
			errorType:      TypeMethodNotAllowed,
			message:        "Custom unprocessable entity error",
			expectedMsg:    "Custom unprocessable entity error",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "NewMethodNotAllowed with default message",
			fn:             NewMethodNotAllowed,
			errorType:      TypeMethodNotAllowed,
			message:        "",
			expectedMsg:    "Метод не поддерживается",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "NewUnknown with custom message",
			fn:             NewUnknown,
			errorType:      TypeUnknown,
			message:        "Custom unknown error",
			expectedMsg:    "Custom unknown error",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "NewUnknown with default message",
			fn:             NewUnknown,
			errorType:      TypeUnknown,
			message:        "",
			expectedMsg:    "Неизвестная ошибка",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "NewTooManyRequests with custom message",
			fn:             NewTooManyRequests,
			errorType:      TypeTooManyRequests,
			message:        "Custom too many requests error",
			expectedMsg:    "Custom too many requests error",
			expectedStatus: http.StatusTooManyRequests,
		},
		{
			name:           "NewTooManyRequests with default message",
			fn:             NewTooManyRequests,
			errorType:      TypeTooManyRequests,
			message:        "",
			expectedMsg:    "Превышено допустимое количество запросов",
			expectedStatus: http.StatusTooManyRequests,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn(tt.message, errOriginal)

			assert.Equal(t, tt.errorType, err.Type(), "unexpected error type")
			assert.Equal(t, tt.expectedMsg, err.Message(), "unexpected error message")
			assert.Equal(t, tt.expectedStatus, err.HTTPStatusCode(), "unexpected HTTP status code")
			assert.Equal(t, errOriginal, err.Unwrap(), "unexpected original error")
		})
	}
}

func TestErrorShortcutFunctions(t *testing.T) {
	tests := []struct {
		name           string
		fn             func(error) *AppError
		expectedType   ErrorType
		expectedMsg    string
		expectedStatus int
	}{
		{
			name:           "Validation shortcut function",
			fn:             Validation,
			expectedType:   TypeNotValid,
			expectedMsg:    "Ошибка валидации данных",
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "BadRequest shortcut function",
			fn:             BadRequest,
			expectedType:   TypeBadRequest,
			expectedMsg:    "Некорректные данные запроса",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Unauthorized shortcut function",
			fn:             Unauthorized,
			expectedType:   TypeUnauthorized,
			expectedMsg:    "Пользователь не авторизован",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Forbidden shortcut function",
			fn:             Forbidden,
			expectedType:   TypeForbidden,
			expectedMsg:    "Недостаточно прав для выполнения операции",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "NotFound shortcut function",
			fn:             NotFound,
			expectedType:   TypeNotFound,
			expectedMsg:    "Ресурс не найден",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Internal shortcut function",
			fn:             Internal,
			expectedType:   TypeInternal,
			expectedMsg:    "Внутренняя ошибка сервера",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "TemporaryUnavailable shortcut function",
			fn:             TemporaryUnavailable,
			expectedType:   TypeTemporaryUnavailable,
			expectedMsg:    "Сервис временно недоступен",
			expectedStatus: http.StatusServiceUnavailable,
		},
		{
			name:           "TooManyRequests shortcut function",
			fn:             TooManyRequests,
			expectedType:   TypeTooManyRequests,
			expectedMsg:    "Превышено допустимое количество запросов",
			expectedStatus: http.StatusTooManyRequests,
		},
		{
			name:           "UnprocessableEntity shortcut function",
			fn:             UnprocessableEntity,
			expectedType:   TypeUnprocessableEntity,
			expectedMsg:    "Невозможно обработать запрос",
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "MethodNotAllowed shortcut function",
			fn:             MethodNotAllowed,
			expectedType:   TypeMethodNotAllowed,
			expectedMsg:    "Метод не поддерживается",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Unknown shortcut function",
			fn:             Unknown,
			expectedType:   TypeUnknown,
			expectedMsg:    "Неизвестная ошибка",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn(errOriginal)

			assert.Equal(t, tt.expectedType, err.Type(), "unexpected error type")
			assert.Equal(t, tt.expectedMsg, err.Message(), "unexpected error message")
			assert.Equal(t, tt.expectedStatus, err.HTTPStatusCode(), "unexpected HTTP status code")
			assert.Equal(t, errOriginal, err.Unwrap(), "unexpected original error")
		})
	}
}
