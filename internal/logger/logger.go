package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func init() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if os.Getenv("ENV") == "development" {
		config = zap.NewDevelopmentConfig()
	}

	var err error
	log, err = config.Build()
	if err != nil {
		panic(err)
	}
}

func Info(msg string, fields ...zapcore.Field) {
	log.Info(msg, fields...)
}

func Error(msg string, fields ...zapcore.Field) {
	log.Error(msg, fields...)
}

func Debug(msg string, fields ...zapcore.Field) {
	log.Debug(msg, fields...)
}

func Fatal(msg string, fields ...zapcore.Field) {
	log.Fatal(msg, fields...)
}

func With(fields ...zapcore.Field) *zap.Logger {
	return log.With(fields...)
}
