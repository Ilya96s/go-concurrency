package logger

import (
	"errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"in-memory-kv/internal/config"
	"os"
	"strings"
)

// Создание кастомного логера
func NewLogger(cfg config.Config) (*zap.Logger, error) {
	//Уровень логирования
	var level zapcore.Level
	switch strings.ToLower(cfg.Logging.Level) {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn", "warning":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		return nil, errors.New("unknown logging level: " + cfg.Logging.Level)
	}

	// настройка формата вывода
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.CallerKey = "caller"
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// куда писать (stdout)
	var writer zapcore.WriteSyncer
	switch strings.ToLower(cfg.Logging.Output) {
	case "stdout", "":
		writer = zapcore.AddSync(os.Stdout)
	case "stderr":
		writer = zapcore.AddSync(os.Stderr)
	default:
		f, err := os.OpenFile(cfg.Logging.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return nil, err
		}
		writer = zapcore.AddSync(f)
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writer,
		level,
	)

	// Опция для логирования места вызова
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return logger, nil
}
