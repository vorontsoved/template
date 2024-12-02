package logging

import (
	"context"
	"log/slog"
	"os"
	"{{cookiecutter.project_slug}}/internal/pkg/constants"

	"github.com/google/uuid"
)

type Logger interface {
	Info(ctx context.Context, msg string, keysAndValues ...any)
	InfoWithoutContext(msg string, keysAndValues ...any)
	Error(ctx context.Context, msg string, err error, keysAndValues ...any)
	ErrorWithoutContext(msg string, err error, keysAndValues ...any)
	Debug(ctx context.Context, msg string, keysAndValues ...any)
	DebugWithoutContext(msg string, keysAndValues ...any)
	Warn(ctx context.Context, msg string, keysAndValues ...any)
	WarnWithoutContext(msg string, keysAndValues ...any)
	With(keysAndValues ...any) Logger
	WithLayer(layer string) Logger
}

type SlogWrapper struct {
	logger *slog.Logger
}

func NewSlogWrapper(baseLogger *slog.Logger) Logger {
	return &SlogWrapper{logger: baseLogger}
}

func (s *SlogWrapper) Info(ctx context.Context, msg string, keysAndValues ...any) {
	s.logger.InfoContext(ctx, msg, keysAndValues...)
}

func (s *SlogWrapper) InfoWithoutContext(msg string, keysAndValues ...any) {
	s.logger.Info(msg, keysAndValues...)
}

func (s *SlogWrapper) Error(ctx context.Context, msg string, err error, keysAndValues ...any) {
	s.logger.ErrorContext(ctx, msg, append(keysAndValues, "error", err)...)
}

func (s *SlogWrapper) ErrorWithoutContext(msg string, err error, keysAndValues ...any) {
	s.logger.Error(msg, append(keysAndValues, "error", err)...)
}

func (s *SlogWrapper) Debug(ctx context.Context, msg string, keysAndValues ...any) {
	s.logger.DebugContext(ctx, msg, keysAndValues...)
}

func (s *SlogWrapper) DebugWithoutContext(msg string, keysAndValues ...any) {
	s.logger.Debug(msg, keysAndValues...)
}

func (s *SlogWrapper) Warn(ctx context.Context, msg string, keysAndValues ...any) {
	s.logger.WarnContext(ctx, msg, keysAndValues...)
}

func (s *SlogWrapper) WarnWithoutContext(msg string, keysAndValues ...any) {
	s.logger.Warn(msg, keysAndValues...)
}

func (s *SlogWrapper) With(keysAndValues ...any) Logger {
	return &SlogWrapper{logger: s.logger.With(keysAndValues...)}
}

func (s *SlogWrapper) WithLayer(layer string) Logger {
	return &SlogWrapper{logger: s.logger.With("layer", layer)}
}

func GetLoggerFromContext(ctx context.Context, defaultLogger Logger) Logger {
	if logger, ok := ctx.Value(constants.LoggerKey).(Logger); ok {
		return logger
	}
	return defaultLogger
}

func NewBaseLogger(level slog.Level) *slog.Logger {
	opts := slog.HandlerOptions{Level: level}
	handler := slog.NewJSONHandler(os.Stdout, &opts)
	return slog.New(handler)
}

func WithRequestID(logger Logger) Logger {
	return logger.With("requestID", uuid.New().String())
}
