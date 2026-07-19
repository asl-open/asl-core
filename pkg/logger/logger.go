// Package logger provides the application's structured logger.
package logger

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/asl-open/asl-core/pkg/config"
)

var Module = fx.Provide(New)

// Logger's methods take ctx so per-request fields (e.g. the request ID
// added in #20) can be attached automatically once that lands - for now ctx
// is unused.
type Logger interface {
	Info(ctx context.Context, msg string, fields ...interface{})
	Debug(ctx context.Context, msg string, fields ...interface{})
	Warn(ctx context.Context, msg string, fields ...interface{})
	Error(ctx context.Context, msg string, fields ...interface{})
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
				if err := log.Sync(); err != nil {
					return fmt.Errorf("failed to sync logger: %w", err)
				}
				return nil
			},
		},
	)

	return &logger{lg: log.Sugar()}, nil
}

func (l *logger) Info(ctx context.Context, msg string, fields ...interface{}) {
	l.lg.Infow(msg, fields...)
}

func (l *logger) Debug(ctx context.Context, msg string, fields ...interface{}) {
	l.lg.Debugw(msg, fields...)
}

func (l *logger) Warn(ctx context.Context, msg string, fields ...interface{}) {
	l.lg.Warnw(msg, fields...)
}

func (l *logger) Error(ctx context.Context, msg string, fields ...interface{}) {
	l.lg.Errorw(msg, fields...)
}

func getLevel(config config.Config) zapcore.Level {
	switch config.GetString("logger.level") {
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
