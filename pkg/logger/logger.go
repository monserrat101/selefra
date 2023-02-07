package logger

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/selefra/selefra/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"log"
	"os"
	"path/filepath"
)

type Logger struct {
	logger *zap.Logger
	config *Config
	name   string
}

// Logger impl the hclog.Logger interface
var _ hclog.Logger = (*Logger)(nil)

func (l *Logger) Log(level hclog.Level, msg string, args ...interface{}) {
	switch level {
	case hclog.NoLevel:
		return
	case hclog.Trace:
		l.Trace(msg, args...)
	case hclog.Debug:
		l.Debug(msg, args...)
	case hclog.Info:
		l.Info(msg, args...)
	case hclog.Warn:
		l.Warn(msg, args...)
	case hclog.Error:
		l.Error(msg, args...)
	}
}

func (l *Logger) Trace(msg string, args ...interface{}) {
	l.logger.Debug(fmt.Sprintf(msg, args...))
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.logger.Debug(fmt.Sprintf(msg, args...))
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(msg, args...))
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	l.logger.Warn(fmt.Sprintf(msg, args...))
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.logger.Error(fmt.Sprintf(msg, args...))
}

func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.logger.Fatal(fmt.Sprintf(msg, args...))
}

func (l *Logger) IsTrace() bool {
	return false
}

func (l *Logger) IsDebug() bool {
	return l.config.TranslationLevel() <= zapcore.DebugLevel
}

func (l *Logger) IsInfo() bool {
	return l.config.TranslationLevel() <= zapcore.InfoLevel
}

func (l *Logger) IsWarn() bool {
	return l.config.TranslationLevel() <= zapcore.WarnLevel
}

func (l *Logger) IsError() bool {
	return l.config.TranslationLevel() <= zapcore.ErrorLevel
}

func (l *Logger) ImpliedArgs() []interface{} {
	return nil
}

func (l *Logger) With(args ...interface{}) hclog.Logger {
	return l
}

func (l *Logger) Name() string {
	return l.name
}

func (l *Logger) Named(name string) hclog.Logger {
	l.name = name
	return l
}

func (l *Logger) ResetNamed(name string) hclog.Logger {
	return l
}

func (l *Logger) SetLevel(level hclog.Level) {
	return
}

func (l *Logger) StandardLogger(opts *hclog.StandardLoggerOptions) *log.Logger {
	return log.New(l.StandardWriter(opts), "", 0)
}

func (l *Logger) StandardWriter(opts *hclog.StandardLoggerOptions) io.Writer {
	return os.Stdin
}

var defaultLogger *Logger

func NewLogger(c Config) (*Logger, error) {
	logDir := filepath.Join(global.WorkSpace(), c.Directory)
	_, err := os.Stat(logDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(logDir, 0755)
	}
	if err != nil {
		return nil, nil
	}
	errorStack := zap.AddStacktrace(zap.ErrorLevel)

	development := zap.Development()

	logger := zap.New(zapcore.NewTee(c.GetEncoderCore()...), errorStack, development)

	if c.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}

	defaultLogger = &Logger{logger: logger, config: &c}

	return defaultLogger, nil
}
