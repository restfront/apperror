// Package apperror
package apperror

import (
	"context"
)

// ContextKey тип для ключей контекста.
type ContextKey string

// Предопределённые ключи контекста.
const (
	ContextKeyRequestID ContextKey = "request_id"
	ContextKeyUserID    ContextKey = "user_id"
	ContextKeyTraceID   ContextKey = "trace_id"
	ContextKeySpanID    ContextKey = "span_id"
)

// FromContext создаёт AppError из обычной ошибки, обогащая её данными из контекста.
func FromContext(ctx context.Context, err error) *AppError {
	appErr := Wrap(err)
	if appErr == nil {
		return nil
	}

	// Извлекаем request_id
	if requestID, ok := ctx.Value(ContextKeyRequestID).(string); ok && requestID != "" {
		appErr.WithRequestID(requestID)
	}

	// Извлекаем user_id
	if userID, ok := ctx.Value(ContextKeyUserID).(string); ok && userID != "" {
		appErr.WithDetail("user_id", userID)
	}

	// Извлекаем trace_id
	if traceID, ok := ctx.Value(ContextKeyTraceID).(string); ok && traceID != "" {
		appErr.WithDetail("trace_id", traceID)
	}

	// Извлекаем span_id
	if spanID, ok := ctx.Value(ContextKeySpanID).(string); ok && spanID != "" {
		appErr.WithDetail("span_id", spanID)
	}

	return appErr
}

// FromContextWithType создаёт AppError указанного типа из контекста.
func FromContextWithType(ctx context.Context, err error, errorType ErrorType) *AppError {
	appErr := WrapWithType(err, errorType)
	if appErr == nil {
		return nil
	}

	return enrichFromContext(ctx, appErr)
}

// FromContextWithCode создаёт AppError с указанным кодом из контекста.
func FromContextWithCode(ctx context.Context, err error, code ErrorCode) *AppError {
	appErr := WrapWithCode(err, code)
	if appErr == nil {
		return nil
	}

	return enrichFromContext(ctx, appErr)
}

// enrichFromContext обогащает ошибку данными из контекста.
func enrichFromContext(ctx context.Context, appErr *AppError) *AppError {
	if requestID, ok := ctx.Value(ContextKeyRequestID).(string); ok && requestID != "" {
		appErr.WithRequestID(requestID)
	}

	if userID, ok := ctx.Value(ContextKeyUserID).(string); ok && userID != "" {
		appErr.WithDetail("user_id", userID)
	}

	if traceID, ok := ctx.Value(ContextKeyTraceID).(string); ok && traceID != "" {
		appErr.WithDetail("trace_id", traceID)
	}

	if spanID, ok := ctx.Value(ContextKeySpanID).(string); ok && spanID != "" {
		appErr.WithDetail("span_id", spanID)
	}

	return appErr
}

// WithContext добавляет данные из контекста в существующую ошибку.
func (e *AppError) WithContext(ctx context.Context) *AppError {
	return enrichFromContext(ctx, e)
}

// ContextWithRequestID добавляет request_id в контекст.
func ContextWithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, ContextKeyRequestID, requestID)
}

// ContextWithUserID добавляет user_id в контекст.
func ContextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, ContextKeyUserID, userID)
}

// ContextWithTraceID добавляет trace_id в контекст.
func ContextWithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, ContextKeyTraceID, traceID)
}

// ContextWithSpanID добавляет span_id в контекст.
func ContextWithSpanID(ctx context.Context, spanID string) context.Context {
	return context.WithValue(ctx, ContextKeySpanID, spanID)
}

// RequestIDFromContext извлекает request_id из контекста.
func RequestIDFromContext(ctx context.Context) string {
	if requestID, ok := ctx.Value(ContextKeyRequestID).(string); ok {
		return requestID
	}
	return ""
}

// UserIDFromContext извлекает user_id из контекста.
func UserIDFromContext(ctx context.Context) string {
	if userID, ok := ctx.Value(ContextKeyUserID).(string); ok {
		return userID
	}
	return ""
}
