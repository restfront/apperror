// Package apperror
package apperror

import (
	"fmt"
	"io"
	"runtime"
	"strings"
)

// Frame представляет один кадр стека вызовов.
type Frame struct {
	Function string
	File     string
	Line     int
}

// String возвращает строковое представление кадра.
func (f Frame) String() string {
	return fmt.Sprintf("%s\n\t%s:%d", f.Function, f.File, f.Line)
}

// Stack возвращает стек вызовов ошибки.
func (e *AppError) Stack() []uintptr {
	return e.stack
}

// HasStack проверяет, есть ли у ошибки стек вызовов.
func (e *AppError) HasStack() bool {
	return len(e.stack) > 0
}

// WithStack захватывает текущий стек вызовов.
func (e *AppError) WithStack() *AppError {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(2, pcs[:])
	e.stack = make([]uintptr, n)
	copy(e.stack, pcs[:n])
	return e
}

// StackFrames возвращает стек вызовов в виде структурированных кадров.
func (e *AppError) StackFrames() []Frame {
	if len(e.stack) == 0 {
		return nil
	}

	frames := runtime.CallersFrames(e.stack)
	result := make([]Frame, 0, len(e.stack))

	for {
		frame, more := frames.Next()
		result = append(result, Frame{
			Function: frame.Function,
			File:     frame.File,
			Line:     frame.Line,
		})
		if !more {
			break
		}
	}

	return result
}

// StackTrace возвращает форматированный стек вызовов.
func (e *AppError) StackTrace() string {
	if len(e.stack) == 0 {
		return ""
	}

	var buf strings.Builder
	frames := runtime.CallersFrames(e.stack)

	for {
		frame, more := frames.Next()
		fmt.Fprintf(&buf, "%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}

	return buf.String()
}

// Format реализует интерфейс fmt.Formatter.
// Поддерживает следующие форматы:
//   - %s, %v: выводит сообщение об ошибке
//   - %+v: выводит сообщение об ошибке и стек вызовов
func (e *AppError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, e.Error())
			if e.HasStack() {
				io.WriteString(s, "\n")
				io.WriteString(s, e.StackTrace())
			}
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}

// Verbose возвращает подробное описание ошибки, включая стек.
func (e *AppError) Verbose() string {
	var buf strings.Builder

	buf.WriteString("Error: ")
	buf.WriteString(e.Error())
	buf.WriteString("\n")

	buf.WriteString("Type: ")
	buf.WriteString(e.errorType.String())
	buf.WriteString("\n")

	if !e.code.IsEmpty() {
		buf.WriteString("Code: ")
		buf.WriteString(string(e.code))
		buf.WriteString("\n")
	}

	buf.WriteString("Message: ")
	buf.WriteString(e.Message())
	buf.WriteString("\n")

	buf.WriteString("HTTP Status: ")
	fmt.Fprintf(&buf, "%d", e.HTTPStatusCode())
	buf.WriteString("\n")

	if e.requestID != "" {
		buf.WriteString("Request ID: ")
		buf.WriteString(e.requestID)
		buf.WriteString("\n")
	}

	if !e.timestamp.IsZero() {
		buf.WriteString("Timestamp: ")
		buf.WriteString(e.timestamp.String())
		buf.WriteString("\n")
	}

	if len(e.details) > 0 {
		buf.WriteString("Details:\n")
		for k, v := range e.details {
			fmt.Fprintf(&buf, "  %s: %v\n", k, v)
		}
	}

	if len(e.fields) > 0 {
		buf.WriteString("Field Errors:\n")
		for _, f := range e.fields {
			fmt.Fprintf(&buf, "  %s (%s): %s\n", f.Field, f.Code, f.Message)
		}
	}

	if e.HasStack() {
		buf.WriteString("Stack Trace:\n")
		buf.WriteString(e.StackTrace())
	}

	return buf.String()
}
