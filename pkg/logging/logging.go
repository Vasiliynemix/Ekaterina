package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

const (
	debug = "debug"
	info  = "info"
	prod  = "prod"
)

type Logger struct {
	*zap.Logger
}

func NewLogger(
	logger *zap.Logger,
) *Logger {
	return &Logger{
		Logger: logger,
	}
}

func InitLogger(
	env string,
	structDateFormat string,
	pathToInfoLogs string,
	pathToDebugLogs string,
) *zap.Logger {
	logger := getFileLogger(env, structDateFormat, pathToInfoLogs, pathToDebugLogs)
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)
	return logger
}

func getFileLogger(
	env string,
	structDateFormat string,
	pathToInfoLogs string,
	pathToDebugLogs string,
) *zap.Logger {
	cfgLogger := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		CallerKey:      "caller",
		MessageKey:     "message",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout(structDateFormat),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	fileEncoder := zapcore.NewJSONEncoder(cfgLogger)
	consoleEncoder := zapcore.NewConsoleEncoder(cfgLogger)
	consoleWriter := zapcore.AddSync(os.Stdout)

	logFileInfo, _ := os.OpenFile(pathToInfoLogs, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	logFileDebug, _ := os.OpenFile(pathToDebugLogs, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	logLevelCfg := debug

	if env == prod {
		logLevelCfg = info
	}

	var writer zapcore.WriteSyncer
	var logLevel zapcore.Level

	switch logLevelCfg {
	case debug:
		writer = zapcore.AddSync(logFileDebug)
		logLevel = zapcore.DebugLevel
	case info:
		writer = zapcore.AddSync(logFileInfo)
		logLevel = zapcore.InfoLevel
	}

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, logLevel),
		zapcore.NewCore(consoleEncoder, consoleWriter, logLevel),
	)

	log := zap.New(core, zap.AddCaller())

	return log
}
