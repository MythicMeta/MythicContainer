package logging

import (
	"errors"
	"os"
	"runtime"
	"time"

	"github.com/MythicMeta/MythicContainer/config"
	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog"
)

var (
	logger logr.Logger
)

func init() {
	zerologr.NameFieldName = "logger"
	zerologr.NameSeparator = "/"
	var zl zerolog.Logger
	switch config.MythicConfig.DebugLevel {
	case "warning":
		zl = zerolog.New(os.Stdout).Level(zerolog.WarnLevel)
	case "info":
		zl = zerolog.New(os.Stdout).Level(zerolog.InfoLevel)
	case "debug":
		zl = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).Level(zerolog.DebugLevel)
	case "trace":
		zl = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).Level(zerolog.TraceLevel)
	default:
		zl = zerolog.New(os.Stdout).Level(zerolog.InfoLevel)
	}
	zl = zl.With().Timestamp().Logger()
	logger = zerologr.New(&zl)
	logger.Info("Logging Initialized")
}

func LogFatalError(err error, message string, messages ...interface{}) {
	LogError(err, message, messages...)
	os.Exit(1)
}

func LogWarning(message string, messages ...interface{}) {
	if pc, _, line, ok := runtime.Caller(1); ok {
		logger.V(-1).Info(message, append([]interface{}{"func", runtime.FuncForPC(pc).Name(), "line", line}, messages...)...)
	} else {
		logger.V(-1).Info(message, messages...)
	}
}

func LogTrace(message string, messages ...interface{}) {
	if pc, _, line, ok := runtime.Caller(1); ok {
		logger.V(2).Info(message, append([]interface{}{"func", runtime.FuncForPC(pc).Name(), "line", line}, messages...)...)
	} else {
		logger.V(2).Info(message, messages...)
	}
}

func LogDebug(message string, messages ...interface{}) {
	if pc, _, line, ok := runtime.Caller(1); ok {
		logger.V(1).Info(message, append([]interface{}{"func", runtime.FuncForPC(pc).Name(), "line", line}, messages...)...)
	} else {
		logger.V(1).Info(message, messages...)
	}
}

func LogInfo(message string, messages ...interface{}) {
	if pc, _, line, ok := runtime.Caller(1); ok {
		logger.V(0).Info(message, append([]interface{}{"func", runtime.FuncForPC(pc).Name(), "line", line}, messages...)...)
	} else {
		logger.V(0).Info(message, messages...)
	}
}

func LogError(err error, message string, messages ...interface{}) {
	if pc, _, line, ok := runtime.Caller(1); ok {
		if err == nil {
			logger.Error(errors.New(message), "", messages...)
		} else {
			logger.Error(err, message, append([]interface{}{"func", runtime.FuncForPC(pc).Name(), "line", line}, messages...)...)
		}
	} else {
		if err == nil {
			logger.Error(errors.New(message), "", messages...)
		} else {
			logger.Error(err, message, messages...)
		}
	}
}
