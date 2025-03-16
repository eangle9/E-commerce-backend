package logit

import (
	"context"
	"fmt"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

// Printf is the kafka logger function implementation.
func (l *Logit) Printf(msg string, fields ...interface{}) {
	l.Info(context.Background(), fmt.Sprintf(msg, fields...))
}

// Log is an implementation for pgx logger.
//
// Note:
// the time field is remapped to pgx_time to avoid conflicts with the logger time field.
// the args field is formatted to []string to insure a valid json encoding.
func (l *Logit) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
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
