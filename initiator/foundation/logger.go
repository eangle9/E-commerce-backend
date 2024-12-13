package foundation

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.Logger
}

func InitLogger() *Logger {
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

	return &Logger{
		logger: logger,
	}
}

func (l *Logger) GetLogger() *zap.Logger {
	return l.logger
}
