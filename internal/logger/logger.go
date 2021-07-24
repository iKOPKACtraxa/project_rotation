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
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
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
func (l Zap) Debug(msg string) {
	l.logger.Debug(msg)
	fmt.Println(msg)
}

// Info make a message in logger at Info-level.
func (l Zap) Info(msg string) {
	l.logger.Info(msg)
	fmt.Println(msg)
}

// Warn make a message in logger at Warn-level.
func (l Zap) Warn(msg string) {
	l.logger.Warn(msg)
	fmt.Println(msg)
}

// Error make a message in logger at Error-level.
func (l Zap) Error(msg string) {
	l.logger.Error(msg)
	fmt.Println(msg)
}
