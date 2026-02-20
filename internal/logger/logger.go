package logger

import (
	"github.com/wb-go/wbf/zlog"
)

func SetLogger(l zlog.Zerolog) {
}

func Debug(msg string, args ...interface{}) {
	zlog.Logger.Debug().Fields(args).Msg(msg)
}

func Info(msg string, args ...interface{}) {
	zlog.Logger.Info().Fields(args).Msg(msg)
}

func Warn(msg string, args ...interface{}) {
	zlog.Logger.Warn().Fields(args).Msg(msg)
}

func Error(msg string, args ...interface{}) {
	zlog.Logger.Error().Fields(args).Msg(msg)
}

func Fatal(msg string, args ...interface{}) {
	zlog.Logger.Fatal().Fields(args).Msg(msg)
}
