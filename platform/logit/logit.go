package logit

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

// Logger is the interface for the logger.
type Logger interface {
	// GetZapLogger returns the underlying zap logger.
	GetZapLogger() *zap.Logger

	// Named returns a new named logger.
	Named(s string) *Logit

	// With returns a new logger with the given fields.
	With(fields ...zap.Field) *Logit

	// Debug logs a debug message.
	Debug(ctx context.Context, msg string, fields ...zap.Field)

	// Info logs an info message.
	Info(ctx context.Context, msg string, fields ...zap.Field)

	// Warn logs a warning message.
	Warn(ctx context.Context, msg string, fields ...zap.Field)

	// Error logs an error message with stack trace.
	Error(ctx context.Context, msg string, fields ...zap.Field)

	// Panic logs a panic message and panics.
	Panic(ctx context.Context, msg string, fields ...zap.Field)

	// Fatal logs a fatal message and exits with os.Exit(1).
	Fatal(ctx context.Context, msg string, fields ...zap.Field)

	// Log is an implementation for pgx logger.
	//
	// Note:
	// the time field is remapped to pgx_time to avoid conflicts with the logger time field.
	// the args field is formatted to []string to insure a valid json encoding.
	Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{})

	extract(ctx context.Context) []zap.Field
}

// ExtractField is a struct to define the transformation of a field in the context.
type ExtractField struct {
	// KeyInContext is the value of the key of the field in the context.
	KeyInContext any
	// Func is the transformation function to be applied to the value of the field in the context.
	Func func(any) zap.Field
}

// Options is a struct to define the options for the logger.
type Options struct {
	// ExtractFields is a list of fields to extract from the context.
	ExtractFields []ExtractField
	// IgnoreDefaultExtractFields is a flag to ignore the default extract fields.
	IgnoreDefaultExtractFields bool
	// IgnoreDefaultTimeField is a flag to ignore the default time field.
	IgnoreDefaultTimeField bool
}

type Logit struct {
	logger  *zap.Logger
	options Options
}

// New initializes a new logger from a zap logger instance.
func New(l *zap.Logger, options Options) Logger {
	if !options.IgnoreDefaultExtractFields {
		options.ExtractFields = append(options.ExtractFields, []ExtractField{
			{
				KeyInContext: "x-request-id",
				Func: func(v any) zap.Field {
					if vString, ok := v.(string); ok {
						return zap.String("x-request-id", vString)
					}

					return zap.Skip()
				},
			},
			{
				KeyInContext: "x-user-id",
				Func: func(v any) zap.Field {
					if vString, ok := v.(string); ok {
						return zap.String("x-user-id", vString)
					}

					return zap.Skip()
				},
			},
			{
				KeyInContext: "request-start-time",
				Func: func(v any) zap.Field {
					if vTime, ok := v.(time.Time); ok {
						return zap.Float64("time-since-request", float64(time.Since(vTime).Milliseconds()))
					}

					return zap.Skip()
				},
			},
			{
				KeyInContext: "x-ws-request-id",
				Func: func(v any) zap.Field {
					if vString, ok := v.(string); ok {
						return zap.String("x-ws-request-id", vString)
					}

					return zap.Skip()
				},
			},
		}...)
	}

	return &Logit{
		logger:  l,
		options: options,
	}
}

// GetZapLogger returns the underlying zap logger.
func (l *Logit) GetZapLogger() *zap.Logger {
	return l.logger
}

// Named returns a new named logger.
func (l *Logit) Named(s string) *Logit {
	l2 := l.logger.Named(s)

	return &Logit{
		logger:  l2,
		options: l.options,
	}
}

// With returns a new logger with the given fields.
func (l *Logit) With(fields ...zap.Field) *Logit {
	l2 := l.logger.With(fields...)

	return &Logit{
		logger:  l2,
		options: l.options,
	}
}

// Debug logs a debug message.
func (l *Logit) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.With(l.extract(ctx)...).Debug(msg, fields...)
}

// Info logs an info message.
func (l *Logit) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.With(l.extract(ctx)...).Info(msg, fields...)
}

// Warn logs a warning message.
func (l *Logit) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.With(l.extract(ctx)...).Warn(msg, fields...)
}

// Error logs an error message with stack trace.
func (l *Logit) Error(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.With(l.extract(ctx)...).Error(msg, fields...)
}

// Panic logs a panic message and panics.
func (l *Logit) Panic(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.With(l.extract(ctx)...).Panic(msg, fields...)
}

// Fatal logs a fatal message and exits with os.Exit(1).
func (l *Logit) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.With(l.extract(ctx)...).Fatal(msg, fields...)
}

func (l *Logit) extract(ctx context.Context) []zap.Field {
	var fields []zap.Field

	if !l.options.IgnoreDefaultTimeField {
		fields = append(fields, zap.String("time", time.Now().Format(time.RFC3339)))
	}

	if ctx != nil {
		for _, field := range l.options.ExtractFields {
			if v := ctx.Value(field.KeyInContext); v != nil {
				fields = append(fields, field.Func(v))
			}
		}
	}

	return fields
}
