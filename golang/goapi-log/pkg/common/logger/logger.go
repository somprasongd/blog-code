package logger

import (
	"os"

	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Field = zapcore.Field

type Interface interface {
	Info(msg string, fileds ...Field)
	Debug(msg string, fileds ...Field)
	Error(msg string, fileds ...Field)
	Warn(msg string, fileds ...Field)
	Fatal(msg string, fileds ...Field)
	Panic(msg string, fileds ...Field)
}

var zlog *zap.Logger
var Default Interface

func init() {
	var err error

	mode := os.Getenv("APP_MODE")
	mode = "production"

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

	Default = newZapLogger(zlog)
}

func New(fileds ...Field) Interface {
	nlog := zlog.With(fileds...)
	return newZapLogger(nlog)
}

func ToFields(source map[string]interface{}) []Field {
	fields := []Field{}
	for k, v := range source {
		fields = append(fields, zap.Any(k, v))
	}
	return fields
}
