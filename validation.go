// Package apperror
package apperror

// FieldError представляет ошибку валидации конкретного поля.
type FieldError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Value   any    `json:"value,omitempty"`
}

// Fields возвращает список ошибок валидации полей.
func (e *AppError) Fields() []FieldError {
	return e.fields
}

// WithField добавляет ошибку валидации поля.
func (e *AppError) WithField(field, code, message string) *AppError {
	e.fields = append(e.fields, FieldError{
		Field:   field,
		Code:    code,
		Message: message,
	})
	return e
}

// WithFieldValue добавляет ошибку валидации поля с указанием значения.
func (e *AppError) WithFieldValue(field, code, message string, value any) *AppError {
	e.fields = append(e.fields, FieldError{
		Field:   field,
		Code:    code,
		Message: message,
		Value:   value,
	})
	return e
}

// WithFields добавляет несколько ошибок валидации полей.
func (e *AppError) WithFields(fields []FieldError) *AppError {
	e.fields = append(e.fields, fields...)
	return e
}

// HasFieldErrors проверяет, есть ли ошибки валидации полей.
func (e *AppError) HasFieldErrors() bool {
	return len(e.fields) > 0
}

// ValidationBuilder предоставляет удобный интерфейс для построения ошибок валидации.
type ValidationBuilder struct {
	err *AppError
}

// NewValidationBuilder создаёт новый билдер для ошибок валидации.
func NewValidationBuilder() *ValidationBuilder {
	return &ValidationBuilder{
		err: TypeNotValid.New("", nil).WithCode(CodeValidationFailed),
	}
}

// AddField добавляет ошибку валидации поля.
func (v *ValidationBuilder) AddField(field, code, message string) *ValidationBuilder {
	v.err.WithField(field, code, message)
	return v
}

// AddFieldValue добавляет ошибку валидации поля с указанием значения.
func (v *ValidationBuilder) AddFieldValue(field, code, message string, value any) *ValidationBuilder {
	v.err.WithFieldValue(field, code, message, value)
	return v
}

// AddRequired добавляет ошибку обязательного поля.
func (v *ValidationBuilder) AddRequired(field, message string) *ValidationBuilder {
	if message == "" {
		message = "Поле обязательно для заполнения"
	}
	return v.AddField(field, string(CodeRequiredField), message)
}

// AddInvalidFormat добавляет ошибку неверного формата.
func (v *ValidationBuilder) AddInvalidFormat(field, message string, value any) *ValidationBuilder {
	if message == "" {
		message = "Неверный формат"
	}
	return v.AddFieldValue(field, string(CodeInvalidFormat), message, value)
}

// AddTooShort добавляет ошибку слишком короткого значения.
func (v *ValidationBuilder) AddTooShort(field, message string, minLength int) *ValidationBuilder {
	if message == "" {
		message = "Значение слишком короткое"
	}
	v.err.WithField(field, string(CodeTooShort), message)
	if v.err.details == nil {
		v.err.details = make(map[string]any)
	}
	return v
}

// AddTooLong добавляет ошибку слишком длинного значения.
func (v *ValidationBuilder) AddTooLong(field, message string, maxLength int) *ValidationBuilder {
	if message == "" {
		message = "Значение слишком длинное"
	}
	v.err.WithField(field, string(CodeTooLong), message)
	return v
}

// AddOutOfRange добавляет ошибку выхода за пределы диапазона.
func (v *ValidationBuilder) AddOutOfRange(field, message string, value any) *ValidationBuilder {
	if message == "" {
		message = "Значение вне допустимого диапазона"
	}
	return v.AddFieldValue(field, string(CodeOutOfRange), message, value)
}

// AddInvalidEmail добавляет ошибку неверного email.
func (v *ValidationBuilder) AddInvalidEmail(field string, value any) *ValidationBuilder {
	return v.AddFieldValue(field, string(CodeInvalidEmail), "Неверный формат email", value)
}

// AddInvalidPhone добавляет ошибку неверного телефона.
func (v *ValidationBuilder) AddInvalidPhone(field string, value any) *ValidationBuilder {
	return v.AddFieldValue(field, string(CodeInvalidPhone), "Неверный формат телефона", value)
}

// WithMessage устанавливает общее сообщение об ошибке валидации.
func (v *ValidationBuilder) WithMessage(message string) *ValidationBuilder {
	v.err.WithMessage(message)
	return v
}

// WithRequestID устанавливает идентификатор запроса.
func (v *ValidationBuilder) WithRequestID(requestID string) *ValidationBuilder {
	v.err.WithRequestID(requestID)
	return v
}

// HasErrors проверяет, есть ли ошибки валидации.
func (v *ValidationBuilder) HasErrors() bool {
	return v.err.HasFieldErrors()
}

// Build возвращает построенную ошибку или nil, если ошибок нет.
func (v *ValidationBuilder) Build() *AppError {
	if !v.HasErrors() {
		return nil
	}
	return v.err
}

// Error возвращает ошибку независимо от наличия полей.
// Используйте Build() если хотите получить nil при отсутствии ошибок.
func (v *ValidationBuilder) Error() *AppError {
	return v.err
}

// ErrorOrNil возвращает ошибку или nil, если ошибок нет.
// Алиас для Build().
func (v *ValidationBuilder) ErrorOrNil() error {
	if !v.HasErrors() {
		return nil
	}
	return v.err
}
