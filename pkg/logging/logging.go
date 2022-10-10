package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func New() *Logger {
	return &Logger{logrus.New()}
}

func (l *Logger) SetConfig(lvl string) error {
	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		return fmt.Errorf("invalid level: %w", err)
	}
	l.SetLevel(level)

	l.SetReportCaller(true)

	l.Formatter = &logrus.TextFormatter{
		DisableColors:   true,
		FullTimestamp:   true,
		TimestampFormat: "02-01-2006 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			filename := filepath.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
	}

	l.SetOutput(os.Stdout)

	return nil
}
