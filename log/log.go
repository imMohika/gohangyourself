package log

import (
	"bytes"
	"fmt"
	"github.com/pterm/pterm"
	"log/slog"
)

var handler = pterm.NewSlogHandler(&pterm.DefaultLogger)
var logger = slog.New(handler)

// todo)) get this from global flag
var i = true

func Error(err error, msg string, args ...any) {
	allArgs := make([]any, 0, len(args)+2)
	allArgs = append(allArgs, args...)
	allArgs = append(allArgs, "err", err.Error())

	if i {
		pterm.Error.Print(makeString(msg, allArgs...))
	} else {
		logger.Error(msg, allArgs...)
	}
}

func ErrorMsg(msg string, args ...any) {
	if i {
		pterm.Error.Print(makeString(msg, args...))
	} else {
		logger.Error(msg, args...)
	}
}

func Fatal(err error, msg string, args ...any) {
	allArgs := make([]any, 0, len(args)+2)
	allArgs = append(allArgs, args...)
	allArgs = append(allArgs, "err", err.Error())

	if i {
		pterm.Fatal.Print(makeString(msg, allArgs...))
	} else {
		logger.Error(msg, allArgs...)
	}
}

func FatalMsg(msg string, args ...any) {
	if i {
		pterm.Fatal.Print(makeString(msg, args...))
	} else {
		logger.Error(msg, args...)
	}
}

func Info(msg string, args ...any) {
	if i {
		pterm.Info.Print(makeString(msg, args...))
	} else {
		logger.Info(msg, args...)
	}
}

func Warn(msg string, args ...any) {
	if i {
		pterm.Warning.Print(makeString(msg, args...))
	} else {
		logger.Warn(msg, args...)
	}
}

func Debug(msg string, args ...any) {
	if i {
		pterm.Debug.Print(makeString(msg, args...))
	} else {
		logger.Debug(msg, args...)
	}
}

func makeString(msg string, args ...any) string {
	var tmp bytes.Buffer

	tmp.WriteString(msg)
	tmp.WriteString("\n")

	for i := 0; i < len(args); i += 2 {
		key := args[i]
		value := "!BADKEY"

		if i+1 < len(args) {
			value = fmt.Sprintf("%v", args[i+1])
		}

		tmp.WriteString(fmt.Sprintf("%v: %v\n", key, value))
	}

	return tmp.String()
}
