package utils

import "log/slog"

type Logger interface {
	Info(msg, source, method, path, agent string, data any)
	Error(msg, source, method, path, agent string, data any)
}

type SlogLogger struct{}

func (l *SlogLogger) Info(msg, source, method, path, agent string, data any) {
	slog.Info(
		msg,
		"source", source,
		"method", method,
		"path", path,
		"user_agent", agent,
		"data", data,
	)
}

func (l *SlogLogger) Error(msg, source, method, path, agent string, data any) {
	slog.Error(
		msg,
		"source", source,
		"method", method,
		"path", path,
		"user_agent", agent,
		"data", data,
	)
}
