package todo

//TODO: Improve logger - look at Sprint.F and slog handler

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

func SetUpLogger(f *os.File) {
	logger = slog.New(slog.NewJSONHandler(f, nil))
}

func InfoLog(loc, msg string) {
	if logger != nil {
		logger.Info(loc, "Message", msg)
	}
}

func ErrorLog(loc, msg string) {
	if logger != nil {
		logger.Error(loc, "Message", msg)
	}
}

func WarnLog(loc, msg string) {
	if logger != nil {
		logger.Warn(loc, "Message", msg)
	}
}
