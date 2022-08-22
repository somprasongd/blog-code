package logger

import "go.uber.org/zap"

type zapLogger struct {
	*zap.Logger
}

func newZapLogger(logger *zap.Logger) Interface {
	return &zapLogger{logger}
}

func (l *zapLogger) Info(msg string, fileds ...Field) {
	l.Logger.Info(msg, fileds...)
}

func (l *zapLogger) Debug(msg string, fileds ...Field) {
	l.Logger.Debug(msg, fileds...)
}

func (l *zapLogger) Error(msg string, fileds ...Field) {
	l.Logger.Error(msg, fileds...)
}

func (l *zapLogger) Warn(msg string, fileds ...Field) {
	l.Logger.Warn(msg, fileds...)
}

func (l *zapLogger) Fatal(msg string, fileds ...Field) {
	l.Logger.Fatal(msg, fileds...)
}

func (l *zapLogger) Panic(msg string, fileds ...Field) {
	l.Logger.Panic(msg, fileds...)
}
