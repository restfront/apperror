// Package apperror
package apperror

// ErrorCode представляет машиночитаемый код ошибки.
// Используется для идентификации конкретного типа ошибки клиентом.
type ErrorCode string

// String возвращает строковое представление кода ошибки.
func (c ErrorCode) String() string {
	return string(c)
}

// IsEmpty проверяет, пустой ли код ошибки.
func (c ErrorCode) IsEmpty() bool {
	return c == ""
}

// Предопределённые коды ошибок.
const (
	// Общие ошибки
	CodeUnknown       ErrorCode = "unknown_error"
	CodeInternalError ErrorCode = "internal_error"

	// Аутентификация и авторизация
	CodeUnauthorized       ErrorCode = "unauthorized"
	CodeInvalidToken       ErrorCode = "invalid_token"
	CodeTokenExpired       ErrorCode = "token_expired"
	CodeInvalidCredentials ErrorCode = "invalid_credentials"
	CodeAccessDenied       ErrorCode = "access_denied"
	CodePermissionDenied   ErrorCode = "permission_denied"
	CodeSessionExpired     ErrorCode = "session_expired"

	// Ресурсы
	CodeNotFound         ErrorCode = "not_found"
	CodeResourceNotFound ErrorCode = "resource_not_found"
	CodeUserNotFound     ErrorCode = "user_not_found"
	CodeDeviceNotFound   ErrorCode = "device_not_found"
	CodeFileNotFound     ErrorCode = "file_not_found"
	CodeAlreadyExists    ErrorCode = "already_exists"
	CodeDuplicateEntry   ErrorCode = "duplicate_entry"
	CodeConflict         ErrorCode = "conflict"

	// Валидация
	CodeValidationFailed ErrorCode = "validation_failed"
	CodeInvalidInput     ErrorCode = "invalid_input"
	CodeInvalidFormat    ErrorCode = "invalid_format"
	CodeRequiredField    ErrorCode = "required_field"
	CodeInvalidEmail     ErrorCode = "invalid_email"
	CodeInvalidPhone     ErrorCode = "invalid_phone"
	CodeInvalidURL       ErrorCode = "invalid_url"
	CodeTooShort         ErrorCode = "too_short"
	CodeTooLong          ErrorCode = "too_long"
	CodeOutOfRange       ErrorCode = "out_of_range"
	CodeInvalidValue     ErrorCode = "invalid_value"

	// Лимиты и ограничения
	CodeRateLimited     ErrorCode = "rate_limited"
	CodeQuotaExceeded   ErrorCode = "quota_exceeded"
	CodePayloadTooLarge ErrorCode = "payload_too_large"
	CodeTooManyRequests ErrorCode = "too_many_requests"
	CodeLimitExceeded   ErrorCode = "limit_exceeded"

	// Внешние сервисы и инфраструктура
	CodeServiceUnavailable ErrorCode = "service_unavailable"
	CodeTimeout            ErrorCode = "timeout"
	CodeConnectionFailed   ErrorCode = "connection_failed"
	CodeDependencyFailed   ErrorCode = "dependency_failed"
	CodeDatabaseError      ErrorCode = "database_error"
	CodeCacheError         ErrorCode = "cache_error"
	CodeExternalAPIError   ErrorCode = "external_api_error"

	// Операции
	CodeOperationFailed   ErrorCode = "operation_failed"
	CodeOperationCanceled ErrorCode = "operation_canceled"
	CodeOperationTimeout  ErrorCode = "operation_timeout"
	CodeNotImplemented    ErrorCode = "not_implemented"
	CodeMethodNotAllowed  ErrorCode = "method_not_allowed"

	// Данные
	CodeDataCorrupted        ErrorCode = "data_corrupted"
	CodeDataInconsistent     ErrorCode = "data_inconsistent"
	CodeSerializationError   ErrorCode = "serialization_error"
	CodeDeserializationError ErrorCode = "deserialization_error"
)
