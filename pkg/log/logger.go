package log

import (
	"io"
	"log/slog"
	"os"
)

type Logger struct {
	slog *slog.Logger
}

func NewFile(filename string, level Level, format Format) (*Logger, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return NewWithFormat(file, level, format), nil
}

func NewLsp(level Level, format Format) *Logger {
	return NewWithFormat(os.Stderr, level, format)
}

func New(w io.Writer, level Level) *Logger {
	return NewWithFormat(w, level, FormatText)
}

func NewWithFormat(w io.Writer, level Level, format Format) *Logger {
	opts := &slog.HandlerOptions{
		Level: level.ToSlogLevel(),
	}

	var handler slog.Handler
	switch format {
	case FormatJSON:
		handler = slog.NewJSONHandler(w, opts)
	default:
		handler = slog.NewTextHandler(w, opts)
	}

	return &Logger{
		slog: slog.New(handler),
	}
}

func NewDefault() *Logger {
	return New(os.Stderr, Info)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.slog.Debug(msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.slog.Info(msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.slog.Warn(msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.slog.Error(msg, args...)
}

func (l *Logger) With(args ...any) *Logger {
	return &Logger{
		slog: l.slog.With(args...),
	}
}

func (l *Logger) WithScope(scope string) *Logger {
	return &Logger{
		slog: l.slog.With("scope", scope),
	}
}
