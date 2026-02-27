package apperror

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

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

func TestErrorCode(t *testing.T) {
	t.Run("ErrorCode String method", func(t *testing.T) {
		code := CodeDeviceNotFound
		assert.Equal(t, "device_not_found", code.String())
	})

	t.Run("ErrorCode IsEmpty", func(t *testing.T) {
		var emptyCode ErrorCode
		assert.True(t, emptyCode.IsEmpty())
		assert.False(t, CodeDeviceNotFound.IsEmpty())
	})
}

func TestAppErrorWithCode(t *testing.T) {
	t.Run("WithCode sets error code", func(t *testing.T) {
		err := NotFound(errOriginal).WithCode(CodeDeviceNotFound)

		assert.Equal(t, CodeDeviceNotFound, err.Code())
		assert.Equal(t, TypeNotFound, err.Type())
	})

	t.Run("Code returns empty for error without code", func(t *testing.T) {
		err := NotFound(errOriginal)

		assert.True(t, err.Code().IsEmpty())
	})

	t.Run("Fluent API chain", func(t *testing.T) {
		err := NotFound(errOriginal).
			WithCode(CodeDeviceNotFound).
			WithMessage("Устройство не найдено")

		assert.Equal(t, CodeDeviceNotFound, err.Code())
		assert.Equal(t, "Устройство не найдено", err.Message())
		assert.Equal(t, TypeNotFound, err.Type())
		assert.Equal(t, http.StatusNotFound, err.HTTPStatusCode())
	})
}

func TestAppErrorDetails(t *testing.T) {
	t.Run("WithDetail adds single detail", func(t *testing.T) {
		err := NotFound(errOriginal).
			WithCode(CodeDeviceNotFound).
			WithDetail("device_id", "abc-123")

		details := err.Details()
		assert.Equal(t, "abc-123", details["device_id"])
	})

	t.Run("WithDetails adds multiple details", func(t *testing.T) {
		err := NotFound(errOriginal).
			WithDetails(map[string]any{
				"device_id": "abc-123",
				"user_id":   "user-456",
			})

		details := err.Details()
		assert.Equal(t, "abc-123", details["device_id"])
		assert.Equal(t, "user-456", details["user_id"])
	})

	t.Run("Detail returns value and ok", func(t *testing.T) {
		err := NotFound(errOriginal).WithDetail("key", "value")

		val, ok := err.Detail("key")
		assert.True(t, ok)
		assert.Equal(t, "value", val)

		_, ok = err.Detail("nonexistent")
		assert.False(t, ok)
	})

	t.Run("Details returns nil for error without details", func(t *testing.T) {
		err := NotFound(errOriginal)

		assert.Nil(t, err.Details())
	})

	t.Run("WithDetail and WithDetails can be chained", func(t *testing.T) {
		err := NotFound(errOriginal).
			WithDetail("key1", "value1").
			WithDetails(map[string]any{"key2": "value2"}).
			WithDetail("key3", "value3")

		details := err.Details()
		assert.Equal(t, "value1", details["key1"])
		assert.Equal(t, "value2", details["key2"])
		assert.Equal(t, "value3", details["key3"])
	})
}

func TestAppErrorTimestamp(t *testing.T) {
	t.Run("Timestamp is set on creation", func(t *testing.T) {
		before := time.Now()
		err := NotFound(errOriginal)
		after := time.Now()

		assert.False(t, err.Timestamp().IsZero())
		assert.True(t, err.Timestamp().After(before) || err.Timestamp().Equal(before))
		assert.True(t, err.Timestamp().Before(after) || err.Timestamp().Equal(after))
	})
}

func TestAppErrorRequestID(t *testing.T) {
	t.Run("WithRequestID sets request ID", func(t *testing.T) {
		err := NotFound(errOriginal).WithRequestID("req-123")

		assert.Equal(t, "req-123", err.RequestID())
	})

	t.Run("RequestID returns empty for error without request ID", func(t *testing.T) {
		err := NotFound(errOriginal)

		assert.Empty(t, err.RequestID())
	})
}

func TestCompleteFluentAPI(t *testing.T) {
	t.Run("Complete fluent API usage", func(t *testing.T) {
		err := NotFound(errOriginal).
			WithCode(CodeDeviceNotFound).
			WithMessage("Устройство не найдено").
			WithRequestID("req-abc-123").
			WithDetail("device_id", "dev-456").
			WithDetails(map[string]any{
				"user_id":    "user-789",
				"ip_address": "192.168.1.1",
			})

		assert.Equal(t, TypeNotFound, err.Type())
		assert.Equal(t, CodeDeviceNotFound, err.Code())
		assert.Equal(t, "Устройство не найдено", err.Message())
		assert.Equal(t, http.StatusNotFound, err.HTTPStatusCode())
		assert.Equal(t, "req-abc-123", err.RequestID())
		assert.Equal(t, errOriginal, err.Unwrap())

		details := err.Details()
		assert.Equal(t, "dev-456", details["device_id"])
		assert.Equal(t, "user-789", details["user_id"])
		assert.Equal(t, "192.168.1.1", details["ip_address"])
	})
}

func TestJSONSerialization(t *testing.T) {
	t.Run("ToResponse creates correct response", func(t *testing.T) {
		err := NotFound(errOriginal).
			WithCode(CodeDeviceNotFound).
			WithMessage("Устройство не найдено").
			WithRequestID("req-123").
			WithDetail("device_id", "dev-456")

		resp := err.ToResponse()

		assert.Equal(t, "device_not_found", resp.Code)
		assert.Equal(t, "Устройство не найдено", resp.Message)
		assert.Equal(t, http.StatusNotFound, resp.Status)
		assert.Equal(t, "req-123", resp.RequestID)
		assert.Equal(t, "dev-456", resp.Details["device_id"])
	})

	t.Run("ToJSON returns valid JSON", func(t *testing.T) {
		err := NotFound(errOriginal).
			WithCode(CodeDeviceNotFound).
			WithMessage("Устройство не найдено")

		jsonBytes, jsonErr := err.ToJSON()

		assert.NoError(t, jsonErr)
		assert.Contains(t, string(jsonBytes), "device_not_found")
		assert.Contains(t, string(jsonBytes), "Устройство не найдено")
	})

	t.Run("MarshalJSON works correctly", func(t *testing.T) {
		err := NotFound(errOriginal).WithCode(CodeNotFound)

		jsonBytes, jsonErr := err.MarshalJSON()

		assert.NoError(t, jsonErr)
		assert.NotEmpty(t, jsonBytes)
	})
}

func TestValidationBuilder(t *testing.T) {
	t.Run("ValidationBuilder creates validation error", func(t *testing.T) {
		v := NewValidationBuilder()
		v.AddField("email", "required", "Email обязателен")
		v.AddField("password", "too_short", "Минимум 8 символов")

		assert.True(t, v.HasErrors())

		err := v.Build()
		assert.NotNil(t, err)
		assert.Equal(t, TypeNotValid, err.Type())
		assert.Equal(t, CodeValidationFailed, err.Code())
		assert.Len(t, err.Fields(), 2)
	})

	t.Run("ValidationBuilder returns nil when no errors", func(t *testing.T) {
		v := NewValidationBuilder()

		assert.False(t, v.HasErrors())
		assert.Nil(t, v.Build())
		assert.Nil(t, v.ErrorOrNil())
	})

	t.Run("ValidationBuilder helper methods", func(t *testing.T) {
		v := NewValidationBuilder()
		v.AddRequired("name", "")
		v.AddInvalidEmail("email", "invalid")
		v.AddInvalidFormat("phone", "", "+123")

		assert.True(t, v.HasErrors())
		assert.Len(t, v.Build().Fields(), 3)
	})

	t.Run("WithField adds field error to AppError", func(t *testing.T) {
		err := NewValidation("", nil).
			WithField("email", "required", "Email обязателен").
			WithFieldValue("age", "out_of_range", "Возраст вне диапазона", 150)

		assert.True(t, err.HasFieldErrors())
		assert.Len(t, err.Fields(), 2)
	})
}

func TestHelperFunctions(t *testing.T) {
	t.Run("Is checks error type", func(t *testing.T) {
		err := NotFound(errOriginal)

		assert.True(t, Is(err, TypeNotFound))
		assert.False(t, Is(err, TypeBadRequest))
	})

	t.Run("IsCode checks error code", func(t *testing.T) {
		err := NotFound(errOriginal).WithCode(CodeDeviceNotFound)

		assert.True(t, IsCode(err, CodeDeviceNotFound))
		assert.False(t, IsCode(err, CodeUserNotFound))
	})

	t.Run("HasCode checks if code is set", func(t *testing.T) {
		errWithCode := NotFound(errOriginal).WithCode(CodeDeviceNotFound)
		errWithoutCode := NotFound(errOriginal)

		assert.True(t, HasCode(errWithCode))
		assert.False(t, HasCode(errWithoutCode))
	})

	t.Run("GetCode extracts code", func(t *testing.T) {
		err := NotFound(errOriginal).WithCode(CodeDeviceNotFound)

		code, ok := GetCode(err)
		assert.True(t, ok)
		assert.Equal(t, CodeDeviceNotFound, code)
	})

	t.Run("GetType extracts type", func(t *testing.T) {
		err := NotFound(errOriginal)

		errType, ok := GetType(err)
		assert.True(t, ok)
		assert.Equal(t, TypeNotFound, errType)
	})

	t.Run("Wrap wraps regular error", func(t *testing.T) {
		wrapped := Wrap(errOriginal)

		assert.NotNil(t, wrapped)
		assert.Equal(t, TypeInternal, wrapped.Type())
		assert.Equal(t, errOriginal, wrapped.Unwrap())
	})

	t.Run("Wrap returns nil for nil error", func(t *testing.T) {
		assert.Nil(t, Wrap(nil))
	})

	t.Run("Wrap returns AppError unchanged", func(t *testing.T) {
		original := NotFound(errOriginal)
		wrapped := Wrap(original)

		assert.Same(t, original, wrapped)
	})

	t.Run("IsNotFound helper", func(t *testing.T) {
		assert.True(t, IsNotFound(NotFound(errOriginal)))
		assert.False(t, IsNotFound(BadRequest(errOriginal)))
	})

	t.Run("IsUnauthorized helper", func(t *testing.T) {
		assert.True(t, IsUnauthorized(Unauthorized(errOriginal)))
		assert.False(t, IsUnauthorized(NotFound(errOriginal)))
	})
}

func TestMultiError(t *testing.T) {
	t.Run("MultiError collects errors", func(t *testing.T) {
		m := NewMultiError()
		m.Add(NotFound(errOriginal))
		m.Add(BadRequest(errOriginal))

		assert.True(t, m.HasErrors())
		assert.Equal(t, 2, m.Len())
	})

	t.Run("MultiError.ErrorOrNil returns nil when empty", func(t *testing.T) {
		m := NewMultiError()

		assert.False(t, m.HasErrors())
		assert.Nil(t, m.ErrorOrNil())
	})

	t.Run("MultiError.First and Last", func(t *testing.T) {
		m := NewMultiError()
		first := NotFound(errOriginal)
		last := BadRequest(errOriginal)
		m.Add(first)
		m.Add(last)

		assert.Same(t, first, m.First())
		assert.Same(t, last, m.Last())
	})

	t.Run("MultiError.Error formats correctly", func(t *testing.T) {
		m := NewMultiError()
		m.AddError(errors.New("error 1"))
		m.AddError(errors.New("error 2"))

		errStr := m.Error()
		assert.Contains(t, errStr, "error 1")
		assert.Contains(t, errStr, "error 2")
	})
}

func TestStackTrace(t *testing.T) {
	t.Run("WithStack captures stack", func(t *testing.T) {
		err := NotFound(errOriginal).WithStack()

		assert.True(t, err.HasStack())
		assert.NotEmpty(t, err.Stack())
	})

	t.Run("StackTrace returns formatted string", func(t *testing.T) {
		err := NotFound(errOriginal).WithStack()

		trace := err.StackTrace()
		assert.NotEmpty(t, trace)
		assert.Contains(t, trace, "TestStackTrace")
	})

	t.Run("StackFrames returns structured frames", func(t *testing.T) {
		err := NotFound(errOriginal).WithStack()

		frames := err.StackFrames()
		assert.NotEmpty(t, frames)
		assert.NotEmpty(t, frames[0].Function)
		assert.NotEmpty(t, frames[0].File)
		assert.Greater(t, frames[0].Line, 0)
	})

	t.Run("Verbose returns detailed info", func(t *testing.T) {
		err := NotFound(errOriginal).
			WithCode(CodeDeviceNotFound).
			WithMessage("Устройство не найдено").
			WithRequestID("req-123").
			WithDetail("device_id", "dev-456").
			WithStack()

		verbose := err.Verbose()
		assert.Contains(t, verbose, "NotFound")
		assert.Contains(t, verbose, "device_not_found")
		assert.Contains(t, verbose, "req-123")
		assert.Contains(t, verbose, "device_id")
		assert.Contains(t, verbose, "Stack Trace")
	})
}

func TestErrorTypeString(t *testing.T) {
	tests := []struct {
		errorType ErrorType
		expected  string
	}{
		{TypeUnknown, "Unknown"},
		{TypeNotValid, "NotValid"},
		{TypeBadRequest, "BadRequest"},
		{TypeUnauthorized, "Unauthorized"},
		{TypeForbidden, "Forbidden"},
		{TypeNotFound, "NotFound"},
		{TypeInternal, "Internal"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.errorType.String())
		})
	}
}

func TestLogging(t *testing.T) {
	t.Run("LogFields returns correct fields", func(t *testing.T) {
		err := NotFound(errOriginal).
			WithCode(CodeDeviceNotFound).
			WithMessage("Устройство не найдено").
			WithRequestID("req-123")

		fields := err.LogFields()

		assert.Equal(t, "NotFound", fields["error_type"])
		assert.Equal(t, "device_not_found", fields["error_code"])
		assert.Equal(t, http.StatusNotFound, fields["http_status"])
		assert.Equal(t, "req-123", fields["request_id"])
	})

	t.Run("LogValue returns slog.Value", func(t *testing.T) {
		err := NotFound(errOriginal).WithCode(CodeDeviceNotFound)

		logValue := err.LogValue()
		assert.NotEmpty(t, logValue.String())
	})
}

func TestContext(t *testing.T) {
	t.Run("FromContext enriches error with context data", func(t *testing.T) {
		ctx := context.Background()
		ctx = ContextWithRequestID(ctx, "req-123")
		ctx = ContextWithUserID(ctx, "user-456")

		err := FromContext(ctx, errOriginal)

		assert.Equal(t, "req-123", err.RequestID())
		userID, ok := err.Detail("user_id")
		assert.True(t, ok)
		assert.Equal(t, "user-456", userID)
	})

	t.Run("WithContext adds context data to existing error", func(t *testing.T) {
		ctx := context.Background()
		ctx = ContextWithRequestID(ctx, "req-789")

		err := NotFound(errOriginal).WithContext(ctx)

		assert.Equal(t, "req-789", err.RequestID())
	})

	t.Run("RequestIDFromContext extracts request ID", func(t *testing.T) {
		ctx := ContextWithRequestID(context.Background(), "req-abc")

		assert.Equal(t, "req-abc", RequestIDFromContext(ctx))
	})
}
