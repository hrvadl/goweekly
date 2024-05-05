package logger

import (
	"fmt"
	"log/slog"
	"os"
)

func Infof(msg string, args ...any) {
	slog.Info(fmt.Sprintf(msg, args...))
}

func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

func Debugf(msg string, args ...any) {
	slog.Debug(fmt.Sprintf(msg, args...))
}

func Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

func Errorf(msg string, args ...any) {
	slog.Error(fmt.Sprintf(msg, args...))
}

func Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

func Fatalf(msg string, args ...any) {
	Errorf(msg, args...)
	os.Exit(1)
}

func Fatal(msg string, args ...any) {
	Error(msg, args...)
	os.Exit(1)
}
