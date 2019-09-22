package log

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LoggingFactory wrapper to hold
type LoggingFactory struct {
	logger *zap.Logger
}

// NewLoggerFactory creates a logger factory
func NewLoggerFactory(logger *zap.Logger) LoggingFactory {
	return LoggingFactory{logger: logger}
}

// Bg creates a context-unaware logger.
func (b LoggingFactory) Bg() Logger {
	return logger(b)
}

// For returns a context-aware Logger. If the context
// contains an OpenTracing span, all logging calls are also
func (b LoggingFactory) For(ctx context.Context) Logger {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		return spanLogger{span: span, logger: b.logger}
	}
	return b.Bg()
}

// With creates a child logger, and optionally adds some context fields to that logger.
func (b LoggingFactory) With(fields ...zapcore.Field) LoggingFactory {
	return LoggingFactory{logger: b.logger.With(fields...)}
}
