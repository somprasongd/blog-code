package logger

import (
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
)

var log *zap.Logger

func init() {
	var err error

	config := zap.NewProductionConfig()
	config.EncoderConfig = ecszap.ECSCompatibleEncoderConfig(config.EncoderConfig)
	log, err = config.Build(ecszap.WrapCoreOption(), zap.AddCallerSkip(1))

	if err != nil {
		panic(err)
	}
}

func Info(message string, fileds ...zap.Field) {
	log.Info(message, fileds...)
}

func Debug(message string, fileds ...zap.Field) {
	log.Debug(message, fileds...)
}

func Error(message string, fileds ...zap.Field) {
	log.Error(message, fileds...)
}
