// Copyright 2025 Notedown Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600) // #nosec G304 - filename is provided by caller for log file creation
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
