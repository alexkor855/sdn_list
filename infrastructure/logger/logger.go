package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxLogger struct{}

func InitLogger(serviceName string, logLevel zapcore.Level) (*zap.Logger, error) {
	logConfig := zap.NewDevelopmentConfig()

	logConfig.Level = zap.NewAtomicLevelAt(logLevel)
	logConfig.EncoderConfig.LineEnding = "\n\n"

	logger, err := logConfig.Build()
	if err != nil {
		return nil, err
	}

	return logger.Named(serviceName), nil
}

// adds logger to context
func ContextWithLogger(ctx context.Context, l *zap.Logger) context.Context {
    return context.WithValue(ctx, ctxLogger{}, l)
}

// returns logger from context
func LoggerFromContext(ctx context.Context) *zap.Logger {
    if l, ok := ctx.Value(ctxLogger{}).(*zap.Logger); ok {
        return l
    }
    return zap.L()
}
