// Package apperror
package apperror

import (
	"log/slog"
)

// LogFields возвращает поля для структурированного логирования.
func (e *AppError) LogFields() map[string]any {
	fields := map[string]any{
		"error_type":  e.errorType.String(),
		"http_status": e.HTTPStatusCode(),
	}

	if !e.code.IsEmpty() {
		fields["error_code"] = string(e.code)
	}

	if e.message != "" {
		fields["message"] = e.message
	}

	if e.requestID != "" {
		fields["request_id"] = e.requestID
	}

	if e.original != nil {
		fields["original_error"] = e.original.Error()
	}

	if !e.timestamp.IsZero() {
		fields["timestamp"] = e.timestamp
	}

	if len(e.details) > 0 {
		fields["details"] = e.details
	}

	if len(e.fields) > 0 {
		fields["validation_fields"] = e.fields
	}

	return fields
}

// LogValue реализует интерфейс slog.LogValuer для интеграции с slog.
func (e *AppError) LogValue() slog.Value {
	attrs := []slog.Attr{
		slog.String("type", e.errorType.String()),
		slog.Int("status", e.HTTPStatusCode()),
	}

	if !e.code.IsEmpty() {
		attrs = append(attrs, slog.String("code", string(e.code)))
	}

	if e.message != "" {
		attrs = append(attrs, slog.String("message", e.message))
	}

	if e.requestID != "" {
		attrs = append(attrs, slog.String("request_id", e.requestID))
	}

	if e.original != nil {
		attrs = append(attrs, slog.String("original", e.original.Error()))
	}

	return slog.GroupValue(attrs...)
}

// SlogAttrs возвращает атрибуты для slog в виде среза.
func (e *AppError) SlogAttrs() []slog.Attr {
	attrs := []slog.Attr{
		slog.String("error_type", e.errorType.String()),
		slog.Int("http_status", e.HTTPStatusCode()),
	}

	if !e.code.IsEmpty() {
		attrs = append(attrs, slog.String("error_code", string(e.code)))
	}

	if e.message != "" {
		attrs = append(attrs, slog.String("message", e.message))
	}

	if e.requestID != "" {
		attrs = append(attrs, slog.String("request_id", e.requestID))
	}

	if e.original != nil {
		attrs = append(attrs, slog.String("original_error", e.original.Error()))
	}

	if len(e.details) > 0 {
		attrs = append(attrs, slog.Any("details", e.details))
	}

	if len(e.fields) > 0 {
		attrs = append(attrs, slog.Any("validation_fields", e.fields))
	}

	return attrs
}
