// Package apperror
package apperror

import (
	"errors"
)

// Is проверяет, является ли ошибка AppError с указанным типом.
func Is(err error, errorType ErrorType) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type() == errorType
	}
	return false
}

// IsCode проверяет, имеет ли ошибка указанный код.
func IsCode(err error, code ErrorCode) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code() == code
	}
	return false
}

// HasCode проверяет, установлен ли код у ошибки.
func HasCode(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return !appErr.code.IsEmpty()
	}
	return false
}

// GetCode извлекает код ошибки.
// Возвращает код и true, если ошибка является AppError с установленным кодом.
func GetCode(err error) (ErrorCode, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) && !appErr.code.IsEmpty() {
		return appErr.code, true
	}
	return "", false
}

// GetType извлекает тип ошибки.
// Возвращает тип и true, если ошибка является AppError.
func GetType(err error) (ErrorType, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type(), true
	}
	return TypeUnknown, false
}

// GetHTTPStatus извлекает HTTP-статус из ошибки.
// Возвращает статус и true, если ошибка является AppError.
func GetHTTPStatus(err error) (int, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.HTTPStatusCode(), true
	}
	return 0, false
}

// Wrap оборачивает любую ошибку в AppError.
// Если ошибка уже является AppError, возвращает её без изменений.
// Если ошибка nil, возвращает nil.
// Иначе оборачивает в Internal error.
func Wrap(err error) *AppError {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	return Internal(err)
}

// WrapWithType оборачивает ошибку с указанным типом.
func WrapWithType(err error, errorType ErrorType) *AppError {
	if err == nil {
		return nil
	}
	return errorType.New("", err)
}

// WrapWithCode оборачивает ошибку с указанным кодом.
func WrapWithCode(err error, code ErrorCode) *AppError {
	if err == nil {
		return nil
	}
	return Internal(err).WithCode(code)
}

// WrapWithTypeAndCode оборачивает ошибку с указанным типом и кодом.
func WrapWithTypeAndCode(err error, errorType ErrorType, code ErrorCode) *AppError {
	if err == nil {
		return nil
	}
	return errorType.New("", err).WithCode(code)
}

// IsNotFound проверяет, является ли ошибка ошибкой "не найдено".
func IsNotFound(err error) bool {
	return Is(err, TypeNotFound)
}

// IsUnauthorized проверяет, является ли ошибка ошибкой авторизации.
func IsUnauthorized(err error) bool {
	return Is(err, TypeUnauthorized)
}

// IsForbidden проверяет, является ли ошибка ошибкой доступа.
func IsForbidden(err error) bool {
	return Is(err, TypeForbidden)
}

// IsBadRequest проверяет, является ли ошибка ошибкой запроса.
func IsBadRequest(err error) bool {
	return Is(err, TypeBadRequest)
}

// IsValidation проверяет, является ли ошибка ошибкой валидации.
func IsValidation(err error) bool {
	return Is(err, TypeNotValid)
}

// IsInternal проверяет, является ли ошибка внутренней ошибкой.
func IsInternal(err error) bool {
	return Is(err, TypeInternal)
}

// IsTemporaryUnavailable проверяет, является ли ошибка временной недоступностью.
func IsTemporaryUnavailable(err error) bool {
	return Is(err, TypeTemporaryUnavailable)
}

// IsTooManyRequests проверяет, является ли ошибка превышением лимита запросов.
func IsTooManyRequests(err error) bool {
	return Is(err, TypeTooManyRequests)
}

// MultiError представляет коллекцию ошибок.
type MultiError struct {
	errors []*AppError
}

// NewMultiError создаёт новую коллекцию ошибок.
func NewMultiError() *MultiError {
	return &MultiError{}
}

// Add добавляет AppError в коллекцию.
func (m *MultiError) Add(err *AppError) *MultiError {
	if err != nil {
		m.errors = append(m.errors, err)
	}
	return m
}

// AddError добавляет любую ошибку в коллекцию.
// Если ошибка не является AppError, она оборачивается в Internal.
func (m *MultiError) AddError(err error) *MultiError {
	if err != nil {
		m.errors = append(m.errors, Wrap(err))
	}
	return m
}

// HasErrors проверяет, есть ли ошибки в коллекции.
func (m *MultiError) HasErrors() bool {
	return len(m.errors) > 0
}

// Errors возвращает список ошибок.
func (m *MultiError) Errors() []*AppError {
	return m.errors
}

// Len возвращает количество ошибок.
func (m *MultiError) Len() int {
	return len(m.errors)
}

// Error реализует интерфейс error.
func (m *MultiError) Error() string {
	if len(m.errors) == 0 {
		return ""
	}
	if len(m.errors) == 1 {
		return m.errors[0].Error()
	}

	result := ""
	for i, err := range m.errors {
		if i > 0 {
			result += "; "
		}
		result += err.Error()
	}
	return result
}

// ErrorOrNil возвращает nil, если ошибок нет, иначе MultiError.
func (m *MultiError) ErrorOrNil() error {
	if !m.HasErrors() {
		return nil
	}
	return m
}

// First возвращает первую ошибку или nil.
func (m *MultiError) First() *AppError {
	if len(m.errors) == 0 {
		return nil
	}
	return m.errors[0]
}

// Last возвращает последнюю ошибку или nil.
func (m *MultiError) Last() *AppError {
	if len(m.errors) == 0 {
		return nil
	}
	return m.errors[len(m.errors)-1]
}
