# apperror

Пакет для типизированных ошибок приложения на Go с поддержкой кодов ошибок, метаданных, стек-трейсов и JSON-сериализации.

## Установка

```bash
go get github.com/restfront/apperror
```

## Возможности

- **Типы ошибок** — категоризация ошибок по HTTP-статусам (NotFound, BadRequest, Unauthorized и т.д.)
- **Коды ошибок** — машиночитаемые коды для клиентов (`device_not_found`, `invalid_token`)
- **Метаданные** — произвольные данные об ошибке (details)
- **Валидация** — ошибки валидации с информацией по полям
- **Стек-трейсы** — захват стека вызовов для отладки
- **JSON-сериализация** — готовые структуры для HTTP API
- **Логирование** — интеграция с `slog`
- **Контекст** — автоматическое обогащение из `context.Context`

## Быстрый старт

### Базовое использование

```go
import "github.com/restfront/apperror"

// Простая ошибка
err := apperror.NotFound(originalErr)

// С кодом и сообщением
err := apperror.NotFound(originalErr).
    WithCode(apperror.CodeDeviceNotFound).
    WithMessage("Устройство не найдено")

// С метаданными
err := apperror.NotFound(originalErr).
    WithCode(apperror.CodeDeviceNotFound).
    WithMessage("Устройство не найдено").
    WithDetail("device_id", "abc-123").
    WithRequestID("req-456")
```

### Валидация

```go
// Использование ValidationBuilder
v := apperror.NewValidationBuilder()

if user.Email == "" {
    v.AddRequired("email", "Email обязателен")
} else if !isValidEmail(user.Email) {
    v.AddInvalidEmail("email", user.Email)
}

if len(user.Password) < 8 {
    v.AddField("password", "too_short", "Минимум 8 символов")
}

if v.HasErrors() {
    return v.Build()
}

// Или напрямую через AppError
err := apperror.NewValidation("Ошибка валидации", nil).
    WithField("email", "required", "Email обязателен").
    WithFieldValue("age", "out_of_range", "Возраст вне диапазона", 150)
```

### JSON-ответы для HTTP API

```go
func handleError(w http.ResponseWriter, err error) {
    appErr := apperror.AsAppError(err)
    appErr.WriteJSON(w)
}

// Или с request ID
func handleErrorWithRequestID(w http.ResponseWriter, err error, requestID string) {
    apperror.HTTPErrorWithRequestID(w, err, requestID)
}
```

Пример JSON-ответа:

```json
{
    "code": "device_not_found",
    "message": "Устройство не найдено",
    "status": 404,
    "request_id": "req-abc-123",
    "details": {
        "device_id": "dev-456"
    },
    "timestamp": "2024-01-15T10:30:00Z"
}
```

### Стек-трейсы

```go
// Захват стека
err := apperror.Internal(originalErr).WithStack()

// Получение стека
if err.HasStack() {
    fmt.Println(err.StackTrace())
}

// Подробный вывод
fmt.Println(err.Verbose())

// Форматирование с %+v
fmt.Printf("%+v\n", err)
```

### Проверка типа ошибки

```go
// Проверка типа
if apperror.IsNotFound(err) {
    // обработка 404
}

if apperror.Is(err, apperror.TypeUnauthorized) {
    // обработка 401
}

// Проверка кода
if apperror.IsCode(err, apperror.CodeDeviceNotFound) {
    // специфичная обработка
}

// Извлечение данных
if code, ok := apperror.GetCode(err); ok {
    log.Printf("Error code: %s", code)
}

if status, ok := apperror.GetHTTPStatus(err); ok {
    log.Printf("HTTP status: %d", status)
}
```

### Оборачивание ошибок

```go
// Автоматическое оборачивание
appErr := apperror.Wrap(err) // Internal если не AppError

// С указанием типа
appErr := apperror.WrapWithType(err, apperror.TypeNotFound)

// С указанием кода
appErr := apperror.WrapWithCode(err, apperror.CodeDatabaseError)
```

### Множественные ошибки

```go
m := apperror.NewMultiError()

for _, item := range items {
    if err := process(item); err != nil {
        m.AddError(err)
    }
}

if m.HasErrors() {
    return m.ErrorOrNil()
}
```

### Интеграция с контекстом

```go
// Добавление данных в контекст
ctx = apperror.ContextWithRequestID(ctx, "req-123")
ctx = apperror.ContextWithUserID(ctx, "user-456")

// Создание ошибки с данными из контекста
err := apperror.FromContext(ctx, originalErr)
// err автоматически содержит request_id и user_id

// Или добавление контекста к существующей ошибке
err := apperror.NotFound(originalErr).WithContext(ctx)
```

### Логирование с slog

```go
import "log/slog"

err := apperror.NotFound(originalErr).
    WithCode(apperror.CodeDeviceNotFound).
    WithRequestID("req-123")

// Как LogValuer
slog.Error("operation failed", "error", err)

// Или получение полей
fields := err.LogFields()
slog.Error("operation failed", slog.Any("error_details", fields))
```

## Типы ошибок

| Тип | HTTP Status | Описание |
| --- | ----------- | -------- |
| `TypeNotValid` | 422 | Ошибка валидации |
| `TypeBadRequest` | 400 | Некорректный запрос |
| `TypeUnauthorized` | 401 | Не авторизован |
| `TypeForbidden` | 403 | Доступ запрещён |
| `TypeNotFound` | 404 | Не найдено |
| `TypeMethodNotAllowed` | 405 | Метод не разрешён |
| `TypeUnprocessableEntity` | 422 | Невозможно обработать |
| `TypeTooManyRequests` | 429 | Слишком много запросов |
| `TypeInternal` | 500 | Внутренняя ошибка |
| `TypeNotImplemented` | 501 | Не реализовано |
| `TypeBadGateway` | 502 | Ошибка шлюза |
| `TypeTemporaryUnavailable` | 503 | Временно недоступен |
| `TypeGatewayTimeout` | 504 | Таймаут шлюза |

## Коды ошибок

### Аутентификация

- `CodeUnauthorized` — не авторизован
- `CodeInvalidToken` — неверный токен
- `CodeTokenExpired` — токен истёк
- `CodeInvalidCredentials` — неверные учётные данные
- `CodeAccessDenied` — доступ запрещён

### Ресурсы

- `CodeNotFound` — не найдено
- `CodeResourceNotFound` — ресурс не найден
- `CodeUserNotFound` — пользователь не найден
- `CodeDeviceNotFound` — устройство не найдено
- `CodeAlreadyExists` — уже существует
- `CodeDuplicateEntry` — дубликат записи

### Валидация данных

- `CodeValidationFailed` — ошибка валидации
- `CodeInvalidInput` — неверный ввод
- `CodeRequiredField` — обязательное поле
- `CodeInvalidEmail` — неверный email
- `CodeTooShort` — слишком короткое
- `CodeTooLong` — слишком длинное
- `CodeOutOfRange` — вне диапазона

### Лимиты

- `CodeRateLimited` — превышен лимит запросов
- `CodeQuotaExceeded` — квота превышена
- `CodePayloadTooLarge` — слишком большой payload

### Инфраструктура

- `CodeServiceUnavailable` — сервис недоступен
- `CodeTimeout` — таймаут
- `CodeDatabaseError` — ошибка БД
- `CodeExternalAPIError` — ошибка внешнего API

## Структура файлов

```text
apperror/
├── errors.go      # Основная структура AppError
├── codes.go       # ErrorCode и предопределённые коды
├── validation.go  # FieldError и ValidationBuilder
├── json.go        # JSON-сериализация и HTTP
├── stack.go       # Стек-трейсы
├── helpers.go     # Утилиты (Is, Wrap, MultiError)
├── logging.go     # Интеграция с slog
├── context.go     # Интеграция с context.Context
└── errors_test.go # Тесты
```

## Лицензия

MIT
