package util

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"os"
)

type LogType int

const (
	LogInfo LogType = iota + 1
	LogWarn
	LogErr
	LogData
)

func (lt LogType) String() string {
	return []string{"LogInfo", "LogWarn", "LogErr", "LogData"}[lt]
}

func CreateStdGoKitLog(serviceName string, debug bool) log.Logger {
	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.NewSyncLogger(logger)
	logger = log.With(
		logger,
		"service", serviceName,
		"time", log.DefaultTimestampUTC,
		"caller", log.Caller(3),
	)
	if debug {
		logger = level.NewFilter(logger, level.AllowDebug())
	}
	return logger
}

func ConsoleLog() {

}
