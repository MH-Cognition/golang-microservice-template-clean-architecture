package logger

import (
	"context"
	"log/slog"
	"os"
)

type SlogLogger struct {
	log *slog.Logger
	env string
}

func New(env string) *SlogLogger {
	var handler slog.Handler

	if env == "prod" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelWarn, // 🔴 only WARN & ERROR
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug, // 🔵 INFO + DEBUG
		})
	}

	return &SlogLogger{
		log: slog.New(handler),
		env: env,
	}
}

func (l *SlogLogger) Debug(ctx context.Context, msg string, fields ...any) {
	l.log.DebugContext(ctx, msg, fields...)
}

func (l *SlogLogger) Info(ctx context.Context, msg string, fields ...any) {
	l.log.InfoContext(ctx, msg, fields...)
}

func (l *SlogLogger) Warn(ctx context.Context, msg string, fields ...any) {
	l.log.WarnContext(ctx, msg, fields...)
}

func (l *SlogLogger) Error(ctx context.Context, msg string, fields ...any) {
	l.log.ErrorContext(ctx, msg, fields...)
}
