package log

import "log/slog"

type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
)

func (l Level) String() string {
	switch l {
	case Debug:
		return "debug"
	case Info:
		return "info"
	case Warn:
		return "warn"
	case Error:
		return "error"
	default:
		return "unknown"
	}
}

func (l Level) ToSlogLevel() slog.Level {
	switch l {
	case Debug:
		return slog.LevelDebug
	case Info:
		return slog.LevelInfo
	case Warn:
		return slog.LevelWarn
	case Error:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func ParseLevel(s string) Level {
	switch s {
	case "debug":
		return Debug
	case "info":
		return Info
	case "warn":
		return Warn
	case "error":
		return Error
	default:
		return Info
	}
}

type Format int

const (
	FormatText Format = iota
	FormatJSON
)

func ParseFormat(s string) Format {
	switch s {
	case "json":
		return FormatJSON
	case "text":
		return FormatText
	default:
		return FormatText
	}
}