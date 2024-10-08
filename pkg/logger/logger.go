package logger

import (
	"log/slog"
	"os"

	"github.com/pttrulez/invest_telega/pkg/logger/slogpretty"
)

type Logger struct {
	Slog *slog.Logger
}

func NewLogger(log *slog.Logger) *Logger {
	if log == nil {
		log = slog.Default()
	}
	return &Logger{
		Slog: log,
	}
}

func (l *Logger) Debug(s string) {
	l.Slog.Debug(s)
}

func (l *Logger) Error(err error) {
	l.Slog.Error(err.Error())
}

func (l *Logger) Info(s string) {
	l.Slog.Info(s)
}

func SetupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
