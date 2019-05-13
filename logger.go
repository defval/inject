package inject

import (
	"log"
)

// Logger
type Logger interface {
	// Printf
	Printf(format string, args ...interface{})
}

// defaultLogger
type defaultLogger struct {
}

func (l *defaultLogger) Printf(format string, args ...interface{}) {
	log.Fatalf(format+"\n", args...)
}

// nopLogger
type nopLogger struct {
}

func (l *nopLogger) Printf(format string, args ...interface{}) {}
