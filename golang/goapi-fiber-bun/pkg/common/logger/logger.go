package logger

import (
	"os"

	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zlog *zap.Logger

func init() {
	var err error

	mode := os.Getenv("APP_MODE")
	var config zap.Config
	if mode == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	config.EncoderConfig = ecszap.ECSCompatibleEncoderConfig(config.EncoderConfig)
	zlog, err = config.Build(ecszap.WrapCoreOption(), zap.AddCallerSkip(1))

	if err != nil {
		panic(err)
	}
}

func WithFields(source map[string]interface{}) *zap.Logger {
	fields := ToFields(source)
	for k, v := range source {
		fields = append(fields, zap.Any(k, v))
	}
	return zlog.With(fields...)
}

func ToFields(source map[string]interface{}) []zapcore.Field {
	fields := []zapcore.Field{}
	for k, v := range source {
		fields = append(fields, zap.Any(k, v))
	}
	return fields
}

func Info(message string, fileds ...zap.Field) {
	zlog.Info(message, fileds...)
}

func Debug(message string, fileds ...zap.Field) {
	zlog.Debug(message, fileds...)
}

func Error(message string, fileds ...zap.Field) {
	zlog.Error(message, fileds...)
}

func Warn(message string, fileds ...zap.Field) {
	zlog.Warn(message, fileds...)
}

func Fatal(message string, fileds ...zap.Field) {
	zlog.Fatal(message, fileds...)
}

func Panic(message string, fileds ...zap.Field) {
	zlog.Panic(message, fileds...)
}
