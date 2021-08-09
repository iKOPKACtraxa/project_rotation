package logger

import (
	"fmt"
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Zap struct {
	logger *zap.SugaredLogger
}

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

const (
	debug      = "debug"
	info       = "info"
	warn       = "warn"
	errorlevel = "error"
)

// New returns a new Logger object.
func New(logfile, level string) *Zap {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{logfile}
	switch level {
	case debug:
		cfg.Level.SetLevel(zapcore.DebugLevel)
	case info:
		cfg.Level.SetLevel(zapcore.InfoLevel)
	case warn:
		cfg.Level.SetLevel(zapcore.WarnLevel)
	case errorlevel:
		cfg.Level.SetLevel(zapcore.ErrorLevel)
	}
	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("logger can't start working: %v", err)
	}
	sugar := logger.Sugar()
	return &Zap{sugar}
}

// Debug make a message in logger at Debug-level.
func (l Zap) Debug(args ...interface{}) {
	l.logger.Debug(args)
	fmt.Println(args...)
}

// Info make a message in logger at Info-level.
func (l Zap) Info(args ...interface{}) {
	l.logger.Info(args)
	fmt.Println(args...)
}

// Warn make a message in logger at Warn-level.
func (l Zap) Warn(args ...interface{}) {
	l.logger.Warn(args)
	fmt.Println(args...)
}

// Error make a message in logger at Error-level.
func (l Zap) Error(args ...interface{}) {
	l.logger.Error(args)
	fmt.Println(args...)
}
