package log

import (
	"fmt"
	"github.com/ohwin/micro/core/config"
	. "github.com/ohwin/micro/core/constant"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

var encoderDevConfig = zapcore.EncoderConfig{
	TimeKey:        "time",
	LevelKey:       "level",
	NameKey:        "logger",
	CallerKey:      "caller",
	FunctionKey:    zapcore.OmitKey,
	MessageKey:     "msg",
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.CapitalColorLevelEncoder,
	EncodeTime:     zapcore.TimeEncoderOfLayout("[2006-01-02T15:04:05.000]"),
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeCaller:   callerEncoder,
}

var encoderProdConfig = zapcore.EncoderConfig{
	TimeKey:        "time",
	LevelKey:       "level",
	NameKey:        "logger",
	CallerKey:      "caller",
	FunctionKey:    zapcore.OmitKey,
	MessageKey:     "msg",
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.CapitalLevelEncoder,
	EncodeTime:     zapcore.TimeEncoderOfLayout("[2006-01-02T15:04:05.000]"),
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeCaller:   callerEncoder,
}

func Init() {
	zap.NewProductionConfig()
	logger = zap.New(core(), zap.AddCaller(), zap.AddCallerSkip(2))
}

func core() zapcore.Core {
	var level zapcore.Level
	var encoder zapcore.Encoder

	switch config.App.Env {
	case Release:
		level = zap.InfoLevel
		encoder = zapcore.NewJSONEncoder(encoderProdConfig)
	default:
		level = zap.DebugLevel
		encoder = zapcore.NewConsoleEncoder(encoderDevConfig)
	}
	return zapcore.NewCore(encoder, os.Stdout, level)
}

func callerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(caller.FullPath())
}

func Debug(a ...any) {
	write(zap.DebugLevel, a)
}

func Error(a ...any) {
	write(zap.ErrorLevel, a)
}

func Info(a ...any) {
	write(zap.InfoLevel, a)
}
func Warn(a ...any) {
	write(zap.WarnLevel, a)
}

func Debugf(format string, a ...any) {
	write(zap.DebugLevel, fmt.Sprintf(format, a...))
}

func Errorf(format string, a ...any) {
	write(zap.ErrorLevel, fmt.Sprintf(format, a...))
}

func Infof(format string, a ...any) {
	write(zap.InfoLevel, fmt.Sprintf(format, a...))
}

func Warnf(format string, a ...any) {
	write(zap.WarnLevel, fmt.Sprintf(format, a...))
}

func write(level zapcore.Level, a ...any) {
	logger.Log(level, fmt.Sprint(a...))
}
