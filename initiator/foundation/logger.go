package foundation

import (
	"fmt"
	"os"

	"context"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

type Log struct {
	logger *zap.Logger
}

func InitLogger() *Log {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	config.EncoderConfig.NameKey = "logger"
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.LevelKey = "level"

	logger, err := config.Build()
	if err != nil {
		fmt.Printf(`{level:"fatal,"msg":"failed to initialize logger: %v"}`, err)
		os.Exit(1)
	}

	return &Log{
		logger: logger,
	}
}

func (l *Log) GetLogger() *zap.Logger {
	return l.logger
}

// Logger is the interface for the logger.
type Logger interface {
	// GetZapLogger returns the underlying zap logger.
	GetZapLogger() *zap.Logger

	// Named returns a new named logger.
	Named(s string) *logger

	// With returns a new logger with the given fields.
	With(fields ...zap.Field) *logger

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

type logger struct {
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

	return &logger{
		logger:  l,
		options: options,
	}
}

// GetZapLogger returns the underlying zap logger.
func (l *logger) GetZapLogger() *zap.Logger {
	return l.logger
}

// Named returns a new named logger.
func (l *logger) Named(s string) *logger {
	l2 := l.logger.Named(s)

	return &logger{
		logger:  l2,
		options: l.options,
	}
}

// With returns a new logger with the given fields.
func (l *logger) With(fields ...zap.Field) *logger {
	l2 := l.logger.With(fields...)

	return &logger{
		logger:  l2,
		options: l.options,
	}
}

// Debug logs a debug message.
func (l *logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.With(l.extract(ctx)...).Debug(msg, fields...)
}

// Info logs an info message.
func (l *logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.With(l.extract(ctx)...).Info(msg, fields...)
}

// Warn logs a warning message.
func (l *logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.With(l.extract(ctx)...).Warn(msg, fields...)
}

// Error logs an error message with stack trace.
func (l *logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.With(l.extract(ctx)...).Error(msg, fields...)
}

// Panic logs a panic message and panics.
func (l *logger) Panic(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.With(l.extract(ctx)...).Panic(msg, fields...)
}

// Fatal logs a fatal message and exits with os.Exit(1).
func (l *logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.With(l.extract(ctx)...).Fatal(msg, fields...)
}

func (l *logger) extract(ctx context.Context) []zap.Field {
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

// Printf is the kafka logger function implementation.
func (l *logger) Printf(msg string, fields ...interface{}) {
	l.Info(context.Background(), fmt.Sprintf(msg, fields...))
}

// Log is an implementation for pgx logger.
//
// Note:
// the time field is remapped to pgx_time to avoid conflicts with the logger time field.
// the args field is formatted to []string to insure a valid json encoding.
func (l *logger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	fields := make([]zap.Field, 0, len(data))

	// since the logger might have a `time` field, we have to remap the pgx time
	data["pgx_time"] = data["time"]
	delete(data, "time")

	for k, v := range data {
		// format args values to []string
		// this is to insure a valid json encoding
		if k == "args" {
			if args, ok := v.([]interface{}); ok {
				var argsStr []string

				for _, arg := range args {
					if argByte, ok := arg.(pgtype.JSON); ok {
						arg = string(argByte.Bytes)
					}

					argsStr = append(argsStr, fmt.Sprintf("%v", arg))
				}

				v = argsStr
			}
		}

		fields = append(fields, zap.Any(k, v))
	}

	switch level {
	case pgx.LogLevelInfo:
		l.Info(ctx, msg, fields...)
	case pgx.LogLevelWarn:
		l.Warn(ctx, msg, fields...)
	case pgx.LogLevelError:
		l.Error(ctx, msg, fields...)
	default:
		l.Debug(ctx, msg, fields...)
	}
}
