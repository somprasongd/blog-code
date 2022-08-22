package logger

import (
	"os"

	"go.elastic.co/ecszap"
	"go.uber.org/zap"
)

var log *zap.Logger

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
	log, err = config.Build(ecszap.WrapCoreOption(), zap.AddCallerSkip(1))

	if err != nil {
		panic(err)
	}
}

func Info(message string, fileds ...zap.Field) {
	log.Info(message, fileds...)
}

func InfoWithReqId(message string, reqId string, fileds ...zap.Field) {
	filed := zap.String("requestId", reqId)
	fileds = append(fileds, filed)
	log.Info(message, fileds...)
}

func Debug(message string, fileds ...zap.Field) {
	log.Debug(message, fileds...)
}

func DebugWithReqId(message string, reqId string, fileds ...zap.Field) {
	filed := zap.String("requestId", reqId)
	fileds = append(fileds, filed)
	log.Debug(message, fileds...)
}

func Error(message string, fileds ...zap.Field) {
	log.Error(message, fileds...)
}

func ErrorWithReqId(message string, reqId string, fileds ...zap.Field) {
	filed := zap.String("requestId", reqId)
	fileds = append(fileds, filed)
	log.Error(message, fileds...)
}

func Fatal(message string, fileds ...zap.Field) {
	log.Error(message, fileds...)
	panic(message)
}
