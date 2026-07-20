// Package logger provides the application's structured logger.
package logger

import (
	"context"
	"errors"
	"fmt"
	"os"
	"syscall"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/asl-open/asl-core/pkg/config"
	"github.com/asl-open/asl-core/pkg/requestid"
)

var Module = fx.Provide(New)

// Logger's methods take ctx so the request ID (see pkg/requestid) is
// attached to every log line automatically, without every call site
// passing it explicitly.
type Logger interface {
	Info(ctx context.Context, msg string, fields ...any)
	Debug(ctx context.Context, msg string, fields ...any)
	Warn(ctx context.Context, msg string, fields ...any)
	Error(ctx context.Context, msg string, fields ...any)
}

type Params struct {
	fx.In
	fx.Lifecycle

	Config config.Config
}

type logger struct {
	lg *zap.SugaredLogger
}

func New(p Params) (Logger, error) {
	level := getLevel(p.Config)

	var encoder zapcore.Encoder
	if p.Config.GetString("logger.format") == "json" {
		encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	} else {
		encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	}

	core := zapcore.NewTee(
		zapcore.NewCore(
			encoder,
			zapcore.Lock(os.Stdout),
			level,
		))

	log := zap.New(core, zap.AddCaller())

	p.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				if err := log.Sync(); err != nil && !isIgnorableSyncError(err) {
					return fmt.Errorf("failed to sync logger: %w", err)
				}
				return nil
			},
		},
	)

	return &logger{lg: log.Sugar()}, nil
}

func (l *logger) Info(ctx context.Context, msg string, fields ...any) {
	l.lg.Infow(msg, withRequestID(ctx, fields)...)
}

func (l *logger) Debug(ctx context.Context, msg string, fields ...any) {
	l.lg.Debugw(msg, withRequestID(ctx, fields)...)
}

func (l *logger) Warn(ctx context.Context, msg string, fields ...any) {
	l.lg.Warnw(msg, withRequestID(ctx, fields)...)
}

func (l *logger) Error(ctx context.Context, msg string, fields ...any) {
	l.lg.Errorw(msg, withRequestID(ctx, fields)...)
}

// withRequestID appends a "request_id" field to fields if ctx carries
// one (see pkg/requestid), without mutating fields' underlying array.
func withRequestID(ctx context.Context, fields []any) []any {
	id, ok := requestid.FromContext(ctx)
	if !ok {
		return fields
	}

	out := make([]any, 0, len(fields)+2)
	out = append(out, fields...)
	return append(out, "request_id", id)
}

// isIgnorableSyncError reports whether err is one of the well-known cases
// where Sync() on stdout/stderr fails even though nothing is actually
// wrong - e.g. inside a container, stdout is a pipe/character device, and
// fsync on it returns EINVAL or ENOTTY. Observed as "sync /dev/stdout:
// invalid argument" when running under Docker.
func isIgnorableSyncError(err error) bool {
	return errors.Is(err, syscall.EINVAL) || errors.Is(err, syscall.ENOTTY)
}

func getLevel(cfg config.Config) zapcore.Level {
	switch cfg.GetString("logger.level") {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.DebugLevel
	}
}
