package todo

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

func SetUpLogger(f *os.File) {
	logger = slog.New(slog.NewJSONHandler(f, nil))
}

func InfoLog(loc, msg string) {
	logger.Info(loc, "Message", msg)
}

func ErrorLog(loc, msg string) {
	logger.Error(loc, "Message", msg)
}

func WarnLog(loc, msg string) {
	logger.Warn(loc, "Message", msg)
}
