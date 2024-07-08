package log

import (
	"github.com/pterm/pterm"
	"log/slog"
	"os"
)

var handler = pterm.NewSlogHandler(&pterm.DefaultLogger)
var logger = slog.New(handler)

func Error(err error, msg string, args ...any) {
	if err != nil {
		allArgs := make([]any, 0, len(args)+2)
		allArgs = append(allArgs, args...)
		allArgs = append(allArgs, "err", err.Error())
		logger.Error(msg, allArgs...)
	}
}

func Fatal(err error, msg string, args ...any) {
	if err != nil {
		allArgs := make([]any, 0, len(args)+2)
		allArgs = append(allArgs, args...)
		allArgs = append(allArgs, "err", err.Error())
		logger.Error(msg, allArgs...)
		os.Exit(1)
	}
}

func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	logger.Warn(msg, args...)
}

func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}
